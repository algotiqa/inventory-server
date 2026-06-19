//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/inventory-server/pkg/business/importexport"
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/types"
)

//=============================================================================

type ConnectionSpec struct {
	Code               string `json:"code"       binding:"required"`
	Name               string `json:"name"       binding:"required"`
	SystemCode         string `json:"systemCode" binding:"required"`
	SystemConfigParams string `json:"systemConfigParams"`
}

//=============================================================================

type TradingSystemSpec struct {
	DataProductId    uint       `json:"dataProductId"     binding:"required"`
	BrokerProductId  uint       `json:"brokerProductId"   binding:"required"`
	TradingSessionId uint       `json:"tradingSessionId"  binding:"required"`
	AgentProfileId   *uint      `json:"agentProfileId"`
	Name             string     `json:"name"              binding:"required"`
	Timeframe        int        `json:"timeframe"         binding:"min=1,max=1440"`
	StrategyType     string     `json:"strategyType"      binding:"required"`
	Overnight        bool       `json:"overnight"`
	Tags             string     `json:"tags"`
	ExternalRef      string     `json:"externalRef"`
	InSampleFrom     types.Date `json:"inSampleFrom"      binding:"required"`
	InSampleTo       types.Date `json:"inSampleTo"        binding:"required"`
	EngineCode       string     `json:"engineCode"        binding:"required"`
}

//=============================================================================

type BrokerProductSpec struct {
	ConnectionId     uint    `json:"connectionId"     binding:"required"`
	ExchangeId       uint    `json:"exchangeId"       binding:"required"`
	Symbol           string  `json:"symbol"           binding:"required"`
	Name             string  `json:"name"             binding:"required"`
	PointValue       float32 `json:"pointValue"       binding:"min=0,max=1000000"`
	CostPerOperation float32 `json:"costPerOperation" binding:"min=0,max=10000"`
	MarginValue      float32 `json:"marginValue"      binding:"min=0,max=1000000"`
	Increment        float64 `json:"increment"        binding:"min=0,max=1"`
	MarketType       string  `json:"marketType"       binding:"required"`
	ProductType      string  `json:"productType"      binding:"required"`
}

//=============================================================================

type TradingSession struct {
	db.Common
	Username string                `json:"username"`
	Name     string                `json:"name"`
	Session  *types.TradingSession `json:"session"`
}

//=============================================================================
//===
//=== TradingSystemReloadResponse
//===
//=============================================================================

type TradingSystemReloadResponse struct {
	TradeCount map[string]int `json:"tradeCount"`
}

//=============================================================================

func NewTradingSystemReloadResponse() *TradingSystemReloadResponse {
	return &TradingSystemReloadResponse{
		TradeCount: make(map[string]int),
	}
}

//=============================================================================
//===
//=== ProductBroker & ProductData composite structs
//===
//=============================================================================

type BrokerProductExt struct {
	db.BrokerProduct
	Connection     *db.Connection          `json:"connection"`
	Exchange       *db.Exchange            `json:"exchange"`
	TradingSystems *[]db.TradingSystemFull `json:"tradingSystems"`
}

//=============================================================================

type DataProductExt struct {
	db.DataProduct
	Connection db.Connection `json:"connection,omitempty"`
	Exchange   db.Exchange   `json:"exchange,omitempty"`
}

//=============================================================================

type ImportOverviewSpec struct {
}

//=============================================================================

type ImportExecutionSpec struct {
	Plan *importexport.ImportPlan
}

//=============================================================================
//===
//=== Connections
//===
//=============================================================================

type ConnectionExt struct {
	db.Connection
	DataProducts   []db.DataProductFull	  `json:"dataProducts"`
	BrokerProducts []db.BrokerProductFull `json:"brokerProducts"`
}

//=============================================================================

const (
	DeleteStatusOk             = "ok"
	DeleteStatusConnected      = "connected"
	DeleteStatusDataProducts   = "dataProducts"
	DeleteStatusBrokerProducts = "brokerProducts"
	DeleteStatusTradingSystems = "tradingSystems"
)

//=============================================================================
