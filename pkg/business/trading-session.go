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
	"encoding/json"
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/sick-engine/session"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSessions(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]TradingSession, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	list, err := db.GetTradingSessions(tx, filter, offset, limit)

	if err != nil {
		return nil, err
	}

	var res []TradingSession

	for _, dbTs := range *list {
		var sickTs session.TradingSession

		err = json.Unmarshal([]byte(dbTs.Config), &sickTs)
		if err != nil {
			c.Log.Error("GetTradingSessions: Invalid session config", "error", err.Error())
			return nil, err
		}

		busTs := TradingSession{
			Common:   dbTs.Common,
			Name:     dbTs.Name,
			Username: dbTs.Username,
			Session:  &sickTs,
		}

		res = append(res, busTs)
	}

	return &res, nil
}

//=============================================================================
