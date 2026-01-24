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
	"time"
)

//=============================================================================

const (
	UrlGetTradingSystems   = "/api/v1/trading-systems"
	UrlReloadTradingSystem = "/api/v1/trading-systems/reload"
	UrlSystemPing          = "/api/v1/system/ping"
)

//=============================================================================

type TradingSystem struct {
	Name       string `json:"name"`
	DataSymbol string `json:"dataSymbol"`

	TradeLists []*TradeList `json:"tradeLists"`
}

//=============================================================================

type TradeList struct {
	FileName     string         `json:"fileName"`
	Trades       []*Trade       `json:"trades"`
	DailyProfits []*DailyProfit `json:"dailyProfits"`
}

//=============================================================================

type Trade struct {
	EntryDate   int     `json:"entryDate"`
	EntryTime   int     `json:"entryTime"`
	EntryPrice  float64 `json:"entryPrice"`
	EntryLabel  string  `json:"entryLabel"`
	ExitDate    int     `json:"exitDate"`
	ExitTime    int     `json:"exitTime"`
	ExitPrice   float64 `json:"exitPrice"`
	ExitLabel   string  `json:"exitLabel"`
	GrossProfit float64 `json:"grossProfit"`
	Contracts   int     `json:"contracts"`
	Position    int     `json:"position"`
}

//=============================================================================

type DailyProfit struct {
	Date        int     `json:"date"`
	Time        int     `json:"time"`
	GrossProfit float64 `json:"grossProfit"`
	Trades      int     `json:"trades"`
}

//=============================================================================

type TradeListMessage struct {
	TradingSystemId uint               `json:"tradingSystemId"`
	Reload          bool               `json:"reload"`
	Trades          []*TradeItem       `json:"trades"`
	DailyProfits    []*DailyProfitItem `json:"dailyProfits"`
}

//=============================================================================

const (
	TradeTypeLong  = "LO"
	TradeTypeShort = "SH"
)

//=============================================================================

type TradeItem struct {
	TradeType   string     `json:"tradeType"`
	EntryDate   *time.Time `json:"entryDate"`
	EntryPrice  float64    `json:"entryPrice"`
	EntryLabel  string     `json:"entryLabel"`
	ExitDate    *time.Time `json:"exitDate"`
	ExitPrice   float64    `json:"exitPrice"`
	ExitLabel   string     `json:"exitLabel"`
	GrossProfit float64    `json:"grossProfit"`
	Contracts   int        `json:"contracts"`
}

//=============================================================================

type DailyProfitItem struct {
	Day         int     `json:"day"`
	GrossProfit float64 `json:"grossProfit"`
	Trades      int     `json:"trades"`
}

//=============================================================================
