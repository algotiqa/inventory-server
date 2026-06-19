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

func GetDataProducts(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]DataProductFull, error) {
	var list []DataProductFull
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetDataProductsFull(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]DataProductFull, error) {
	var list []DataProductFull

	res := tx.Model(&DataProduct{}).Select("data_product.*, " +
		"connection.code as connection_code, connection.name as connection_name, connection.system_code as system_code, " +
		"exchange.code as exchange_code").
		Joins("LEFT JOIN connection ON data_product.connection_id = connection.id").
		Joins("LEFT JOIN exchange   ON data_product.exchange_id   = exchange.id").
		Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetDataProductById(tx *gorm.DB, id uint) (*DataProduct, error) {
	var list []DataProduct
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

func AddDataProduct(tx *gorm.DB, ts *DataProduct) error {
	return tx.Create(ts).Error
}

//=============================================================================

func UpdateDataProduct(tx *gorm.DB, ts *DataProduct) error {
	return tx.Save(ts).Error
}

//=============================================================================

func DeleteDataProduct(tx *gorm.DB, id uint) error {
	return tx.Delete(&DataProduct{}, id).Error
}

//=============================================================================
