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

package service

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/inventory-server/pkg/business"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func getConnections(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
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
		err = db.RunInTransaction(func(tx *gorm.DB) error {
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
		err = db.RunInTransaction(func(tx *gorm.DB) error {
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
			err = db.RunInTransaction(func(tx *gorm.DB) error {
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
		err = db.RunInTransaction(func(tx *gorm.DB) error {
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
