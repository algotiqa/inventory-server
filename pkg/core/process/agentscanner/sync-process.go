//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package agentscanner

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

const (
	DevString = ".dev."
)

//=============================================================================

var agentMap map[uint]int = map[uint]int{}

//=============================================================================

func Init(cfg *app.Config) *time.Ticker {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		time.Sleep(2 * time.Second)
		run()

		for range ticker.C {
			run()
		}
	}()

	return ticker
}

//=============================================================================

func run() {
	agents, err := getAgentProfiles()
	if err != nil {
		slog.Error("Cannot retrieve agent profiles", "error", err)
		return
	}

	slog.Info("Starting sync process with agents")

	for _, ap := range *agents {
		runAgent(&ap)
	}

	slog.Info("Ending sync process")
}

//=============================================================================

func getAgentProfiles() (*[]db.AgentProfile, error) {
	filter := map[string]any{}
	var list *[]db.AgentProfile

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		var err error
		list, err = db.GetAgentProfiles(tx, filter, 0, 100000)
		return err
	})

	return list, err
}

//=============================================================================

func runAgent(ap *db.AgentProfile) {
	delay, found := agentMap[ap.Id]
	if !found {
		agentMap[ap.Id] = ap.ScanInterval
		delay = ap.ScanInterval
	}

	//--- 0 disables scanning

	if delay != 0 {
		delay--

		if delay == 0 {
			delay = ap.ScanInterval
			collectFromAgent(ap)
		}

		agentMap[ap.Id] = delay
	}
}

//=============================================================================

func collectFromAgent(ap *db.AgentProfile) {
	client := CreateClient(ap.SslCertRef, ap.SslKeyRef, "ca.crt")
	if client == nil {
		return
	}

	var names []string

	err := req.DoGet(client, ap.RemoteUrl + UrlTradingSystems, &names, "")
	if err != nil {
		slog.Error("Cannot connect to agent", "error", err.Error())
		return
	}

	slog.Info("Trading system names successfully retrieved from agent", "username", ap.Username, "systems", len(names), "agent", ap.Name)

	var ts TradingSystem

	for _, name := range names {
		rq := &TradingSystemRequest{ name }
		err = req.DoPost(client, ap.RemoteUrl + UrlTradingSystems, &rq, &ts, "")
		if err != nil {
			slog.Error("Cannot connect to agent", "error", err.Error(), "system", name)
			continue
		}

		_ = dbms.RunInTransaction(func(tx *gorm.DB) error {
			return enqueueAgentTrades(tx, ap, &ts)
		})
	}
}

//=============================================================================

func CreateClient(agentCert string, agentKey string, caCert string) *http.Client {
	path := "certificate/"

	cert, err := os.ReadFile(path + caCert)
	if err != nil {
		slog.Error("Cannot read agent CA certificate: ", "path", path+caCert)
		return nil
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair(path+agentCert, path+agentKey)
	if err != nil {
		slog.Error("Cannot read agent certificate/private key: ", "certificate", path+agentCert, "key", path+agentKey)
		return nil
	}

	return &http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certificate},
			},
		},
	}
}

//=============================================================================
//=== An error aborts the transaction

func enqueueAgentTrades(tx *gorm.DB, ap *db.AgentProfile, ats *TradingSystem) error {
	ts, err := db.GetTradingSystemByExtRef(tx, ap.Username, ats.Name)

	if err != nil {
		slog.Error("enqueueAgentTrades: Cannot find trading system", "externalRef", ats.Name, "error", err.Error())
		return err
	}

	if ts == nil {
		slog.Warn("Trading system was not found. Skipping", "externalRef", ats.Name, "username", ap.Username)
		return nil
	}

	location, err := GetLocation(tx, ts)
	if err != nil {
		slog.Warn("Cannot retrieve timezone for trading system. Skipping", "externalRef", ats.Name, "username", ap.Username, "error", err)
		return nil
	}

	err = SendTrades(ts, ats, location, false)
	if err != nil {
		return err
	}

	return nil
}

//=============================================================================

func GetLocation(tx *gorm.DB, ts *db.TradingSystem) (*time.Location, error) {
	bp, err := db.GetBrokerProductById(tx, ts.BrokerProductId)
	if err != nil {
		slog.Error("getLocation: Could not retrieve broker product of TS", "error", err.Error(), "id", ts.Id)
		return nil, err
	}
	if bp == nil {
		slog.Error("getLocation: Could not retrieve broker product of TS", "id", ts.Id)
		return nil, errors.New("Could not retrieve broker product of TS")
	}

	ex, err := db.GetExchangeById(tx, bp.ExchangeId)
	if err != nil {
		slog.Error("getLocation: Could not retrieve exchange of TS", "error", err.Error(), "id", ts.Id)
		return nil, err
	}
	if ex == nil {
		slog.Error("getLocation: Could not retrieve exchange of TS", "id", ts.Id)
		return nil, errors.New("Could not retrieve exchange of TS")
	}

	return time.LoadLocation(ex.Timezone)
}

