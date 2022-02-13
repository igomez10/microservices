package coindesk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const API_URL = "http://api.coindesk.com/v1/bpi/currentprice.json"

var httpClient = http.Client{Timeout: 10 * time.Second}

type CoindDeskProvider struct {
}

type CoinDeskResponse struct {
	Time struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
		Updateduk  string    `json:"updateduk"`
	} `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        struct {
		USD struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
		GBP struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"GBP"`
		EUR struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"EUR"`
	} `json:"bpi"`
}

func (c CoindDeskProvider) GetPitcoinPriceInUSD() (float64, error) {
	req, err := http.NewRequest(http.MethodGet, API_URL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %+v", err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %+v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("failed to read response: %+v", err)
	}

	var coinDeskResponse CoinDeskResponse
	if err := json.Unmarshal(body, &coinDeskResponse); err != nil {
		return 0, fmt.Errorf("failed to parse response into a type 'CoinDeskResponse' request: %+v", err)
	}

	return coinDeskResponse.Bpi.USD.RateFloat, nil
}
