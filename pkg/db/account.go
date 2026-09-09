//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
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

func GetAccounts(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]AccountFull, error) {
	var list []AccountFull
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetAccountsFull(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]AccountFull, error) {
	var list []AccountFull
	res := tx.Model(&Account{}).Select("account.*, " +
		"currency.code as currency_code, " +
		"connection.code as connection_code, connection.name as connection_name, connection.system_code as system_code").
		Joins("LEFT JOIN connection ON account.connection_id = connection.id").
		Joins("LEFT JOIN currency   ON account.currency_id   = currency.id").
		Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetAccountById(tx *gorm.DB, id uint) (*Account, error) {
	var list []Account
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

func AddAccount(tx *gorm.DB, a *Account) error {
	return tx.Create(a).Error
}

//=============================================================================

func UpdateAccount(tx *gorm.DB, a *Account) error {
	return tx.Save(a).Error
}

//=============================================================================

func DeleteAccount(tx *gorm.DB, id uint) error {
	return tx.Delete(&Account{}, id).Error
}

//=============================================================================
