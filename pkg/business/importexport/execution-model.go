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
