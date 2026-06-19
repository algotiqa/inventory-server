//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
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

func getExchanges(c *auth.Context) {
	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		list, err := business.GetExchanges(tx)

		if err != nil {
			return err
		}

		return c.ReturnList(list, 0, len(*list), len(*list))
	})

	c.ReturnError(err)
}

//=============================================================================