//=============================================================================

func SendTrades(ts *db.TradingSystem, ats *TradingSystem, location *time.Location, reload bool) error {
	//--- We want backtested data first and then last period data

	sortTrades(ats)

	//--- Collect trades

	var tradeList []*TradeItem

	for _, tl := range ats.TradeLists {
		for _, atr := range tl.Trades {
			tr := createTrade(ats.Name, atr, location)
			if tr == nil {
				return errors.New("aborted")
			}
			tradeList = append(tradeList, tr)
		}
	}

	var equity []*EquityBarItem

	if len(ats.TradeLists) > 0 {
		tl := ats.TradeLists[len(ats.TradeLists) -1]
		if tl.OpenTrade != nil {
			equity = createEquity(tl.OpenTrade, location)
			if equity == nil {
				return errors.New("aborted")
			}
		}
	}

	//--- Send message

	message := TradeListMessage{
		TradingSystemId: ts.Id,
		Reload         : reload,
		Trades         : tradeList,
		OpenTrade      : equity,
	}

	err := msg.SendMessage(msg.ExRuntime, msg.SourceTrade, msg.TypeCreate, message, nil)
	if err != nil {
		slog.Error("SendTrades: Cannot enqueue trades for trading system", "name", ts.Name, "error", err.Error())
		return err
	}

	slog.Info("SendTrades: Enqueued trades for trading system", "name", ts.Name, "username", ts.Username)

	return nil
}

//=============================================================================

func sortTrades(ats *TradingSystem) {
	slices.SortFunc(ats.TradeLists, func(i, j *TradeList) int {
		if strings.Contains(i.FileName, DevString) {
			return -1
		}

		return 0
	})
}

//=============================================================================

func createTrade(extRef string, atr *Trade, loc *time.Location) *TradeItem {
	tradeType := "?"

	if atr.Position == 1 {
		tradeType = TradeTypeLong
	} else if atr.Position == -1 {
		tradeType = TradeTypeShort
	} else {
		slog.Error("createTrade: Unknown trade type!", "tradeType", atr.Position, "name", extRef)
		return nil
	}

	entryDate, err1 := parseDate(atr.EntryDate, atr.EntryTime, loc)
	exitDate,  err2 := parseDate(atr.ExitDate,  atr.ExitTime,  loc)

	if err1 != nil {
		slog.Error("createTrade: Cannot parse entry date/time", "entryDate", atr.EntryDate, "entryTime", atr.EntryTime, "name", extRef)
		return nil
	}

	if err2 != nil {
		slog.Error("createTrade: Cannot parse exit date/time", "exitDate", atr.ExitDate, "exitTime", atr.ExitTime, "name", extRef)
		return nil
	}

	if atr.MaxContracts == 0 {
		slog.Error("createTrade: Cannot manage 0 contracts", "name", extRef)
		return nil
	}

	equity := createEquity(atr.Equity, loc)
	if equity == nil {
		return nil
	}

	return &TradeItem{
		TradeType   : tradeType,
		EntryDate   : &entryDate,
		EntryPrice  : atr.EntryPrice,
		EntryLabel  : atr.EntryLabel,
		ExitDate    : &exitDate,
		ExitPrice   : atr.ExitPrice,
		ExitLabel   : atr.ExitLabel,
		GrossReturn : atr.GrossReturn,
		MaxContracts: atr.MaxContracts,
		Equity      : equity,
	}
}

//=============================================================================

func createEquity(list []*EquityBar, loc *time.Location) []*EquityBarItem {
	var result []*EquityBarItem

	for _, eb := range list {
		date, err := parseDate(eb.Date, eb.Time, loc)
		if err != nil {
			slog.Error("createEquity: Cannot parse date/time", "date", eb.Date, "time", eb.Time, "error", err)
			return nil
		}

		ebi := &EquityBarItem{
			Date       : date,
			GrossReturn: eb.GrossReturn,
			Contracts  : eb.Contracts,
		}

		result = append(result, ebi)
	}

	return result
}

//=============================================================================

func parseDate(date int, tim int, loc *time.Location) (time.Time, error) {
	sDate := DateToString(date)
	sTime := TimeToString(tim)

	return time.ParseInLocation(time.DateTime, sDate+" "+sTime, loc)
}

//=============================================================================

func DateToString(date int) string {
	y := date / 10000
	m := (date / 100) % 100
	d := date % 100

	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}

//=============================================================================

func TimeToString(t int) string {
	hh := t / 100
	mm := t % 100

	return fmt.Sprintf("%02d:%02d:00", hh, mm)
}

//=============================================================================
