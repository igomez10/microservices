package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

type BitcoinPriceHandlerResponse struct {
	Currency string `json:"currency,omitempty"`
	Value    int    `json:"value,omitempty"`
}

// BitcoinPriceHandler is a handler to serve up
// a the price of bitcoin.
func BitcoinPriceHandler(c buffalo.Context) error {
	body := BitcoinPriceHandlerResponse{Currency: "usd", Value: 3000}
	res := c.Render(http.StatusOK, r.JSON(body))
	return res
}
