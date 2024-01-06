package converter

import (
	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
)

func FromDBUrlToAPIUrl(dbUrl db.Url) server.Url {
	apiUrl := server.Url{
		Alias:     dbUrl.Alias,
		Url:       dbUrl.Url,
		CreatedAt: dbUrl.CreatedAt,
		UpdatedAt: dbUrl.UpdatedAt,
	}

	return apiUrl
}
