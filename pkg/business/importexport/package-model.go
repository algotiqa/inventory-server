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

import (
	"time"

	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/types"
)

//=============================================================================
//===
//=== In memory structure for package
//===
//=============================================================================

type InMemoryPackage struct {
	Metadata  *Metadata
	Data      *ExportedData
}

//=============================================================================
//===
//=== Metadata
//===
//=============================================================================

type Metadata struct {
	Version    *Version  `json:"version"`
	ExportDate time.Time `json:"exportDate"`
}

//=============================================================================

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

//=============================================================================
//===
//=== ExportData
//===
//=============================================================================

type ExportedData struct {
	DataProducts    []*DataProduct    `json:"dataProducts"`
	BrokerProducts  []*BrokerProduct  `json:"brokerProducts"`
	TradingSessions []*TradingSession `json:"tradingSessions"`
	AgentProfiles   []*AgentProfile   `json:"agentProfiles"`
	TradingSystems  []*TradingSystem  `json:"tradingSystems"`
}

//=============================================================================

type DataProduct struct {
	Id              uint             `json:"id"`
	ExchangeId      uint             `json:"exchangeId"`
	ExchangeCode    string           `json:"exchangeCode"`
	Symbol          string           `json:"symbol"`
	Name            string           `json:"name"`
	MarketType      string           `json:"marketType"`
	ProductType     string           `json:"productType"`
	Months          string           `json:"months"`
	RolloverTrigger db.DPRollTrigger `json:"rolloverTrigger"`
	ConnectionCode  string           `json:"connectionCode,omitempty"`
	SystemCode      string           `json:"systemCode,omitempty"`
	SessionId       uint             `json:"sessionId"`
}

//=============================================================================

type BrokerProduct struct {
	Id               uint       `json:"id"`
	ExchangeId       uint       `json:"exchangeId"`
	ExchangeCode     string     `json:"exchangeCode"`
	Symbol           string     `json:"symbol"`
	Name             string     `json:"name"`
	PointValue       float32    `json:"pointValue"`
	CostPerOperation float32    `json:"costPerOperation"`
	MarginValue      float32    `json:"marginValue"`
	Increment        float64    `json:"increment"`
	MarketType       string     `json:"marketType"`
	ProductType      string     `json:"productType"`
	ConnectionCode   string     `json:"connectionCode,omitempty"`
	SystemCode       string     `json:"systemCode,omitempty"`
}

//=============================================================================

type TradingSession struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	Session    string     `json:"session"`
}

//=============================================================================

type AgentProfile struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//=============================================================================

type TradingSystem struct {
	Id               uint       `json:"id"`
	DataProductId    uint       `json:"dataProductId"`
	BrokerProductId  uint       `json:"brokerProductId"`
	TradingSessionId uint       `json:"tradingSessionId"`
	AgentProfileId   *uint      `json:"agentProfileId"`
	Name             string     `json:"name"`
	Timeframe        int        `json:"timeframe"`
	StrategyType     string     `json:"strategyType"`
	Overnight        bool       `json:"overnight"`
	Tags             string     `json:"tags"`
	ExternalRef      string     `json:"externalRef"`
	Finalized        bool       `json:"finalized"`
	InSampleFrom     types.Date `json:"inSampleFrom"`
	InSampleTo       types.Date `json:"inSampleTo"`
	EngineCode       string     `json:"engineCode"`

	portfolioData    []byte
	storageData      []byte
}

//=============================================================================
//===
//=== Constructors
//===
//=============================================================================

func NewDataProduct(dp *db.DataProductFull) *DataProduct {
	return &DataProduct{
		Id             : dp.Id,
		ExchangeId     : dp.ExchangeId,
		ExchangeCode   : dp.ExchangeCode,
		Symbol         : dp.Symbol,
		Name           : dp.Name,
		MarketType     : dp.MarketType,
		ProductType    : dp.ProductType,
		Months         : dp.Months,
		RolloverTrigger: dp.RolloverTrigger,
		ConnectionCode : dp.ConnectionCode,
		SystemCode     : dp.SystemCode,
		SessionId      : dp.SessionId,
	}
}

//=============================================================================

func NewBrokerProduct(bp *db.BrokerProductFull) *BrokerProduct {
	return &BrokerProduct{
		Id              : bp.Id,
		ExchangeId      : bp.ExchangeId,
		ExchangeCode    : bp.ExchangeCode,
		Symbol          : bp.Symbol,
		Name            : bp.Name,
		PointValue      : bp.PointValue,
		CostPerOperation: bp.CostPerOperation,
		MarginValue     : bp.MarginValue,
		Increment       : bp.Increment,
		MarketType      : bp.MarketType,
		ProductType     : bp.ProductType,
		ConnectionCode  : bp.ConnectionCode,
		SystemCode      : bp.SystemCode,
	}
}

//=============================================================================

func NewTradingSession(s *db.TradingSession) *TradingSession {
	return &TradingSession{
		Id        : s.Id,
		Name      : s.Name,
		Session   : s.Session,
	}
}

//=============================================================================

func NewAgentProfile(ap *db.AgentProfile) *AgentProfile {
	return &AgentProfile{
		Id   : ap.Id,
		Name : ap.Name,
	}
}

//=============================================================================

func NewTradingSystem(ts *db.TradingSystem) *TradingSystem {
	return &TradingSystem{
		Id              : ts.Id,
		DataProductId   : ts.DataProductId,
		BrokerProductId : ts.BrokerProductId,
		TradingSessionId: ts.TradingSessionId,
		AgentProfileId  : ts.AgentProfileId,
		Name            : ts.Name,
		Timeframe       : ts.Timeframe,
		StrategyType    : ts.StrategyType,
		Overnight       : ts.Overnight,
		Tags            : ts.Tags,
		ExternalRef     : ts.ExternalRef,
		Finalized       : ts.Finalized,
		InSampleFrom    : ts.InSampleFrom,
		InSampleTo      : ts.InSampleTo,
		EngineCode      : ts.EngineCode,
	}
}

//=============================================================================
