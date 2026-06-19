//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/types"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSessions(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]*TradingSession, error) {
	list, err := db.GetTradingSessions(tx, c.Session.Username)
	if err != nil {
		return nil, err
	}

	res,errs := rebuildSessions(c, list)
	if errs != nil {
		return nil, req.NewServerError("Unable to rebuild the sessions")
	}

	return &res, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func rebuildSessions(c *auth.Context, list *[]db.TradingSession) ([]*TradingSession, error) {
	var result []*TradingSession

	for _, dbTs := range *list {
		sess,err := types.NewTradingSession(dbTs.Session)
		if err != nil {
			c.Log.Error("rebuildSessions: Invalid session config", "error", err.Error(), "config", dbTs.Session)
			return nil,err
		}

		busTs := &TradingSession{
			Common:   dbTs.Common,
			Name:     dbTs.Name,
			Username: dbTs.Username,
			Session:  sess,
		}

		result = append(result, busTs)
	}

	return result,nil
}

//=============================================================================
