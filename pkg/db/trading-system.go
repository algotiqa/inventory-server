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

func GetTradingSystems(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]TradingSystemFull, error) {
	var list []TradingSystemFull
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemById(tx *gorm.DB, id uint) (*TradingSystem, error) {
	var list []TradingSystem
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

func GetTradingSystemsById(tx *gorm.DB, username string, ids []uint) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "username = ? and id in ?", username, ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemsFull(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]TradingSystemFull, error) {
	var list []TradingSystemFull

	res := tx.Model(&TradingSystem{}).Select("trading_system.*, " +
		"data_product.symbol as data_symbol, " +
		"broker_product.symbol as broker_symbol, " +
		"trading_session.name as trading_session ").
		Joins("LEFT JOIN data_product    ON trading_system.data_product_id   = data_product.id").
		Joins("LEFT JOIN broker_product  ON trading_system.broker_product_id = broker_product.id").
		Joins("LEFT JOIN trading_session ON trading_system.trading_session_id= trading_session.id").
		Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemByExtRef(tx *gorm.DB, username string, externalRef string) (*TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "external_ref = ? and username = ?", externalRef, username)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	return nil, nil
}

//=============================================================================

func AddTradingSystem(tx *gorm.DB, ts *TradingSystem) error {
	return tx.Create(ts).Error
}

//=============================================================================

func UpdateTradingSystem(tx *gorm.DB, ts *TradingSystem) error {
	return tx.Save(ts).Error
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, id uint) error {
	return tx.Delete(&TradingSystem{}, id).Error
}

//=============================================================================
