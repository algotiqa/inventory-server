//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package service

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/inventory-server/pkg/business"
	"gorm.io/gorm"
)

//=============================================================================

func getPortfolios(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		var details bool
		details, err = c.GetParamAsBool("details", false)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				list, terr := business.GetPortfolios(tx, c, filter, offset, limit, details)

				if terr != nil {
					return terr
				}

				return c.ReturnList(list, offset, limit, len(*list))
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func getPortfolioById(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			pb, err := business.GetPortfolioById(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(&pb)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func addPortfolio(c *auth.Context) {
	var ps business.PortfolioSpec
	err := c.BindParamsFromBody(&ps)

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err := business.AddPortfolio(tx, c, &ps)

			if err != nil {
				return err
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func updatePortfolio(c *auth.Context) {
	var ps business.PortfolioSpec
	err := c.BindParamsFromBody(&ps)

	if err == nil {
		id, err := c.GetIdFromUrl()

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				ts, err := business.UpdatePortfolio(tx, c, id, &ps)

				if err != nil {
					return err
				}

				return c.ReturnObject(ts)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func deletePortfolio(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err := business.DeletePortfolio(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================
