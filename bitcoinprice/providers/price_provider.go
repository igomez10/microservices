package providers

type PriceProvider interface {
	GetPitcoinPriceInUSD() (float64, error)
}
