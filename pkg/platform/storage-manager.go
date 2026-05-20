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
	client := req.GetClient("bf")
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
