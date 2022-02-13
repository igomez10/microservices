package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func HealthHandler(c buffalo.Context) error {
	response := map[string]bool{"alive": true, "ready": true}
	res := c.Render(http.StatusOK, r.JSON(response))
	return res
}
