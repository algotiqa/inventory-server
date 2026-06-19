//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
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

	client := req.GetDefaultClient()
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
