//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
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
