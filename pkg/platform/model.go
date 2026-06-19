//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================


package platform

//=============================================================================

type SystemList struct {
	Offset   int      `json:"offset"`
	Limit    int      `json:"limit"`
	Overflow bool     `json:"overflow"`
	Result   []System `json:"result"`
}

//=============================================================================

type System struct {
	Code                  string `json:"code"`
	Name                  string `json:"name"`
	SupportsData          bool   `json:"supportsData"`
	SupportsBroker        bool   `json:"supportsBroker"`
	SupportsMultipleData  bool   `json:"supportsMultipleData"`
	SupportsInventory     bool   `json:"supportsInventory"`
}

//=============================================================================
