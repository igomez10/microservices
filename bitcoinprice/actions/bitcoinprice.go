package actions

import (
	"bitcoinprice/providers"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/rs/zerolog/log"
)

type BitcoinPriceHandlerResponse struct {
	Currency string `json:"currency,omitempty"`
	Value    int    `json:"value,omitempty"`
}

var PriceProvider providers.PriceProvider

// BitcoinPriceHandler is a handler to serve up
// a the price of bitcoin.
func BitcoinPriceHandler(c buffalo.Context) error {
	btcPrice, err := PriceProvider.GetPitcoinPriceInUSD()
	if err != nil {
		log.Err(err).Msg("failed to retrieve bitcoin price")
		response := map[string]string{"error": "failed to retrieve bitcoin price"}
		if errRender := c.Render(http.StatusInternalServerError, r.JSON(response)); errRender != nil {
			log.Err(err).Msg("failed to send response")
			// TODO add metrics for failed to send response
		}

		// TODO add metrics for 500
		return err
	}

	btcPriceInCents := btcPrice * 100 // Always use price in cents
	response := BitcoinPriceHandlerResponse{Currency: "USD", Value: int(btcPriceInCents * 100)}
	res := c.Render(http.StatusOK, r.JSON(response))
	return res
}
