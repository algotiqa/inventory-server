//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package tradeimport

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/algotiqa/inventory-server/pkg/core/process/agentscanner"
)

//=============================================================================

const INFO  = "INFO"
const START = "START"
const BAR   = "BAR"
const TRADE = "TRADE"
const DAILY = "DAILY"

//=============================================================================

func importTradeList(spec *Spec, data []byte) (*agentscanner.TradingSystem,error) {
	tl  := &agentscanner.TradeList{}
	ats := &agentscanner.TradingSystem{}
	ats.TradeLists = append(ats.TradeLists, tl)

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		if err := handleLine(tl, scanner.Text()); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ats, nil
}

//=============================================================================

func handleLine(tl *agentscanner.TradeList, line string) error {
	tokens := strings.Split(line, "|")

	switch tokens[0] {
		case INFO:
			handleInfo(tl, tokens)
		case START:
			tl.OpenTrade = nil
		case BAR:
			if err := handleBar(tl, tokens); err != nil {
				return err
			}
		case TRADE:
			if err := handleTrade(tl, tokens); err != nil {
				return err
			}
		case DAILY:
			if err := handleDaily(tl, tokens); err != nil {
				return err
			}
		default:
			return errors.New("Unknown token: " + tokens[0])
	}

	return nil
}

//=============================================================================

func handleInfo(tl *agentscanner.TradeList, tokens []string) {
	//--- We are not interested in mapping the symbol or trading system
}

//=============================================================================

func handleBar(tl *agentscanner.TradeList, tokens []string) error {
	var err error

	ddate       := tokens[1]
	dtime       := tokens[2]
	grossReturn := tokens[3]
	contracts   := tokens[4]

	eb := &agentscanner.EquityBar{}

	//-----------------------------------------

	eb.Date, err = convertDate(ddate)
	if err != nil {
		return err
	}

	eb.Time, err = parseInt(dtime)
	if err != nil {
		return errors.New("Cannot parse time: " + dtime)
	}

	eb.GrossReturn, err = strconv.ParseFloat(grossReturn, 64)
	if err != nil {
		return errors.New("Cannot parse gross return: " + grossReturn)
	}

	eb.Contracts, err = parseInt(contracts)
	if err != nil {
		return errors.New("Cannot parse contracts: " + contracts)
	}

	//-----------------------------------------

	tl.OpenTrade = append(tl.OpenTrade, eb)
	return nil
}

//=============================================================================

func handleDaily(tl *agentscanner.TradeList, tokens []string) error {
	var err error

	ddate       := tokens[1]
	dtime       := tokens[2]
	grossReturn := tokens[3]

	dr := &agentscanner.DailyReturn{}

	//-----------------------------------------

	dr.Date, err = convertDate(ddate)
	if err != nil {
		return err
	}

	dr.Time, err = parseInt(dtime)
	if err != nil {
		return errors.New("Cannot parse time: " + dtime)
	}

	dr.GrossReturn, err = strconv.ParseFloat(grossReturn, 64)
	if err != nil {
		return errors.New("Cannot parse gross return: " + grossReturn)
	}

	//-----------------------------------------

	tl.DailyReturns = append(tl.DailyReturns, dr)
	return nil
}

//=============================================================================

func handleTrade(tl *agentscanner.TradeList, tokens []string) error {
	var err error

	entryDate    := tokens[1]
	entryTime    := tokens[2]
	entryPrice   := tokens[3]
	entryLabel   := tokens[4]
	exitDate     := tokens[5]
	exitTime     := tokens[6]
	exitPrice    := tokens[7]
	exitLabel    := tokens[8]
	grossReturn  := tokens[9]
	maxContracts := tokens[10]
	position     := tokens[11]

	tr := &agentscanner.Trade{}

	//-----------------------------------------

	tr.EntryDate, err = convertDate(entryDate)
	if err != nil {
		return err
	}

	tr.EntryTime, err = parseInt(entryTime)
	if err != nil {
		return errors.New("Cannot parse entry time: " + entryTime)
	}

	tr.EntryPrice, err = strconv.ParseFloat(entryPrice, 64)
	if err != nil {
		return errors.New("Cannot parse entry price: " + entryPrice)
	}

	tr.EntryLabel = entryLabel

	//-----------------------------------------

	tr.ExitDate, err = convertDate(exitDate)
	if err != nil {
		return err
	}

	tr.ExitTime, err = parseInt(exitTime)
	if err != nil {
		return errors.New("Cannot parse exit time: " + exitTime)
	}

	tr.ExitPrice, err = strconv.ParseFloat(exitPrice, 64)
	if err != nil {
		return errors.New("Cannot parse exit price: " + exitPrice)
	}

	tr.ExitLabel = exitLabel

	//-----------------------------------------

	tr.GrossReturn, err = strconv.ParseFloat(grossReturn, 64)
	if err != nil {
		return errors.New("Cannot parse gross return: " + grossReturn)
	}

	tr.MaxContracts, err = parseInt(maxContracts)
	if err != nil {
		return errors.New("Cannot parse max contracts: " + maxContracts)
	}

	tr.Position, err = parseInt(position)
	if err != nil {
		return errors.New("Cannot parse position: " + position)
	}

	//-----------------------------------------

	tr.Equity    = tl.OpenTrade
	tl.Trades    = append(tl.Trades, tr)
	tl.OpenTrade = nil
	return nil
}

//=============================================================================

func convertDate(date string) (int, error) {
	tokens := strings.Split(date, "/")

	if len(tokens) != 3 {
		return 0, errors.New("Bad format for date: " + date)
	}

	value, err := strconv.ParseInt(tokens[2]+tokens[1]+tokens[0], 10, 32)

	if err != nil {
		return 0, errors.New("Cannot convert date to int: " + date)
	}

	if value < 20000000 || value > 30000000 {
		return 0, errors.New("Date out of range: " + date)
	}

	return int(value), nil
}

//=============================================================================

func parseInt(text string) (int,error) {
	res, err := strconv.ParseInt(text, 10, 32)
	return int(res), err
}

//=============================================================================
