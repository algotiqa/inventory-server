//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================


package db

import (
	"github.com/algotiqa/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func GetCurrencies(tx *gorm.DB) (*[]Currency, error) {
	var list []Currency
	res := tx.Find(&list).Order("code")

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetCurrencyById(tx *gorm.DB, id uint) (*Currency, error) {
	var list []Currency
	res := tx.Find(&list, id)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	return nil, nil
}

//=============================================================================

func GetCurrenciesAsMap(tx *gorm.DB) (map[uint]*Currency, error) {
	list,err := GetCurrencies(tx)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*Currency)
	for _,c := range *list {
		result[c.Id] = &c
	}

	return result,nil
}

//=============================================================================

func UpdateCurrency(tx *gorm.DB, c *Currency) error {
	return tx.Save(c).Error
}

//=============================================================================

func AddCurrencyHistory(tx *gorm.DB, ci *CurrencyHistory) error {
	return tx.Create(ci).Error
}

//=============================================================================
