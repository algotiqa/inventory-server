//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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

package currencyupdater

import (
	"log/slog"
	"strings"
	"time"

	"github.com/algotiqa/core/datatype"
	"github.com/algotiqa/inventory-server/pkg/app"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

const (
	BaseCurrency = "USD"
)

var baseUrl string
var apiKey string

//=============================================================================

func Init(cfg *app.Config) {
	baseUrl = cfg.Provider.Currency.BaseUrl
	apiKey = cfg.Provider.Currency.ApiKey

	ticker := time.NewTicker(45 * time.Minute)

	go func() {
		time.Sleep(10 * time.Second)
		run()

		for range ticker.C {
			run()
		}
	}()
}

//=============================================================================

func run() {
	slog.Info("CurrencyUpdater: Starting sync process")

	currencies, err := getCurrencies()
	if err == nil {
		//--- This is BaseCurrency (USD)
		cur := currencies[0]

		var history []*db.CurrencyHistory

		if cur.LastDate.IsNil() {
			history, err = latestUpdate(currencies, datatype.Today(time.UTC).AddDays(-1))
		} else if newLatestDay(cur) {
			history, err = latestUpdate(currencies, cur.LastDate.AddDays(1))
		} else if !cur.HistoryEnded {
			history, err = dateUpdate(currencies, cur.FirstDate.AddDays(-1))
		}

		if err == nil {
			err = saveCurrenciesAndHistory(currencies, history)
		}
	}

	slog.Info("CurrencyUpdater: Ending sync process")
}

//=============================================================================

func getCurrencies() ([]*db.Currency, error) {
	var list *[]db.Currency

	err := db.RunInTransaction(func(tx *gorm.DB) error {
		var err error
		list, err = db.GetCurrencies(tx)
		return err
	})

	if err != nil {
		slog.Error("getCurrencies: Cannot retrieve currencies", "error", err)
	}

	var res []*db.Currency
	for _, cur := range *list {
		res = append(res, &cur)
	}

	return res, err
}

//=============================================================================

func newLatestDay(cur *db.Currency) bool {
	today := datatype.Today(time.UTC)

	return cur.LastDate.AddDays(1) < today
}

//=============================================================================

func latestUpdate(currencies []*db.Currency, date datatype.IntDate) ([]*db.CurrencyHistory, error) {
	fcc := NewFreeCurrencyClient(baseUrl, apiKey)
	res, err := fcc.GetHistoricalValues(date, BaseCurrency, toList(currencies))
	if err != nil {
		slog.Error("latestUpdate: Cannot retrieve currencies from provider", "error", err, "date", date)
		return nil, err
	}

	var history []*db.CurrencyHistory

	for _, cur := range currencies {
		if cur.FirstDate.IsNil() {
			cur.FirstDate = date
		}
		cur.LastDate = date
		value, ok := res.Currencies[cur.Code]
		//--- Skipping BaseCurrency
		if ok {
			cur.LastValue = value

			ci := &db.CurrencyHistory{
				CurrencyId: cur.Id,
				Date:       date,
				Value:      value,
			}

			history = append(history, ci)
		}
	}

	return history, nil
}

//=============================================================================

func dateUpdate(currencies []*db.Currency, date datatype.IntDate) ([]*db.CurrencyHistory, error) {
	fcc := NewFreeCurrencyClient(baseUrl, apiKey)
	res, err := fcc.GetHistoricalValues(date, BaseCurrency, toList(currencies))
	if err != nil {
		slog.Error("dateUpdate: Cannot retrieve currencies from provider", "error", err, "date", date)
		return nil, err
	}

	var history []*db.CurrencyHistory

	for _, cur := range currencies {
		cur.FirstDate = date
		cur.HistoryEnded = date == 20000101

		value, ok := res.Currencies[cur.Code]
		//--- Skipping BaseCurrency
		if ok {
			ci := &db.CurrencyHistory{
				CurrencyId: cur.Id,
				Date:       date,
				Value:      value,
			}

			history = append(history, ci)
		}
	}

	return history, nil
}

//=============================================================================

func saveCurrenciesAndHistory(currencies []*db.Currency, history []*db.CurrencyHistory) error {
	return db.RunInTransaction(func(tx *gorm.DB) error {
		for _, cur := range currencies {
			err := db.UpdateCurrency(tx, cur)
			if err != nil {
				slog.Error("saveCurrenciesAndHistory: Cannot update currencies", "error", err)
				return err
			}
		}

		for _, ci := range history {
			err := db.AddCurrencyHistory(tx, ci)
			if err != nil {
				slog.Error("saveCurrenciesAndHistory: Cannot save currency history", "error", err)
				return err
			}
		}

		return nil
	})
}

//=============================================================================

func toList(list []*db.Currency) string {
	var res strings.Builder

	isFirst := true

	for _, cur := range list {
		if cur.Code != BaseCurrency {
			if !isFirst {
				res.WriteString(",")
			}

			res.WriteString(cur.Code)
			isFirst = false
		}
	}

	return res.String()
}

//=============================================================================
