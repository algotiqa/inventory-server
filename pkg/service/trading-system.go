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
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/business"
	"github.com/algotiqa/inventory-server/pkg/business/importexport"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func getTradingSystems(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		var details bool
		details, err = c.GetParamAsBool("details", false)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				list, terr := business.GetTradingSystems(tx, c, filter, offset, limit, details)

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

func addTradingSystem(c *auth.Context) {
	var tss business.TradingSystemSpec
	err := c.BindParamsFromBody(&tss)

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, terr := business.AddTradingSystem(tx, c, &tss)

			if terr != nil {
				return terr
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func updateTradingSystem(c *auth.Context) {
	var tss business.TradingSystemSpec
	err := c.BindParamsFromBody(&tss)

	if err == nil {
		var id uint
		id, err = c.GetIdFromUrl()

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				var ts *db.TradingSystem
				ts, err = business.UpdateTradingSystem(tx, c, id, &tss)

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

func deleteTradingSystem(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err := business.DeleteTradingSystem(tx, c, id)

			if err != nil {
				return err
			}

			return c.ReturnObject(ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func finalizeTradingSystem(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, terr := business.FinalizeTradingSystem(tx, c, id)

			if terr != nil {
				return terr
			}

			return c.ReturnObject(ts)
		})
	}
	c.ReturnError(err)
}

//=============================================================================

func reloadTradesFromAgent(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			res, terr := business.ReloadTradesFromAgent(tx, c, id)

			if terr != nil {
				return terr
			}

			return c.ReturnObject(res)
		})
	}
	c.ReturnError(err)
}

//=============================================================================

func exportTradingSystems(c *auth.Context) {
	ids,err := c.GetIdsFromUrl()
	if err == nil {
		if len(ids) == 0 {
			err = req.NewBadRequestError("Parameter 'id' is missing or empty")
		} else {
			var data *importexport.TradingSystemsData
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				res, terr := business.CollectTradingSystemsData(tx, c, ids)
				data = res
				return terr
			})

			if err == nil {
				var res []byte
				res,err = business.ExportTradingSystems(c, data)
				if err == nil {
					_=c.ReturnData("application/zip", res)
					return
				}
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func createImportOverview(c *auth.Context) {
	reader, err := c.Gin.Request.MultipartReader()
	if err == nil {
		var part *multipart.Part

		if part, err = reader.NextPart(); err != io.EOF {
			var spec *business.ImportOverviewSpec
			spec, err = retrieveImportOverviewSpec(part)

			if err == nil {
				if part, err = reader.NextPart(); err != io.EOF {
					var data []byte
					data, err = getImportPackage(part)
					_ = part.Close()

					if err == nil {
						var res *importexport.ImportOverviewResponse
						err = dbms.RunInTransaction(func(tx *gorm.DB) error {
							out,errx :=business.CreateImportOverview(tx, c, spec, data)
							res = out
							return errx
						})

						if err == nil {
							_ = c.ReturnObject(res)
							return
						}
					}
				}
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func executeImportPlan(c *auth.Context) {
	reader, err := c.Gin.Request.MultipartReader()
	if err == nil {
		var part *multipart.Part

		if part, err = reader.NextPart(); err != io.EOF {
			var spec *business.ImportExecutionSpec
			spec, err = retrieveImportExecutionSpec(part)

			if err == nil {
				if part, err = reader.NextPart(); err != io.EOF {
					var data []byte
					data, err = getImportPackage(part)
					_ = part.Close()

					if err == nil {
						err = dbms.RunInTransaction(func(tx *gorm.DB) error {
							return business.ExecuteImportPlan(tx, c, spec, data)
						})

						if err == nil {
							_ = c.ReturnObject(true)
							return
						}
					}
				}
			}
		}
	}

	c.ReturnError(err)
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func retrieveImportOverviewSpec(part *multipart.Part) (*business.ImportOverviewSpec, error) {
	data, err := io.ReadAll(part)

	if err == nil {
		var spec business.ImportOverviewSpec

		err = json.Unmarshal(data, &spec)

		if err == nil {
			err = part.Close()

			if err == nil {
				return &spec, nil
			}
		}
	}

	return nil, err
}

//=============================================================================

func retrieveImportExecutionSpec(part *multipart.Part) (*business.ImportExecutionSpec, error) {
	data, err := io.ReadAll(part)

	if err == nil {
		var spec business.ImportExecutionSpec

		err = json.Unmarshal(data, &spec)

		if err == nil {
			err = part.Close()

			if err == nil {
				return &spec, nil
			}
		}
	}

	return nil, err
}

//=============================================================================

func getImportPackage(part io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, part)

	if err != nil {
		return nil,err
	}

	return buf.Bytes(), nil
}

//=============================================================================
