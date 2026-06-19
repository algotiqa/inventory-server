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

func getConnections(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetConnections(tx, c, filter, offset, limit)

			if err != nil {
				return err
			}

			return c.ReturnList(list, offset, limit, len(*list))
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getConnectionById(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			conn, err := business.GetConnectionById(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(conn)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func addConnection(c *auth.Context) {
	var cs business.ConnectionSpec
	err := c.BindParamsFromBody(&cs)

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			conn, err := business.AddConnection(tx, c, &cs)

			if err != nil {
				return err
			}

			return c.ReturnObject(conn)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func updateConnection(c *auth.Context) {
	var cs business.ConnectionSpec
	err := c.BindParamsFromBody(&cs)

	if err == nil {
		var id uint
		id, err = c.GetIdFromUrl()

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				ts, err := business.UpdateConnection(tx, c, id, &cs)

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

func deleteConnection(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err := business.DeleteConnection(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================
