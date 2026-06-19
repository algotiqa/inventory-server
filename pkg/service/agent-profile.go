//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
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

func getAgentProfiles(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetAgentProfiles(tx, c, filter, offset, limit)

			if err != nil {
				return err
			}

			return c.ReturnList(list, offset, limit, len(*list))
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getExternalRefs(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			res, terr := business.GetExternalRefs(tx, c, id)
			if terr != nil {
				return terr
			}

			return c.ReturnObject(res)
		})
	}
	c.ReturnError(err)
}

//=============================================================================
