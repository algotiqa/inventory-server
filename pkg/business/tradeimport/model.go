//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package tradeimport

//=============================================================================

type Spec struct {
	WipeData bool `json:"wipeData"`
}

//=============================================================================

type Response struct {
	Trades       int `json:"trades"`
	Bars         int `json:"bars"`
	DailyReturns int `json:"dailyReturns"`
}

//=============================================================================
