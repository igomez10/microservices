package grifts

import (
	"github.com/igomez10/microservices/bitcoinprice/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
