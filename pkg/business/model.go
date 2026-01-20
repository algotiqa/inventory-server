//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

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

package business

import (
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/tiq-engine/session"
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
	DataProductId    uint   `json:"dataProductId"     binding:"required"`
	BrokerProductId  uint   `json:"brokerProductId"   binding:"required"`
	TradingSessionId uint   `json:"tradingSessionId"  binding:"required"`
	AgentProfileId   *uint  `json:"agentProfileId"`
	Name             string `json:"name"              binding:"required"`
	Timeframe        int    `json:"timeframe"         binding:"min=1,max=1440"`
	StrategyType     string `json:"strategyType"      binding:"required"`
	Overnight        bool   `json:"overnight"`
	Tags             string `json:"tags"`
	ExternalRef      string `json:"externalRef"`
}

//=============================================================================

type DataProductSpec struct {
	ConnectionId    uint             `json:"connectionId"   binding:"required"`
	ExchangeId      uint             `json:"exchangeId"     binding:"required"`
	Symbol          string           `json:"symbol"         binding:"required"`
	Name            string           `json:"name"           binding:"required"`
	MarketType      string           `json:"marketType"     binding:"required"`
	ProductType     string           `json:"productType"    binding:"required"`
	Months          string           `json:"months"`
	RolloverTrigger db.DPRollTrigger `json:"rolloverTrigger"`
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
	Username string                  `json:"username"`
	Name     string                  `json:"name"`
	Session  *session.TradingSession `json:"session"`
}

//=============================================================================
//===
//=== ProductBroker & ProductData composite structs
//===
//=============================================================================

type BrokerProductExt struct {
	db.BrokerProduct
	Connection db.Connection `json:"connection"`
	Exchange   db.Exchange   `json:"exchange"`
}

//=============================================================================

type DataProductExt struct {
	db.DataProduct
	Connection db.Connection `json:"connection,omitempty"`
	Exchange   db.Exchange   `json:"exchange,omitempty"`
}

//=============================================================================
//===
//=== Messages
//===
//=============================================================================

type TradingSystemMessage struct {
	TradingSystem  *db.TradingSystem  `json:"tradingSystem"`
	DataProduct    *db.DataProduct    `json:"dataProduct"`
	BrokerProduct  *db.BrokerProduct  `json:"brokerProduct"`
	Currency       *db.Currency       `json:"currency"`
	TradingSession *db.TradingSession `json:"tradingSession"`
	AgentProfile   *db.AgentProfile   `json:"agentProfile"`
	Exchange       *db.Exchange       `json:"exchange"`
}

//=============================================================================

type DataProductMessage struct {
	DataProduct db.DataProduct `json:"dataProduct"`
	Connection  db.Connection  `json:"connection"`
	Exchange    db.Exchange    `json:"exchange"`
}

//=============================================================================

type BrokerProductMessage struct {
	BrokerProduct db.BrokerProduct `json:"brokerProduct"`
	Connection    db.Connection    `json:"connection"`
	Exchange      db.Exchange      `json:"exchange"`
	Currency      db.Currency      `json:"currency"`
}

//=============================================================================

// TradingSessionMessage TODO: To be implemented
type TradingSessionMessage struct {
	TradingSession db.TradingSession `json:"tradingSession"`
}

//=============================================================================

// AgentProfileMessage TODO: To be implemented
type AgentProfileMessage struct {
	AgentProfile db.AgentProfile `json:"agentProfile"`
}

//=============================================================================
