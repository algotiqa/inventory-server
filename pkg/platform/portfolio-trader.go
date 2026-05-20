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

package platform

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
)

//=============================================================================
//===
//=== ExportedData
//===
//=============================================================================

type ExportedData struct {
	TradingSystems []*EncodedSystem `json:"tradingSystems"`
}

//=============================================================================

type EncodedSystem struct {
	Id       uint   `json:"id"`
	JsonData []byte `json:"jsonData"`
}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func ExportTradingSystemsFromPortfolio(c *auth.Context, ids []uint) (*ExportedData, error) {
	c.Log.Info("ExportTradingSystemsFromPortfolio: Retrieving systems from portfolio trader...")

	var exportedData ExportedData

	client := req.GetClient("bf")
	url := c.Config.(*app.Config).Platform.Portfolio + "/v1/trading-systems/export?"+ addParameters(ids)
	err := req.DoGet(client, url, &exportedData, c.Token)

	if err != nil {
		c.Log.Error("ExportTradingSystemsFromPortfolio: Got an error from portfolio trader", "error", err.Error())
		return nil, req.NewServerError("Cannot communicate with portfolio-trader: %v", err.Error())
	}

	c.Log.Info("ExportTradingSystemsFromPortfolio: Systems loaded", "systems", len(exportedData.TradingSystems))
	return &exportedData,nil
}

//=============================================================================
