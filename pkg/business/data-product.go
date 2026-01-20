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

func AddDataProduct(tx *gorm.DB, c *auth.Context, pds *DataProductSpec) (*db.DataProduct, error) {
	c.Log.Info("AddDataProduct: Adding a new data product", "symbol", pds.Symbol, "name", pds.Name)

	var pd db.DataProduct
	pd.ConnectionId = pds.ConnectionId
	pd.ExchangeId = pds.ExchangeId
	pd.Username = c.Session.Username
	pd.Symbol = pds.Symbol
	pd.Name = pds.Name
	pd.MarketType = pds.MarketType
	pd.ProductType = pds.ProductType
	pd.Months = pds.Months
	pd.RolloverTrigger = pds.RolloverTrigger

	//TODO: validate rollover trigger

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

func UpdateDataProduct(tx *gorm.DB, c *auth.Context, id uint, pds *DataProductSpec) (*db.DataProduct, error) {
	c.Log.Info("UpdateDataProduct: Updating a data product", "id", id, "name", pds.Name)

	pd, err := getDataProductAndCheckAccess(tx, c, id, "UpdateDataProduct")
	if err != nil {
		return nil, err
	}

	//--- We can't change the exchange and the symbol

	pd.Name = pds.Name
	pd.MarketType = pds.MarketType
	pd.ProductType = pds.ProductType

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

func sendDataProductChangeMessage(tx *gorm.DB, c *auth.Context, pd *db.DataProduct, msgType int) error {
	conn, err := db.GetConnectionById(tx, pd.ConnectionId)
	if err != nil {
		c.Log.Error("[Add|Update]DataProduct: Could not retrieve connection", "error", err.Error())
		return err
	}

	exc, err := db.GetExchangeById(tx, pd.ExchangeId)
	if err != nil {
		c.Log.Error("[Add|Update]DataProduct: Could not retrieve exchange", "error", err.Error())
		return err
	}

	pdm := DataProductMessage{*pd, *conn, *exc}
	err = msg.SendMessage(msg.ExInventory, msg.SourceDataProduct, msgType, &pdm)

	if err != nil {
		c.Log.Error("[Add|Update]DataProduct: Could not publish the update message", "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================
