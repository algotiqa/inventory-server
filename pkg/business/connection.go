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
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/inventory-server/pkg/platform"
	"gorm.io/gorm"
)

//=============================================================================

func GetConnections(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]db.Connection, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	return db.GetConnections(tx, filter, offset, limit)
}

//=============================================================================

func GetConnectionById(tx *gorm.DB, c *auth.Context, id uint) (*ConnectionExt, error) {
	conn, err := db.GetConnectionById(tx, id)
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, req.NewNotFoundError("Connection with id='%v' was not found", id)
	}

	if !c.Session.IsAdmin() {
		if c.Session.Username != conn.Username {
			return nil, req.NewForbiddenError("Connection with id='%v' is not owned by the user", id)
		}
	}

	dps,bps,err := getReferences(tx, id)
	if err != nil {
		return nil, err
	}

	return &ConnectionExt{
		Connection    : *conn,
		DataProducts  : *dps,
		BrokerProducts: *bps,
	}, nil
}

//=============================================================================

func AddConnection(tx *gorm.DB, c *auth.Context, cs *ConnectionSpec) (*db.Connection, error) {
	c.Log.Info("AddConnection: Adding a new connection", "code", cs.Code, "name", cs.Name)

	sys, err := platform.GetSystem(c, cs.SystemCode)
	if err != nil {
		c.Log.Info("AddConnection: Unable to retrieve the system", "code", cs.SystemCode)
		return nil, err
	}

	if sys == nil {
		c.Log.Info("AddConnection: System was not found", "code", cs.SystemCode)
		return nil, req.NewNotFoundError("System not found: %v", cs.SystemCode)
	}

	var conn db.Connection
	conn.Username             = c.Session.Username
	conn.Code                 = cs.Code
	conn.Name                 = cs.Name
	conn.SystemCode           = cs.SystemCode
	conn.SystemConfigParams   = cs.SystemConfigParams
	conn.SystemName           = sys.Name
	conn.SupportsData         = sys.SupportsData
	conn.SupportsBroker       = sys.SupportsBroker
	conn.SupportsMultipleData = sys.SupportsMultipleData
	conn.SupportsInventory    = sys.SupportsInventory
	conn.Connected            = conn.SupportsMultipleData

	err = db.AddConnection(tx, &conn)

	if err != nil {
		c.Log.Error("AddConnection: Could not add a new connection", "error", err.Error())
		return nil, err
	}

	c.Log.Info("AddConnection: Connection added", "code", cs.Code, "id", conn.Id)
	return &conn, err
}

//=============================================================================

func UpdateConnection(tx *gorm.DB, c *auth.Context, id uint, cs *ConnectionSpec) (*db.Connection, error) {
	c.Log.Info("UpdateConnection: Updating a connection", "id", id, "name", cs.Name)

	conn, err := getConnection(tx, c, id, "UpdateConnection")
	if err != nil {
		return nil, err
	}

	conn.Name = cs.Name
	conn.SystemConfigParams = cs.SystemConfigParams

	err = db.UpdateConnection(tx, conn)
	if err != nil {
		return nil, err
	}

	c.Log.Info("UpdateConnection: Connection updated", "id", conn.Id, "name", conn.Name)
	return conn, err
}

//=============================================================================

const (
	DeleteStatusOk             = "ok"
	DeleteStatusConnected      = "connected"
	DeleteStatusDataProducts   = "dataProducts"
	DeleteStatusBrokerProducts = "brokerProducts"
)

//-----------------------------------------------------------------------------

func DeleteConnection(tx *gorm.DB, c *auth.Context, id uint) (string, error) {
	c.Log.Info("DeleteConnection: Deleting connection", "id", id)

	conn, err := getConnection(tx, c, id, "DeleteConnection")
	if err != nil {
		return "", err
	}

	//--- We cannot delete when we are connected

	if conn.Connected && !conn.SupportsMultipleData {
		return DeleteStatusConnected, nil
	}

	//--- Check if there are references (not the efficient way, but...)

	dps,bps,err := getReferences(tx, id)
	if err != nil {
		return "", err
	}
	if len(*dps) > 0 {
		return DeleteStatusDataProducts, err
	}
	if len(*bps) > 0 {
		return DeleteStatusBrokerProducts, err
	}

	//--- Proper delete

	err = db.DeleteConnection(tx, id)
	if err != nil {
		c.Log.Error("DeleteConnection: Cannot delete connection", "id", id, "error", err.Error())
		return "", req.NewServerErrorByError(err)
	}

	c.Log.Info("DeleteConnection: Connection deleted", "id", id, "name", conn.Name)
	return DeleteStatusOk, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getConnection(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.Connection, error) {
	conn, err := db.GetConnectionById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve connection", "error", err.Error())
		return nil, err
	}

	if conn == nil {
		c.Log.Error(function+": Connection was not found", "id", id)
		return nil, req.NewNotFoundError("Connection was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if conn.Username != c.Session.Username {
			c.Log.Error(function+": Connection not owned by user", "id", id)
			return nil, req.NewForbiddenError("Connection is not owned by user: %v", id)
		}
	}

	return conn, nil
}

//=============================================================================

func getReferences(tx *gorm.DB, id uint) (*[]db.DataProductFull, *[]db.BrokerProductFull, error) {
	filter := map[string]any{}
	filter["connection_id"] = id

	dps,err := db.GetDataProductsFull(tx, filter, 0, 5000)
	if err != nil {
		return nil, nil, err
	}

	bps,err := db.GetBrokerProductsFull(tx, filter, 0, 5000)
	if err != nil {
		return nil, nil, err
	}

	return dps, bps, nil
}

//=============================================================================
