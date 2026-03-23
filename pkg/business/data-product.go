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

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetDataProducts(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.DataProductFull, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	if details {
		return db.GetDataProductsFull(tx, filter, offset, limit)
	}

	return db.GetDataProducts(tx, filter, offset, limit)
}

//=============================================================================

func GetDataProductById(tx *gorm.DB, c *auth.Context, id uint) (*DataProductExt, error) {
	c.Log.Info("GetDataProductById: Getting a data product", "id", id)

	pd, err := getDataProductAndCheckAccess(tx, c, id, "GetDataProductById")
	if err != nil {
		return nil, err
	}

	//--- Get connection

	conn, err := db.GetConnectionById(tx, pd.ConnectionId)
	if err != nil {
		c.Log.Error("GetDataProductById: Could not retrieve connection", "error", err.Error())
		return nil, err
	}

	//--- Get exchange

	exc, err := db.GetExchangeById(tx, pd.ExchangeId)
	if err != nil {
		c.Log.Error("GetDataProductById: Could not retrieve exchange", "error", err.Error())
		return nil, err
	}

	pde := DataProductExt{
		DataProduct: *pd,
		Connection:  *conn,
		Exchange:    *exc,
	}

	return &pde, nil
}

//=============================================================================

func AddDataProduct(tx *gorm.DB, c *auth.Context, dps *DataProductSpec) (*db.DataProduct, error) {
	c.Log.Info("AddDataProduct: Adding a new data product", "symbol", dps.Symbol, "name", dps.Name)

	if err := dps.validateForAdd(); err != nil {
		return nil, err
	}

	var pd db.DataProduct
	pd.ConnectionId    = dps.ConnectionId
	pd.ExchangeId      = dps.ExchangeId
	pd.Username        = c.Session.Username
	pd.Symbol          = dps.Symbol
	pd.Name            = dps.Name
	pd.MarketType      = dps.MarketType
	pd.ProductType     = dps.ProductType
	pd.Months          = dps.Months
	pd.RolloverTrigger = dps.RolloverTrigger
	pd.SessionId       = dps.SessionId

	err := db.AddDataProduct(tx, &pd)

	if err != nil {
		c.Log.Error("AddDataProduct: Could not add a new data product", "error", err.Error())
		return nil, err
	}

	err = sendDataProductChangeMessage(tx, c, &pd, msg.TypeCreate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("AddDataProduct: Data product added", "symbol", pd.Symbol, "id", pd.Id)
	return &pd, err
}

//=============================================================================

func UpdateDataProduct(tx *gorm.DB, c *auth.Context, id uint, dps *DataProductSpec) (*db.DataProduct, error) {
	c.Log.Info("UpdateDataProduct: Updating a data product", "id", id, "name", dps.Name)

	if err := dps.validateForUpdate(); err != nil {
		return nil, err
	}

	pd, err := getDataProductAndCheckAccess(tx, c, id, "UpdateDataProduct")
	if err != nil {
		return nil, err
	}

	//--- We can't change the exchange and the symbol

	pd.Name        = dps.Name
	pd.MarketType  = dps.MarketType
	pd.ProductType = dps.ProductType

	//TODO: Should we allow to modify these? Some recomputation is required
	//pd.Months          = pds.Months
	//pd.RolloverTrigger = pds.RolloverTrigger

	err = db.UpdateDataProduct(tx, pd)
	if err != nil {
		return nil, err
	}

	err = sendDataProductChangeMessage(tx, c, pd, msg.TypeUpdate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("UpdateDataProduct: Data product updated", "id", pd.Id, "name", pd.Name)
	return pd, err
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getDataProductAndCheckAccess(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.DataProduct, error) {
	pd, err := db.GetDataProductById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve data product", "error", err.Error())
		return nil, err
	}

	if pd == nil {
		c.Log.Error(function+": Data product was not found", "id", id)
		return nil, req.NewNotFoundError("Data product was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if pd.Username != c.Session.Username {
			c.Log.Error(function+": Data product not owned by user", "id", id)
			return nil, req.NewForbiddenError("Data product is not owned by user: %v", id)
		}
	}

	return pd, nil
}

//=============================================================================

func sendDataProductChangeMessage(tx *gorm.DB, c *auth.Context, dp *db.DataProduct, msgType int) error {
	conn, err := db.GetConnectionById(tx, dp.ConnectionId)
	if err != nil {
		c.Log.Error("sendDataProductChangeMessage: Could not retrieve connection", "error", err.Error())
		return err
	}

	exc, err := db.GetExchangeById(tx, dp.ExchangeId)
	if err != nil {
		c.Log.Error("sendDataProductChangeMessage: Could not retrieve exchange", "error", err.Error())
		return err
	}

	sess, err := db.GetTradingSessionById(tx, dp.SessionId)
	if err != nil {
		c.Log.Error("sendDataProductChangeMessage: Could not retrieve trading session", "error", err.Error())
		return err
	}

	if conn == nil || exc == nil || sess == nil {
		c.Log.Error("sendDataProductChangeMessage: Could not retrieve the connection/exchange/session for data product", "id", dp.Id)
		return req.NewServerError("Could not retrieve some information while sending the data product message")
	}

	pdm := DataProductMessage{*dp, *conn, *exc, *sess}
	err = msg.SendMessage(msg.ExInventory, msg.SourceDataProduct, msgType, &pdm)

	if err != nil {
		c.Log.Error("sendDataProductChangeMessage: Could not publish the update message", "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================
