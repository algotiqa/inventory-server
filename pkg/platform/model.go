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
	SupportsAccount       bool   `json:"supportsAccount"`
}

//=============================================================================

type AccountList struct {
	Offset   int       `json:"offset"`
	Limit    int       `json:"limit"`
	Overflow bool      `json:"overflow"`
	Result   []Account `json:"result"`
}

//=============================================================================

type AccountType string

const (
	AccountTypeFutures AccountType = "FE"
)

//-----------------------------------------------------------------------------

type Account struct {
	Code                 string      `json:"code"`
	Type                 AccountType `json:"type"`
	CurrencyCode         string      `json:"currencyCode"`
	CashBalance          float64     `json:"cashBalance"`
	Equity               float64     `json:"equity"`
	RealizedProfitLoss   float64     `json:"realizedProfitLoss"`
	UnrealizedProfitLoss float64     `json:"unrealizedProfitLoss"`
	OpenOrderMargin      float64     `json:"openOrderMargin"`
	InitialMargin        float64     `json:"initialMargin"`
	MaintenanceMargin    float64     `json:"maintenanceMargin"`
	StatusMessage        string      `json:"statusMessage"`
}

//=============================================================================
