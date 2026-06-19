//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
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

func GetExchanges(tx *gorm.DB) (*[]Exchange, error) {
	var list []Exchange
	res := tx.Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetExchangeById(tx *gorm.DB, id uint) (*Exchange, error) {
	var list []Exchange
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

func GetExchangesAsMap(tx *gorm.DB) (map[uint]*Exchange, error) {
	list,err := GetExchanges(tx)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*Exchange)
	for _,e := range *list {
		result[e.Id] = &e
	}

	return result,nil
}

//=============================================================================
