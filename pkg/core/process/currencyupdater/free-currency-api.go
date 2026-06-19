//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package currencyupdater

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/algotiqa/types"
)

//=============================================================================

const (
	Historical = "historical"
)

//=============================================================================

type FreeCurrencyClient struct {
	baseUrl string
	apiKey  string
	client  *http.Client
}

//=============================================================================

func NewFreeCurrencyClient(baseUrl, apiKey string) *FreeCurrencyClient {
	return &FreeCurrencyClient{
		baseUrl: baseUrl,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func (f *FreeCurrencyClient) GetHistoricalValues(date types.Date, baseCurrency, currencyList string) (*HistoricalResponse, error) {
	params := map[string]string{}
	params["date"] = date.String()
	params["base_currency"] = baseCurrency
	params["currencies"] = currencyList

	res, err := f.callAPI(Historical, params)
	if err != nil {
		return nil, err
	}

	var output map[string]interface{}
	err = json.Unmarshal(res, &output)
	if err != nil {
		return nil, err
	}

	return convertHistoricalResponse(output), nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func (f *FreeCurrencyClient) callAPI(service string, params map[string]string) ([]byte, error) {
	url := f.baseUrl + "/" + service + "?" + mapToQueryParams(params)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", f.apiKey)

	response, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

//=============================================================================

func mapToQueryParams(params map[string]string) string {
	var sb strings.Builder

	for k, v := range params {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString("&")
	}

	return strings.TrimRight(sb.String(), "&")
}

//=============================================================================

func convertHistoricalResponse(output map[string]interface{}) *HistoricalResponse {
	res := &HistoricalResponse{
		Currencies: make(map[string]float64),
	}

	val, ok := output["data"]
	if ok {
		mapVal := val.(map[string]interface{})
		for k, v := range mapVal {
			res.Date = k
			mapCur := v.(map[string]interface{})
			for code, value := range mapCur {
				res.Currencies[code] = value.(float64)
			}
		}
	}
	return res
}

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type HistoricalResponse struct {
	Date       string             `json:"date"`
	Currencies map[string]float64 `json:"currencies"`
}

//=============================================================================
