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
	"time"
)

//=============================================================================

const (
	UrlTradingSystems = "/api/v1/trading-systems"
	UrlSystemPing     = "/api/v1/system/ping"
)

//=============================================================================
//===
//=== Request
//===
//=============================================================================

type TradingSystemRequest struct {
	Name string `json:"name"`
}

//=============================================================================
//===
//=== Response
//===
//=============================================================================

type TradingSystem struct {
	Name       string `json:"name"`
	DataSymbol string `json:"dataSymbol"`

	TradeLists []*TradeList `json:"tradeLists"`
}

//=============================================================================

type TradeList struct {
	FileName   string         `json:"fileName"`
	Trades     []*Trade       `json:"trades"`
	OpenTrade  []*EquityBar   `json:"openTrade"`
}

//=============================================================================

type Trade struct {
	EntryDate    int          `json:"entryDate"`
	EntryTime    int          `json:"entryTime"`
	EntryPrice   float64      `json:"entryPrice"`
	EntryLabel   string       `json:"entryLabel"`
	ExitDate     int          `json:"exitDate"`
	ExitTime     int          `json:"exitTime"`
	ExitPrice    float64      `json:"exitPrice"`
	ExitLabel    string       `json:"exitLabel"`
	GrossReturn  float64      `json:"grossReturn"`
	MaxContracts int          `json:"maxContracts"`
	Position     int          `json:"position"`
	Equity       []*EquityBar `json:"equity"`
}

//=============================================================================

type EquityBar struct {
	Date        int     `json:"date"`
	Time        int     `json:"time"`
	GrossReturn float64 `json:"grossReturn"`
	Contracts   int     `json:"contracts"`
}

//=============================================================================
//===
//=== Message
//===
//=============================================================================

type TradeListMessage struct {
	TradingSystemId uint             `json:"tradingSystemId"`
	Reload          bool             `json:"reload"`
	Trades          []*TradeItem     `json:"trades"`
	OpenTrade       []*EquityBarItem `json:"openTrade"`
}

//=============================================================================

const (
	TradeTypeLong  = "LO"
	TradeTypeShort = "SH"
)

//=============================================================================

type TradeItem struct {
	TradeType    string           `json:"tradeType"`
	EntryDate    *time.Time       `json:"entryDate"`
	EntryPrice   float64          `json:"entryPrice"`
	EntryLabel   string           `json:"entryLabel"`
	ExitDate     *time.Time       `json:"exitDate"`
	ExitPrice    float64          `json:"exitPrice"`
	ExitLabel    string           `json:"exitLabel"`
	GrossReturn  float64          `json:"grossReturn"`
	MaxContracts int              `json:"maxContracts"`
	Equity       []*EquityBarItem `json:"equity"`
}

//=============================================================================

type EquityBarItem struct {
	Date        time.Time  `json:"date"`
	GrossReturn float64    `json:"grossReturn"`
	Contracts   int        `json:"contracts"`
}

//=============================================================================
