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

func GetConnections(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]Connection, error) {
	var list []Connection
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetConnectionById(tx *gorm.DB, id uint) (*Connection, error) {
	var list []Connection
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

func GetConnectionByCode(tx *gorm.DB, user, code string) (*Connection, error) {
	var list []Connection
	res := tx.Where("username = ? AND code = ?", user, code).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	return nil, nil
}

//=============================================================================

func AddConnection(tx *gorm.DB, conn *Connection) error {
	return tx.Create(conn).Error
}

//=============================================================================

func UpdateConnection(tx *gorm.DB, conn *Connection) error {
	return tx.Save(conn).Error
}

//=============================================================================

func DeleteConnection(tx *gorm.DB, id uint) error {
	return tx.Delete(&Connection{}, id).Error
}

//=============================================================================

func DisconnectAll(tx *gorm.DB) error {
	return tx.Model(&Connection{}).
		Where("supports_multiple_data = false").
		Update("connected", false).Error
}

//=============================================================================

func SetConnectionStatus(tx *gorm.DB, user, code string, flag bool) error {
	return tx.Model(&Connection{}).
		Where("username = ? AND code = ?", user, code).
		Update("connected", flag).Error
}

//=============================================================================
