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
