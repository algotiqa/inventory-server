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
