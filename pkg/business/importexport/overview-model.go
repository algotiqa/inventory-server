//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

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

package importexport

import "github.com/algotiqa/inventory-server/pkg/db"

//=============================================================================
//===
//=== Response
//===
//=============================================================================

type ImportOverviewResponse struct {
	TradingSystems  []*TradingSystemItem `json:"tradingSystems"`
	ReferencedItems []*ReferencedItem    `json:"referencedItems"`
}

//=============================================================================

type TradingSystemItem struct {
	Id         uint    `json:"id"`
	Name       string  `json:"name"`
	Timeframe  int     `json:"timeframe"`
}

//=============================================================================

type ReferencedItemStatus int

const (
	RIStatusNew          ReferencedItemStatus = 0
	RIStatusExisting     ReferencedItemStatus = 1
	RIStatusNoConnection ReferencedItemStatus = 2
)

//-----------------------------------------------------------------------------

type ReferencedItemType int

const (
	ReferencedItemTypeData    ReferencedItemType = 0
	ReferencedItemTypeBroker  ReferencedItemType = 1
	ReferencedItemTypeProfile ReferencedItemType = 2
	ReferencedItemTypeSession ReferencedItemType = 3
)

//-----------------------------------------------------------------------------

type ReferencedItem struct {
	Id           uint                 `json:"id"`
	Symbol       string               `json:"symbol"`
	Name         string               `json:"name"`
	SystemCode   string               `json:"systemCode"`
	ExchangeCode string               `json:"exchangeCode"`
	ItemType     ReferencedItemType   `json:"itemType"`
	Status       ReferencedItemStatus `json:"status"`
	Options      []*ReferencedOption  `json:"options"`
	MappedTo     uint                 `json:"mappedTo"`

	sessionConfig string
}

//=============================================================================

type ReferencedOption struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	MatchNotes string `json:"matchNotes"`

	dataProduct    *db.DataProduct
	brokerProduct  *db.BrokerProduct
	agentProfile   *db.AgentProfile
	tradingSession *db.TradingSession
}

//=============================================================================
