//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package process

import (
	"github.com/algotiqa/inventory-server/pkg/app"
	"github.com/algotiqa/inventory-server/pkg/core/process/agentscanner"
	"github.com/algotiqa/inventory-server/pkg/core/process/currencyupdater"
)

//=============================================================================

func Init(cfg *app.Config) {
	agentscanner.Init(cfg)
	currencyupdater.Init(cfg)
}

//=============================================================================
