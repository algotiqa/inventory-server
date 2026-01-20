//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
	"strconv"
	"time"

	"github.com/algotiqa/core/datatype"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
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

	err := db.RunInTransaction(func(tx *gorm.DB) error {
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

	delay--

	if delay == 0 {
		agentMap[ap.Id] = ap.ScanInterval
		collectFromAgent(ap)
	} else {
		agentMap[ap.Id] = delay
	}
}

//=============================================================================

func collectFromAgent(ap *db.AgentProfile) {
	client := createClient(ap.SslCertRef, ap.SslKeyRef, "ca.crt")
	if client == nil {
		return
	}

	var data []TradingSystem

	err := req.DoGet(client, ap.RemoteUrl, &data, "")

	if err == nil {
		slog.Info("Trades successfully retrieved from agent", "username", ap.Username, "systems", strconv.Itoa(len(data)), "agent", ap.Name)

		_ = db.RunInTransaction(func(tx *gorm.DB) error {
			return enqueueAgentTrades(tx, ap, data)
		})
	} else {
		slog.Error("Cannot connect to agent", "error", err.Error())
	}
}

//=============================================================================

func createClient(agentCert string, agentKey string, caCert string) *http.Client {
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

func enqueueAgentTrades(tx *gorm.DB, ap *db.AgentProfile, agentTss []TradingSystem) error {
	for _, ats := range agentTss {
		ts, err := db.GetTradingSystemByExtRef(tx, ap.Username, ats.Name)

		if err != nil {
			slog.Error("enqueueAgentTrades: Cannot find trading system", "externalRef", ats.Name, "error", err.Error())
			return err
		}

		if ts == nil {
			slog.Warn("Trading system was not found. Skipping", "externalRef", ats.Name, "username", ap.Username)
			continue
		}

		location, err := getLocation(tx, ts)
		if err != nil {
			slog.Warn("Cannot retrieve timezone for trading system. Skipping", "externalRef", ats.Name, "username", ap.Username, "error", err)
			continue
		}

		for _, tl := range ats.TradeLists {
			err = sendTradeList(ts, ats.Name, tl, location)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//=============================================================================

func getLocation(tx *gorm.DB, ts *db.TradingSystem) (*time.Location, error) {
	bp, err := db.GetBrokerProductById(tx, ts.BrokerProductId)
	if err != nil {
		slog.Error("getTimezone: Could not retrieve broker product of TS", "error", err.Error(), "id", ts.Id)
		return nil, err
	}

	ex, err := db.GetExchangeById(tx, bp.ExchangeId)
	if err != nil {
		slog.Error("getTimezone: Could not retrieve exchange of TS", "error", err.Error(), "id", ts.Id)
		return nil, err
	}

	if ex.Timezone == "utc" {
		return time.UTC, nil
	}

	return time.LoadLocation(ex.Timezone)
}

//=============================================================================

func sendTradeList(ts *db.TradingSystem, extRef string, tl *TradeList, location *time.Location) error {
	//--- Collect trades

	var tradeList []*TradeItem

	for _, atr := range tl.Trades {
		tr := createTrade(extRef, atr, location)
		if tr == nil {
			return errors.New("aborted")
		}
		tradeList = append(tradeList, tr)
	}

	//--- Collect daily profits

	var dayList []*DailyProfitItem

	for _, adp := range tl.DailyProfits {
		dp := createDailyProfit(adp, location)
		if dp == nil {
			return errors.New("aborted")
		}
		dayList = append(dayList, dp)
	}

	//--- Send message

	message := TradeListMessage{
		TradingSystemId: ts.Id,
		Trades:          tradeList,
		DailyProfits:    dayList,
	}

	err := msg.SendMessage(msg.ExRuntime, msg.SourceTrade, msg.TypeCreate, message)
	if err != nil {
		slog.Error("sendTradeList: Cannot enqueue trades for trading system", "name", ts.Name, "error", err.Error())
		return err
	} else {
		slog.Info("sendTradeList: Enqueued trades for trading system", "name", ts.Name, "username", ts.Username)
	}

	return nil
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
	exitDate, err2 := parseDate(atr.ExitDate, atr.ExitTime, loc)

	if err1 != nil {
		slog.Error("createTrade: Cannot parse entry date/time", "entryDate", atr.EntryDate, "entryTime", atr.EntryTime, "name", extRef)
		return nil
	}

	if err2 != nil {
		slog.Error("createTrade: Cannot parse exit date/time", "exitDate", atr.ExitDate, "exitTime", atr.ExitTime, "name", extRef)
		return nil
	}

	if atr.Contracts == 0 {
		slog.Error("createTrade: Cannot manage 0 contracts", "name", extRef)
		return nil
	}

	return &TradeItem{
		TradeType:   tradeType,
		EntryDate:   &entryDate,
		EntryPrice:  atr.EntryPrice,
		EntryLabel:  atr.EntryLabel,
		ExitDate:    &exitDate,
		ExitPrice:   atr.ExitPrice,
		ExitLabel:   atr.ExitLabel,
		GrossProfit: atr.GrossProfit,
		Contracts:   atr.Contracts,
	}
}

//=============================================================================

func createDailyProfit(dp *DailyProfit, loc *time.Location) *DailyProfitItem {
	date, err := parseDate(dp.Date, dp.Time, loc)
	if err != nil {
		slog.Error("createDailyProfit: Cannot parse date/time", "date", dp.Date, "time", dp.Time, "error", err)
		return nil
	}

	return &DailyProfitItem{
		Day:         int(datatype.ToIntDate(&date)),
		GrossProfit: dp.GrossProfit,
		Trades:      dp.Trades,
	}
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
