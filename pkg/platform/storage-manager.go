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
	"io"
	"net/http"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/app"
)

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func ExportTradingSystemsFromStorage(c *auth.Context, ids []uint) ([]byte, error) {
	c.Log.Info("ExportTradingSystemsFromStorage: Retrieving systems from portfolio trader...")

	data,err := getData(c, ids)
	if err != nil {
		c.Log.Error("ExportTradingSystemsFromStorage: Got an error from storage manager", "error", err.Error())
		return nil, req.NewServerError("Cannot communicate with storage-manager: %v", err.Error())
	}

	c.Log.Info("ExportTradingSystemsFromStorage: Systems loaded", "zipSize", len(data))
	return data,nil
}

//=============================================================================

func getData(c *auth.Context, ids []uint) ([]byte, error) {
	client := req.GetDefaultClient()
	url    := c.Config.(*app.Config).Platform.Storage + "/v1/trading-systems/export?"+ addParameters(ids)

	rq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil,err
	}

	rq.Header.Set("Authorization", "Bearer "+ c.Token)

	res, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

//=============================================================================
