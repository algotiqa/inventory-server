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

func GetPortfolios(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]PortfolioFull, error) {
	var list []PortfolioFull
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetPortfoliosFull(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]PortfolioFull, error) {
	var list []PortfolioFull
	res := tx.Model(&Portfolio{}).Select("portfolio.*, " +
		"currency.code as currency_code, account.code as account_code, account.name as account_name").
		Joins("LEFT JOIN account  ON portfolio.account_id = account.id").
		Joins("LEFT JOIN currency ON account.currency_id  = currency.id").
		Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetPortfolioById(tx *gorm.DB, id uint) (*Portfolio, error) {
	var list []Portfolio
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

func AddPortfolio(tx *gorm.DB, p *Portfolio) error {
	return tx.Create(p).Error
}

//=============================================================================

func UpdatePortfolio(tx *gorm.DB, p *Portfolio) error {
	return tx.Save(p).Error
}

//=============================================================================

func DeletePortfolio(tx *gorm.DB, id uint) error {
	return tx.Delete(&Portfolio{}, id).Error
}

//=============================================================================
