//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================


package db

import (
	"strconv"
	"time"

	"github.com/algotiqa/types"
)

//=============================================================================

type Common struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

//=============================================================================

type Currency struct {
	Id           uint       `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Symbol       string     `json:"symbol"`
	FirstDate    types.Date `json:"firstDate"`
	LastDate     types.Date `json:"lastDate"`
	LastValue    float64    `json:"lastValue"`
	HistoryEnded bool       `json:"historyEnded"`
}

//=============================================================================

type CurrencyHistory struct {
	Id         uint       `json:"id"`
	CurrencyId uint       `json:"currencyId"`
	Date       types.Date `json:"date"`
	Value      float64    `json:"value"`
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
	SessionId       uint          `json:"sessionId"`
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
	PointValue       float64 `json:"pointValue"`
	CostPerOperation float64 `json:"costPerOperation"`
	MarginValue      float64 `json:"marginValue"`
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
	Session  string `json:"session"`
}

//=============================================================================

type TradingSystem struct {
	Common
	Username         string     `json:"username"`
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
}

//=============================================================================

type TradingSystemFull struct {
	TradingSystem
	DataSymbol     string `json:"dataSymbol,omitempty"`
	BrokerSymbol   string `json:"brokerSymbol,omitempty"`
	TradingSession string `json:"tradingSession,omitempty"`
}

//=============================================================================

const (
	HostTypeWindows = "windows"
	HostTypeLinux   = "linux"
)

//-----------------------------------------------------------------------------

type AgentProfile struct {
	Common
	Username      string `json:"username"`
	Name          string `json:"name"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	ScanInterval  int    `json:"scanInterval"`
	ScanFolder    string `json:"scanFolder"`
	FileExtension string `json:"fileExtension"`
	HostType      string `json:"hostType"`
	SslKey        []byte `json:"sslKey"`
	SslCert       []byte `json:"sslCert"`
}

//-----------------------------------------------------------------------------

func (ap *AgentProfile) RemoteUrl() string {
	return "https://" + ap.Host + ":" + strconv.Itoa(ap.Port)
}

//=============================================================================
//===
//=== Table names
//===
//=============================================================================

func (Currency)         TableName() string { return "currency"          }
func (CurrencyHistory)  TableName() string { return "currency_history"  }
func (Exchange)         TableName() string { return "exchange"          }
func (Connection)       TableName() string { return "connection"        }
func (AgentProfile)     TableName() string { return "agent_profile"     }
func (DataProduct)      TableName() string { return "data_product"      }
func (BrokerProduct)    TableName() string { return "broker_product"    }
func (BrokerInstrument) TableName() string { return "broker_instrument" }
func (TradingSession)   TableName() string { return "trading_session"   }
func (TradingSystem)    TableName() string { return "trading_system"    }

//=============================================================================
