//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/core/messaging"
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/inventory-server/pkg/platform"
	"gorm.io/gorm"
)

//=============================================================================

func GetAccounts(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.AccountFull, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	if details {
		return db.GetAccountsFull(tx, filter, offset, limit)
	}

	return db.GetAccounts(tx, filter, offset, limit)
}

//=============================================================================

func GetAccountById(tx *gorm.DB, c *auth.Context, id uint) (*AccountExt, error) {
	c.Log.Info("GetAccountById: Getting an account", "id", id)

	a, err := getAccount(tx, c, id, "GetAccountById")
	if err != nil {
		return nil, err
	}

	//--- Get connection

	conn, err := db.GetConnectionById(tx, a.ConnectionId)
	if err != nil {
		c.Log.Error("GetAccountById: Could not retrieve connection", "error", err.Error())
		return nil, err
	}

	//--- Get currency

	curr, err := db.GetCurrencyById(tx, a.CurrencyId)
	if err != nil {
		c.Log.Error("GetAccountById: Could not retrieve currency", "error", err.Error())
		return nil, err
	}

	//--- Get portfolios

	filter := make(map[string]any)
	filter["account_id"] = id
	list,err := db.GetPortfoliosFull(tx, filter, 0, 5000)
	if err != nil {
		c.Log.Error("GetAccountById: Could not retrieve portfolios", "error", err.Error())
		return nil, err
	}

	//--- Put all together

	ae := AccountExt{
		Account   : *a,
		Connection: conn,
		Currency  : curr,
		Portfolios: list,
	}

	return &ae, nil
}

//=============================================================================

func AddAccount(tx *gorm.DB, c *auth.Context, as *AccountSpec) (*db.Account, error) {
	c.Log.Info("AddAccount: Adding a new account", "code", as.Code, "name", as.Name)

	conn, err := getConnection(tx, c, as.ConnectionId, "AddAccount")
	if err != nil {
		return nil,err
	}

	var a db.Account
	a.Username        = c.Session.Username
	a.ConnectionId    = as.ConnectionId
	a.Code            = as.Code
	a.Name            = as.Name
	a.SupportsAccount = conn.SupportsAccount

	if conn.SupportsAccount {
		acc,erra := findAccount(c, conn.Code, as.Code)
		if erra != nil {
			return nil,erra
		}

		curr,errc := findCurrency(tx, acc.CurrencyCode)
		if errc != nil {
			return nil,errc
		}

		a.CurrencyId     = curr.Id
		a.CurrentCapital = acc.Equity
		a.StatusMessage  = acc.StatusMessage
	} else {
		a.CurrencyId     = as.CurrencyId
		a.CurrentCapital = as.CurrentCapital
		a.StatusMessage  = ""
	}

	err = db.AddAccount(tx, &a)

	if err != nil {
		c.Log.Error("AddAccount: Could not add a new account", "error", err.Error())
		return nil, err
	}

	err = sendAccountChangeMessage(tx, c, &a, msg.TypeCreate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("AddAccount: Account added", "code", a.Code, "id", a.Id)
	return &a, err
}

//=============================================================================

func UpdateAccount(tx *gorm.DB, c *auth.Context, id uint, as *AccountSpec) (*db.Account, error) {
	c.Log.Info("UpdateAccount: Updating an account", "id", id, "name", as.Name)

	a, err := getAccount(tx, c, id, "UpdateAccount")
	if err != nil {
		return nil, err
	}

	a.Name = as.Name

	if !a.SupportsAccount {
		a.Code           = as.Code
		a.CurrentCapital = as.CurrentCapital
		a.CurrencyId     = as.CurrencyId
	}

	err = db.UpdateAccount(tx, a)
	if err != nil {
		return nil, err
	}

	err = sendAccountChangeMessage(tx, c, a, msg.TypeUpdate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("UpdateAccount: Account updated", "id", a.Id, "name", a.Name)
	return a, err
}

//=============================================================================

func DeleteAccount(tx *gorm.DB, c *auth.Context, id uint) (string, error) {
	c.Log.Info("DeleteAccount: Deleting account", "id", id)

	a, err := getAccount(tx, c, id, "DeleteAccount")
	if err != nil {
		return "", err
	}

	//--- Check if there are references (not the efficient way, but...)

	filter := map[string]any{}
	filter["account_id"] = id

	ps,err := db.GetPortfolios(tx, filter, 0, 5000)
	if err != nil {
		return "", err
	}
	if len(*ps) > 0 {
		return DeleteStatusPortfolios, err
	}

	//--- Proper delete

	err = db.DeleteAccount(tx, id)
	if err != nil {
		c.Log.Error("DeleteAccount: Cannot delete account", "id", id, "error", err.Error())
		return "", req.NewServerErrorByError(err)
	}

	err = sendAccountChangeMessage(tx, c, a, msg.TypeDelete)
	if err != nil {
		return "", err
	}

	c.Log.Info("DeleteAccount: Account deleted", "id", id, "name", a.Name)
	return DeleteStatusOk, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getAccount(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.Account, error) {
	a, err := db.GetAccountById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve account", "error", err.Error())
		return nil, err
	}

	if a == nil {
		c.Log.Error(function+": Account was not found", "id", id)
		return nil, req.NewNotFoundError("Account was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if a.Username != c.Session.Username {
			c.Log.Error(function+": Account not owned by user", "id", id)
			return nil, req.NewForbiddenError("Account is not owned by user: %v", id)
		}
	}

	return a, nil
}

//=============================================================================

func sendAccountChangeMessage(tx *gorm.DB, c *auth.Context, a *db.Account, msgType int) error {
	var cur *db.Currency

	conn, err := db.GetConnectionById(tx, a.ConnectionId)
	if err != nil {
		c.Log.Error("[Add|Update]Account: Could not retrieve connection", "error", err.Error())
		return err
	}

	cur, err = db.GetCurrencyById(tx, a.CurrencyId)
	if err != nil {
		c.Log.Error("[Add|Update]Account: Could not retrieve currency", "error", err.Error())
		return err
	}

	am := messaging.AccountMessage{
		Account   : a,
		Connection: conn,
		Currency  : cur,
	}

	err = msg.SendMessage(msg.ExInventory, msg.SourceAccount, msgType, &am, tx)

	if err != nil {
		c.Log.Error("[Add|Update]Account: Could not publish the update message", "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================

func findAccount(c *auth.Context, connCode, accCode string) (*platform.Account, error) {
	list,err := platform.GetAccounts(c, connCode)
	if err != nil {
		return nil,err
	}

	for _, a := range *list {
		if a.Code == accCode {
			return &a,nil
		}
	}

	return nil, req.NewNotFoundError("Account not found on connection: %v", accCode)
}

//=============================================================================

func findCurrency(tx *gorm.DB, code string) (*db.Currency,error){
	curr,err := db.GetCurrencyByCode(tx,code)
	if err != nil {
		return nil,err
	}

	if curr == nil {
		return nil,req.NewNotFoundError("Currency not found on platform: %v", code)
	}

	return curr, nil
}

//=============================================================================
