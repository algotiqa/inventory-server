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

func GetTradingSessions(tx *gorm.DB, username string) (*[]TradingSession, error) {
	var list []TradingSession
	res := tx.Find(&list, "username = ? or username is null order by id", username)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSessionById(tx *gorm.DB, id uint) (*TradingSession, error) {
	var list []TradingSession
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

func AddTradingSession(tx *gorm.DB, s *TradingSession) error {
	return tx.Create(s).Error
}

//=============================================================================
