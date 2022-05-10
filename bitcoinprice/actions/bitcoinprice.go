package actions

import (
	"bitcoinprice/providers"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/rs/zerolog/log"
)

type BitcoinPriceHandlerResponse struct {
	Currency string `json:"currency,omitempty"`
	Value    int    `json:"value,omitempty"`
}

type CacheEntry struct {
	createdAt time.Time
	response  BitcoinPriceHandlerResponse
}

var PriceProvider providers.PriceProvider

var mostRecentResponse *CacheEntry

// BitcoinPriceHandler is a handler to serve up
// a the price of bitcoin.
func BitcoinPriceHandler(c buffalo.Context) error {

	// Check cache if exist and is valid
	if mostRecentResponse != nil {
		elapsedTime := time.Now().Sub(mostRecentResponse.createdAt)
		if elapsedTime < time.Duration(10*time.Second) {
			log.Debug().Msg("Serving bitcoin price from cache")
			if err := c.Render(http.StatusOK, r.JSON(mostRecentResponse.response)); err != nil {
				log.Error().Err(err).Msg("Failed to render response")
				return err
			}
			return nil
		}
	}

	//  Issue request to price provider
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

	if err := c.Render(http.StatusOK, r.JSON(response)); err != nil {
		log.Error().Err(err).Msg("Failed to render response")
		return err
	}

	mostRecentResponse = &CacheEntry{
		createdAt: time.Now(),
		response:  response,
	}

	return nil
}
