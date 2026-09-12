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
	"gorm.io/gorm"
)

//=============================================================================

func GetPortfolios(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.PortfolioFull, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	if details {
		return db.GetPortfoliosFull(tx, filter, offset, limit)
	}

	return db.GetPortfolios(tx, filter, offset, limit)
}

//=============================================================================

func GetPortfolioById(tx *gorm.DB, c *auth.Context, id uint) (*PortfolioExt, error) {
	c.Log.Info("GetPortfolioById: Getting a portfolio", "id", id)

	p, err := getPortfolio(tx, c, id, "GetPortfolioById")
	if err != nil {
		return nil, err
	}

	//--- Get account

	acc, err := db.GetAccountById(tx, p.AccountId)
	if err != nil {
		c.Log.Error("GetPortfolioById: Could not retrieve account", "error", err.Error())
		return nil, err
	}

	//--- Get currency

	curr, err := db.GetCurrencyById(tx, acc.CurrencyId)
	if err != nil {
		c.Log.Error("GetPortfolioById: Could not retrieve currency", "error", err.Error())
		return nil, err
	}

	//--- Get trading systems

	filter := make(map[string]any)
	filter["portfolio_id"] = id
	tss,err := db.GetTradingSystemsFull(tx, filter, 0, 5000)
	if err != nil {
		c.Log.Error("GetPortfolioById: Could not retrieve trading systems", "error", err.Error())
		return nil, err
	}

	//--- Put all together

	pe := PortfolioExt{
		Portfolio     : *p,
		Account       : acc,
		Currency      : curr,
		TradingSystems: tss,
	}

	return &pe, nil
}

//=============================================================================

func AddPortfolio(tx *gorm.DB, c *auth.Context, ps *PortfolioSpec) (*db.Portfolio, error) {
	c.Log.Info("AddPortfolio: Adding a new portfolio", "name", ps.Name)

	if err := ps.Validate(); err != nil {
		return nil, err
	}

	//--- Validate if account belong to user
	_,err := getAccount(tx, c, ps.AccountId, "AddPortfolio")
	if err != nil {
		return nil, err
	}

	var p db.Portfolio
	p.Username      = c.Session.Username
	p.Name          = ps.Name
	p.AccountId     = ps.AccountId
	p.Management    = ps.Management
	p.AccountPerc   = ps.AccountPerc
	p.MaxMarginPerc = ps.MaxMarginPerc

	err = db.AddPortfolio(tx, &p)

	if err != nil {
		c.Log.Error("AddPortfolio: Could not add a new portfolio", "error", err.Error())
		return nil, err
	}

	err = sendPortfolioChangeMessage(tx, c, &p, msg.TypeCreate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("AddPortfolio: Portfolio added", "name", p.Name, "id", p.Id)
	return &p, err
}

//=============================================================================

func UpdatePortfolio(tx *gorm.DB, c *auth.Context, id uint, ps *PortfolioSpec) (*db.Portfolio, error) {
	c.Log.Info("UpdatePortfolio: Updating a portfolio", "id", id, "name", ps.Name)

	p, err := getPortfolio(tx, c, id, "UpdatePortfolio")
	if err != nil {
		return nil, err
	}

	if p.Management != ps.Management {
		filter := map[string]any{}
		filter["portfolio_id"] = id

		tss,errs := db.GetTradingSystems(tx, filter, 0, 5000)
		if errs != nil {
			return nil, errs
		}
		if len(*tss) > 0 {
			return nil, req.NewForbiddenError("Cannot change the management of a portfolio if there are trading systems attached")
		}
	}

	p.Name          = ps.Name
	p.Management    = ps.Management
	p.AccountPerc   = ps.AccountPerc
	p.MaxMarginPerc = ps.MaxMarginPerc

	//--- Parent account cannot be changed
	//p.AccountId = ps.AccountId

	err = db.UpdatePortfolio(tx, p)
	if err != nil {
		return nil, err
	}

	err = sendPortfolioChangeMessage(tx, c, p, msg.TypeUpdate)
	if err != nil {
		return nil, err
	}

	c.Log.Info("UpdatePortfolio: Portfolio updated", "id", p.Id, "name", p.Name)
	return p, err
}

//=============================================================================

func DeletePortfolio(tx *gorm.DB, c *auth.Context, id uint) (string, error) {
	c.Log.Info("DeletePortfolio: Deleting portfolio", "id", id)

	p, err := getPortfolio(tx, c, id, "DeletePortfolio")
	if err != nil {
		return "", err
	}

	//--- Check if there are references (not the efficient way, but...)

	filter := map[string]any{}
	filter["portfolio_id"] = id

	tss,err := db.GetTradingSystems(tx, filter, 0, 5000)
	if err != nil {
		return "", err
	}
	if len(*tss) > 0 {
		return DeleteStatusTradingSystems, err
	}

	//--- Proper delete

	err = db.DeletePortfolio(tx, id)
	if err != nil {
		c.Log.Error("DeletePortfolio: Cannot delete portfolio", "id", id, "error", err.Error())
		return "", req.NewServerErrorByError(err)
	}

	err = sendPortfolioChangeMessage(tx, c, p, msg.TypeDelete)
	if err != nil {
		return "", err
	}

	c.Log.Info("DeletePortfolio: Portfolio deleted", "id", id, "name", p.Name)
	return DeleteStatusOk, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getPortfolio(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.Portfolio, error) {
	p, err := db.GetPortfolioById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve portfolio", "error", err.Error())
		return nil, err
	}

	if p == nil {
		c.Log.Error(function+": Portfolio was not found", "id", id)
		return nil, req.NewNotFoundError("Portfolio was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if p.Username != c.Session.Username {
			c.Log.Error(function+": Portfolio not owned by user", "id", id)
			return nil, req.NewForbiddenError("Portfolio is not owned by user: %v", id)
		}
	}

	return p, nil
}

//=============================================================================

func sendPortfolioChangeMessage(tx *gorm.DB, c *auth.Context, p *db.Portfolio, msgType int) error {
	var cur *db.Currency

	acc, err := db.GetAccountById(tx, p.AccountId)
	if err != nil {
		c.Log.Error("[Add|Update]Profile: Could not retrieve account", "error", err.Error())
		return err
	}
	if acc == nil {
		return req.NewNotFoundError("Account was not found: %v", p.AccountId)
	}

	cur, err = db.GetCurrencyById(tx, acc.CurrencyId)
	if err != nil {
		c.Log.Error("[Add|Update]Profile: Could not retrieve currency", "error", err.Error())
		return err
	}

	pbm := messaging.PortfolioMessage{
		Portfolio: p,
		Account  : acc,
		Currency : cur,
	}

	err = msg.SendMessage(msg.ExInventory, msg.SourcePortfolio, msgType, &pbm, tx)

	if err != nil {
		c.Log.Error("[Add|Update]Portfolio: Could not publish the update message", "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================
