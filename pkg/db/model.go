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

package db

import (
	"time"

	"github.com/algotiqa/core/datatype"
)

//=============================================================================

type Common struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

//=============================================================================

type Currency struct {
	Id           uint             `json:"id"`
	Code         string           `json:"code"`
	Name         string           `json:"name"`
	Symbol       string           `json:"symbol"`
	FirstDate    datatype.IntDate `json:"firstDate"`
	LastDate     datatype.IntDate `json:"lastDate"`
	LastValue    float64          `json:"lastValue"`
	HistoryEnded bool             `json:"historyEnded"`
}

//=============================================================================

type CurrencyHistory struct {
	Id         uint             `json:"id"`
	CurrencyId uint             `json:"currencyId"`
	Date       datatype.IntDate `json:"date"`
	Value      float64          `json:"value"`
}

//=============================================================================

type Exchange struct {
	Id         uint   `json:"id"`
	CurrencyId uint   `json:"currencyId"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Timezone   string `json:"timezone"`
	Url        string `json:"url"`
}

//=============================================================================

type Connection struct {
	Common
	Username             string `json:"username"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	SystemCode           string `json:"systemCode"`
	SystemName           string `json:"systemName"`
	SystemConfigParams   string `json:"systemConfigParams"`
	Connected            bool   `json:"connected"`
	SupportsData         bool   `json:"supportsData"`
	SupportsBroker       bool   `json:"supportsBroker"`
	SupportsMultipleData bool   `json:"supportsMultipleData"`
	SupportsInventory    bool   `json:"supportsInventory"`
}

//=============================================================================

type DPRollTrigger string

const (
	DPRollTriggerSD4  = "sd4"
	DPRollTriggerSD6  = "sd6"
	DPRollTriggerSD30 = "sd30"
)

//-----------------------------------------------------------------------------

type DataProduct struct {
	Common
	ConnectionId    uint          `json:"connectionId"`
	ExchangeId      uint          `json:"exchangeId"`
	Username        string        `json:"username"`
	Symbol          string        `json:"symbol"`
	Name            string        `json:"name"`
	MarketType      string        `json:"marketType"`
	ProductType     string        `json:"productType"`
	Months          string        `json:"months"`
	RolloverTrigger DPRollTrigger `json:"rolloverTrigger"`
}

//=============================================================================

type DataProductFull struct {
	DataProduct
	ConnectionCode string `json:"connectionCode,omitempty"`
	ConnectionName string `json:"connectionName,omitempty"`
	SystemCode     string `json:"systemCode,omitempty"`
	ExchangeCode   string `json:"exchangeCode,omitempty"`
}

//=============================================================================

type BrokerProduct struct {
	Common
	ConnectionId     uint    `json:"connectionId"`
	ExchangeId       uint    `json:"exchangeId"`
	Username         string  `json:"username"`
	Symbol           string  `json:"symbol"`
	Name             string  `json:"name"`
	PointValue       float32 `json:"pointValue"`
	CostPerOperation float32 `json:"costPerOperation"`
	MarginValue      float32 `json:"marginValue"`
	Increment        float64 `json:"increment"`
	MarketType       string  `json:"marketType"`
	ProductType      string  `json:"productType"`
}

//=============================================================================

type BrokerProductFull struct {
	BrokerProduct
	CurrencyCode   string `json:"currencyCode,omitempty"`
	ConnectionCode string `json:"connectionCode,omitempty"`
	ConnectionName string `json:"connectionName,omitempty"`
	SystemCode     string `json:"systemCode,omitempty"`
	ExchangeCode   string `json:"exchangeCode,omitempty"`
}

//=============================================================================

type BrokerInstrument struct {
	Id              uint   `json:"id" gorm:"primaryKey"`
	BrokerProductId uint   `json:"brokerProductId"`
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	ExpirationDate  int    `json:"expirationDate"`
}

//=============================================================================

type TradingSession struct {
	Common
	Username string `json:"username"`
	Name     string `json:"name"`
	Config   string `json:"config"`
}

//=============================================================================

type TradingSystem struct {
	Common
	Username         string           `json:"username"`
	DataProductId    uint             `json:"dataProductId"`
	BrokerProductId  uint             `json:"brokerProductId"`
	TradingSessionId uint             `json:"tradingSessionId"`
	AgentProfileId   *uint            `json:"agentProfileId"`
	Name             string           `json:"name"`
	Timeframe        int              `json:"timeframe"`
	StrategyType     string           `json:"strategyType"`
	Overnight        bool             `json:"overnight"`
	Tags             string           `json:"tags"`
	ExternalRef      string           `json:"externalRef"`
	Finalized        bool             `json:"finalized"`
	InSampleFrom     datatype.IntDate `json:"inSampleFrom"`
	InSampleTo       datatype.IntDate `json:"inSampleTo"`
	EngineCode       string           `json:"engineCode"`
}

//=============================================================================

type TradingSystemFull struct {
	TradingSystem
	DataSymbol     string `json:"dataSymbol,omitempty"`
	BrokerSymbol   string `json:"brokerSymbol,omitempty"`
	TradingSession string `json:"tradingSession,omitempty"`
}

//=============================================================================

type AgentProfile struct {
	Common
	Username     string `json:"username"`
	Name         string `json:"name"`
	RemoteUrl    string `json:"remoteUrl"`
	SslKeyRef    string `json:"sslKeyRef"`
	SslCertRef   string `json:"sslCertRef"`
	ScanInterval int    `json:"scanInterval"`
}

//=============================================================================
//===
//=== Table names
//===
//=============================================================================

func (Currency) TableName() string         { return "currency" }
func (CurrencyHistory) TableName() string  { return "currency_history" }
func (Exchange) TableName() string         { return "exchange" }
func (Connection) TableName() string       { return "connection" }
func (AgentProfile) TableName() string     { return "agent_profile" }
func (DataProduct) TableName() string      { return "data_product" }
func (BrokerProduct) TableName() string    { return "broker_product" }
func (BrokerInstrument) TableName() string { return "broker_instrument" }
func (TradingSession) TableName() string   { return "trading_session" }
func (TradingSystem) TableName() string    { return "trading_system" }

//=============================================================================
