//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package currencyupdater

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/algotiqa/core/datatype"
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

func (f *FreeCurrencyClient) GetHistoricalValues(date datatype.IntDate, baseCurrency, currencyList string) (*HistoricalResponse, error) {
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
