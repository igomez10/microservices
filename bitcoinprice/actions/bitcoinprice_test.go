package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (as *ActionSuite) Test_BitcoinPriceHandlerWithBuffalo() {

	res := as.JSON("/api/v1/btc").Get()

	as.Equal(http.StatusOK, res.Code)
	btcResponse := BitcoinPriceHandlerResponse{}
	if err := json.Unmarshal(res.Body.Bytes(), &btcResponse); err != nil {
		as.Errorf(err, "Failed to unmarshal response as BitcoinPriceHandlerResponse")
	}

	if currencyLength := len(btcResponse.Currency); currencyLength != 3 {
		msg := fmt.Sprintf("Unexpected currency %q with length: %d should be 3 characters long", btcResponse.Currency, currencyLength)
		as.T().Error(msg)
	}

	if btcResponse.Value <= 0 {
		msg := fmt.Sprintf("Currency value cannot be negative or equal to 0 but was %d", btcResponse.Value)
		as.T().Error(msg)
	}
}
