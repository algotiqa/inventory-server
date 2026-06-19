//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/inventory-server/pkg/db"
)

//=============================================================================

type DataProductSpec struct {
	ConnectionId    uint             `json:"connectionId"   binding:"required"`
	ExchangeId      uint             `json:"exchangeId"     binding:"required"`
	Symbol          string           `json:"symbol"         binding:"required"`
	Name            string           `json:"name"           binding:"required"`
	MarketType      string           `json:"marketType"     binding:"required"`
	ProductType     string           `json:"productType"    binding:"required"`
	Months          string           `json:"months"`
	RolloverTrigger db.DPRollTrigger `json:"rolloverTrigger"`
	SessionId       uint             `json:"sessionId"      binding:"required"`
}

//=============================================================================

func (s *DataProductSpec) validateForAdd() error {
	//TODO: validate rollover trigger

	return nil
}

//=============================================================================

func (s *DataProductSpec) validateForUpdate() error {
	return nil
}

//=============================================================================
