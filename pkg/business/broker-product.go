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

func GetBrokerProducts(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.BrokerProductFull, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	if details {
		return db.GetBrokerProductsFull(tx, filter, offset, limit)
	}

	return db.GetBrokerProducts(tx, filter, offset, limit)
}

//=============================================================================

func GetBrokerProductById(tx *gorm.DB, c *auth.Context, id uint) (*BrokerProductExt, error) {
	c.Log.Info("GetBrokerProductById: Getting a broker product", "id", id)

	bp, err := getBrokerProductAndCheckAccess(tx, c, id, "GetBrokerProductById")
	if err != nil {
		return nil, err
	}

	//--- Get connection

	conn, err := db.GetConnectionById(tx, bp.ConnectionId)
	if err != nil {
		c.Log.Error("GetBrokerProductById: Could not retrieve connection", "error", err.Error())
		return nil, err
	}

	//--- Get exchange

	ex, err := db.GetExchangeById(tx, bp.ExchangeId)
	if err != nil {
		c.Log.Error("GetBrokerProductById: Could not retrieve exchange", "error", err.Error())
		return nil, err
	}

	//--- Put all together

	bpe := BrokerProductExt{
		BrokerProduct: *bp,
		Connection:    *conn,
		Exchange:      *ex,
	}

	return &bpe, nil
}

//=============================================================================

func AddBrokerProduct(tx *gorm.DB, c *auth.Context, bps *BrokerProductSpec) (*db.BrokerProduct, error) {
	c.Log.Info("AddBrokerProduct: Adding a new broker product", "symbol", bps.Symbol, "name", bps.Name)

	var pb db.BrokerProduct
	pb.ConnectionId = bps.ConnectionId
	pb.ExchangeId = bps.ExchangeId
	pb.Username = c.Session.Username
	pb.Symbol = bps.Symbol
	pb.Name = bps.Name
	pb.PointValue = bps.PointValue
	pb.CostPerOperation = bps.CostPerOperation
	pb.MarginValue = bps.MarginValue
	pb.Increment = bps.Increment
	pb.MarketType = bps.MarketType
	pb.ProductType = bps.ProductType

	err := db.AddBrokerProduct(tx, &pb)

	if err != nil {
		c.Log.Error("AddBrokerProduct: Could not add a new broker product", "error", err.Error())
		return nil, err
	}

	err = sendBrokerProductChangeMessage(tx, c, &pb, msg.TypeCreate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("AddBrokerProduct: Broker product added", "symbol", pb.Symbol, "id", pb.Id)
	return &pb, err
}

//=============================================================================

func UpdateBrokerProduct(tx *gorm.DB, c *auth.Context, id uint, pbs *BrokerProductSpec) (*db.BrokerProduct, error) {
	c.Log.Info("UpdateBrokerProduct: Updating a broker product", "id", id, "name", pbs.Name)

	pb, err := getBrokerProductAndCheckAccess(tx, c, id, "UpdateBrokerProduct")
	if err != nil {
		return nil, err
	}

	pb.ExchangeId = pbs.ExchangeId
	pb.Symbol = pbs.Symbol
	pb.Name = pbs.Name
	pb.PointValue = pbs.PointValue
	pb.CostPerOperation = pbs.CostPerOperation
	pb.MarginValue = pbs.MarginValue
	pb.Increment = pbs.Increment
	pb.MarketType = pbs.MarketType
	pb.ProductType = pbs.ProductType

	err = db.UpdateBrokerProduct(tx, pb)
	if err != nil {
		return nil, err
	}

	err = sendBrokerProductChangeMessage(tx, c, pb, msg.TypeUpdate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("UpdateBrokerProduct: Broker product updated", "id", pb.Id, "name", pb.Name)
	return pb, err
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getBrokerProductAndCheckAccess(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.BrokerProduct, error) {
	pb, err := db.GetBrokerProductById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve broker product", "error", err.Error())
		return nil, err
	}

	if pb == nil {
		c.Log.Error(function+": Broker product was not found", "id", id)
		return nil, req.NewNotFoundError("Broker product was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if pb.Username != c.Session.Username {
			c.Log.Error(function+": Broker product not owned by user", "id", id)
			return nil, req.NewForbiddenError("Broker product is not owned by user: %v", id)
		}
	}

	return pb, nil
}

//=============================================================================

func sendBrokerProductChangeMessage(tx *gorm.DB, c *auth.Context, pb *db.BrokerProduct, msgType int) error {

	var exc *db.Exchange
	var cur *db.Currency

	conn, err := db.GetConnectionById(tx, pb.ConnectionId)
	if err != nil {
		c.Log.Error("[Add|Update]BrokerProduct: Could not retrieve connection", "error", err.Error())
		return err
	}

	exc, err = db.GetExchangeById(tx, pb.ExchangeId)
	if err != nil {
		c.Log.Error("[Add|Update]BrokerProduct: Could not retrieve exchange", "error", err.Error())
		return err
	}

	cur, err = db.GetCurrencyById(tx, exc.CurrencyId)
	if err != nil {
		c.Log.Error("[Add|Update]BrokerProduct: Could not retrieve currency", "error", err.Error())
		return err
	}

	pbm := BrokerProductMessage{*pb, *conn, *exc, *cur}
	err = msg.SendMessage(msg.ExInventory, msg.SourceBrokerProduct, msgType, &pbm)

	if err != nil {
		c.Log.Error("[Add|Update]BrokerProduct: Could not publish the update message", "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================
