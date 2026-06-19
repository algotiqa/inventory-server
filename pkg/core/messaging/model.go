//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package messaging

import (
	"github.com/algotiqa/inventory-server/pkg/db"
)

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
	PortfolioPack  []byte             `json:"portfolioPack"`
	StoragePack    []byte             `json:"storagePack"`
}

//=============================================================================

type DataProductMessage struct {
	DataProduct    *db.DataProduct    `json:"dataProduct"`
	Connection     *db.Connection     `json:"connection"`
	Exchange       *db.Exchange       `json:"exchange"`
	TradingSession *db.TradingSession `json:"tradingSession"`
}

//=============================================================================

type BrokerProductMessage struct {
	BrokerProduct *db.BrokerProduct `json:"brokerProduct"`
	Connection    *db.Connection    `json:"connection"`
	Exchange      *db.Exchange      `json:"exchange"`
	Currency      *db.Currency      `json:"currency"`
}

//=============================================================================

// TradingSessionMessage TODO: To be implemented
type TradingSessionMessage struct {
	TradingSession *db.TradingSession `json:"tradingSession"`
}

//=============================================================================

// AgentProfileMessage TODO: To be implemented
type AgentProfileMessage struct {
	AgentProfile *db.AgentProfile `json:"agentProfile"`
}

//=============================================================================
