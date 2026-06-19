//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package importexport

import (
	"github.com/algotiqa/inventory-server/pkg/db"
)

//=============================================================================
//===
//=== Plan
//===
//=============================================================================

type ImportPlan struct {
	TradingSystems  []*SelectedTradingSystem `json:"tradingSystems"`
	ReferencedItems []*SelectedReference     `json:"referencedItems"`
}

//=============================================================================

type SelectedTradingSystem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//=============================================================================

type SelectedReference struct {
	Id       uint               `json:"id"`
	ItemType ReferencedItemType `json:"itemType"`
	MappedTo uint               `json:"mappedTo"`
}

//=============================================================================
//===
//=== Result
//===
//=============================================================================

type ImportExecutionResult struct {
	Items []*ImportedItem
}

//=============================================================================

type ImportedItem struct {
	System    *db.TradingSystem
	Data      *db.DataProduct
	Broker    *db.BrokerProduct
	Profile   *db.AgentProfile
	Session   *db.TradingSession
	Exchange  *db.Exchange
	Currency  *db.Currency
	Portfolio []byte
	Storage   []byte
}

//=============================================================================
