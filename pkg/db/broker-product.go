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

func GetBrokerProducts(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]BrokerProductFull, error) {
	var list []BrokerProductFull
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetBrokerProductsFull(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]BrokerProductFull, error) {
	var list []BrokerProductFull
	res := tx.Model(&BrokerProduct{}).Select("broker_product.*, " +
		"currency.code as currency_code, " +
		"connection.code as connection_code, connection.name as connection_name, connection.system_code as system_code, " +
		"exchange.code as exchange_code").
		Joins("LEFT JOIN connection ON broker_product.connection_id = connection.id").
		Joins("LEFT JOIN exchange   ON broker_product.exchange_id   = exchange.id").
		Joins("LEFT JOIN currency   ON exchange.currency_id         = currency.id").
		Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetBrokerProductById(tx *gorm.DB, id uint) (*BrokerProduct, error) {
	var list []BrokerProduct
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

func AddBrokerProduct(tx *gorm.DB, pb *BrokerProduct) error {
	return tx.Create(pb).Error
}

//=============================================================================

func UpdateBrokerProduct(tx *gorm.DB, pb *BrokerProduct) error {
	return tx.Save(pb).Error
}

//=============================================================================

func DeleteBrokerProduct(tx *gorm.DB, id uint) error {
	return tx.Delete(&BrokerProduct{}, id).Error
}

//=============================================================================
//===
//=== Broker instruments
//===
//=============================================================================

func GetBrokerInstrumentsByBrokerId(tx *gorm.DB, id uint) (*[]BrokerInstrument, error) {
	var list []BrokerInstrument

	filter := map[string]any{}
	filter["broker_product_id"] = id

	res := tx.Where(filter).Order("expiration_date").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================
