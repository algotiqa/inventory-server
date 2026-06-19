//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
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

func getBrokerProducts(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		var details bool
		details, err = c.GetParamAsBool("details", false)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				list, terr := business.GetBrokerProducts(tx, c, filter, offset, limit, details)

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

func getBrokerProductById(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			pb, err := business.GetBrokerProductById(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(&pb)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func addBrokerProduct(c *auth.Context) {
	var pds business.BrokerProductSpec
	err := c.BindParamsFromBody(&pds)

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err := business.AddBrokerProduct(tx, c, &pds)

			if err != nil {
				return err
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func updateBrokerProduct(c *auth.Context) {
	var pds business.BrokerProductSpec
	err := c.BindParamsFromBody(&pds)

	if err == nil {
		id, err := c.GetIdFromUrl()

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				ts, err := business.UpdateBrokerProduct(tx, c, id, &pds)

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

func deleteBrokerProduct(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err := business.DeleteBrokerProduct(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================
