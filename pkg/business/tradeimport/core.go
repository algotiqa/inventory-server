//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package tradeimport

import (
	"errors"

	"github.com/algotiqa/inventory-server/pkg/core/process/agentscanner"
)

//=============================================================================

func ImportTrades(spec *Spec, data []byte) (*agentscanner.TradingSystem,error) {
	if isTradeListFile(data) {
		return importTradeList(spec, data)
	}

	return nil,errors.New("Unknown trade file format")
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func isTradeListFile(data []byte) bool {
	if len(data) >9 {
		prefix := string(data[:5])
		if prefix == "INFO|" {
			return true
		}
	}

	return false
}

//=============================================================================
