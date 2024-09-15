package url

import (
	"context"
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/db"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	urlClient "github.com/igomez10/microservices/urlshortener/generated/clients/go/client"
)

// implements the URLApiServicer interface
// s *URLApiService openapi.URLApiServicer
type URLApiService struct {
	// legacyDB is the legacy database, now we use the urlshortener microservice
	// DEPRECATED
	DB db.Querier
	// DEPRECATED
	DBConn db.DBTX

	// NEW URLSHORTENER SERVICE
	// urlClient is the Client for the urlshortener service
	Client *urlClient.APIClient

	// feature flags
	UseURLShortenerService bool
}

type URLApiServiceConfig struct {
	DB                     db.Querier
	DBConn                 db.DBTX
	Client                 *urlClient.APIClient
	UseURLShortenerService bool
}

// CreateUrl creates a new url
func (s *URLApiService) CreateUrl(ctx context.Context, newURL openapi.Url) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateUrl")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	var openapiURL openapi.Url
	if s.UseURLShortenerService {
		// create the url
		{
			newURLRequest := urlClient.NewURL(newURL.Url, newURL.Alias)
			u, createRes, err := s.Client.URLAPI.CreateUrl(ctx).
				URL(*newURLRequest).
				Execute()
			if err != nil {
				log.Error().Err(err).Msgf("error creating url %q with alias %q", newURL.Url, newURL.Alias)
				return openapi.ImplResponse{
					Code: http.StatusInternalServerError,
					Body: openapi.Error{
						Message: "error creating url",
						Code:    http.StatusInternalServerError,
					},
				}, err
			}

			switch createRes.StatusCode {
			case http.StatusConflict:
				log.Error().Err(err).Msgf("url with alias %q already exists", newURL.Alias)
				return openapi.ImplResponse{
					Code: http.StatusConflict,
					Body: openapi.Error{
						Message: fmt.Sprintf("url with alias %q already exists", newURL.Alias),
						Code:    http.StatusConflict,
					},
				}, nil
			}

			openapiURL = openapi.Url{
				Url:       u.Url,
				Alias:     u.Alias,
				CreatedAt: *u.CreatedAt,
				UpdatedAt: *u.UpdatedAt,
				DeletedAt: *u.DeletedAt,
			}
		}

	} else {
		// validate we dont have a url with the same alias
		if _, err := s.DB.GetURLFromAlias(ctx, s.DBConn, newURL.Alias); err == nil {
			log.Error().Err(err).Msg("url with alias already exists")
			return openapi.ImplResponse{
				Code: http.StatusConflict,
				Body: openapi.Error{
					Message: "url with alias already exists",
					Code:    http.StatusConflict,
				},
			}, nil
		}

		newURLParams := db.CreateURLParams{
			Alias: newURL.Alias,
			Url:   newURL.Url,
		}
		res, err := s.DB.CreateURL(ctx, s.DBConn, newURLParams)
		if err != nil {
			log.Error().Err(err).Msg("error creating url")
			return openapi.ImplResponse{}, err
		}
		openapiURL = converter.FromDBUrlToAPIUrl(res)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: openapiURL,
	}, nil

}

// DeleteUrl deletes a url
func (s *URLApiService) DeleteUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "DeleteUrl")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	if s.UseURLShortenerService {
		// use the urlshortener service
		res, err := s.Client.URLAPI.DeleteUrl(ctx, alias).Execute()
		if err != nil || res.StatusCode != http.StatusOK {
			log.Error().Err(err).Msg("error deleting url")
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Message: "error deleting url",
					Code:    http.StatusInternalServerError,
				},
			}, err
		}
	} else {
		if err := s.DB.DeleteURL(ctx, s.DBConn, alias); err != nil {
			log.Error().Err(err).Msg("error deleting url")
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Message: "error deleting url",
					Code:    http.StatusInternalServerError,
				},
			}, err
		}
	}
	return openapi.ImplResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrl")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)

	var shortURL string
	if s.UseURLShortenerService {
		// use the urlshortener service
		u, res, err := s.Client.URLAPI.GetUrlData(ctx, alias).Execute()
		if err != nil || res.StatusCode != http.StatusOK {
			switch res.StatusCode {
			case http.StatusNotFound:
				log.Debug().Err(err).Msg("alias does not exist")
				return openapi.ImplResponse{
					Code: http.StatusNotFound,
					Body: openapi.Error{
						Message: "alias does not exist",
						Code:    http.StatusNotFound,
					},
				}, nil
			default:
				log.Error().Err(err).Msg("error getting url from urlshortener service")
				return openapi.ImplResponse{
					Code: http.StatusInternalServerError,
					Body: openapi.Error{
						Message: "error fetching url from downstream service",
						Code:    http.StatusInternalServerError,
					},
				}, err
			}
		}

		shortURL = u.Url
	} else {
		shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
		if err != nil {
			log.Debug().Err(err).Msg("alias does not exist")
			return openapi.ImplResponse{
				Code: http.StatusNotFound,
				Body: openapi.Error{
					Message: "alias does not exist",
					Code:    http.StatusNotFound,
				},
			}, nil
		}
		shortURL = shortedURL.Url
	}

	// add location header for redirect in the response
	res := openapi.ImplResponse{
		Code: http.StatusPermanentRedirect,
		Headers: map[string][]string{
			"Location": {shortURL},
		},
	}

	return res, nil
}

func (s *URLApiService) GetUrlData(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrlData")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)

	var responseUrl openapi.Url
	if s.UseURLShortenerService {
		u, res, err := s.Client.URLAPI.GetUrlData(ctx, alias).Execute()
		if err != nil || res.StatusCode != http.StatusOK {
			log.Error().Err(err).Msg("error getting url from urlshortener service")
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Message: "error fetching url from downstream service",
					Code:    http.StatusInternalServerError,
				},
			}, err
		}

		responseUrl = openapi.Url{
			Alias:     u.Alias,
			Url:       u.Url,
			CreatedAt: *u.CreatedAt,
			UpdatedAt: *u.UpdatedAt,
			DeletedAt: *u.DeletedAt,
		}
	} else {
		// validate we dont have a url with the same alias
		shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
		if err != nil {
			log.Debug().Err(err).Msg("alias does not exist")
			return openapi.ImplResponse{
				Code: http.StatusNotFound,
				Body: openapi.Error{
					Message: "alias does not exist",
					Code:    http.StatusNotFound,
				},
			}, nil
		}

		responseUrl = converter.FromDBUrlToAPIUrl(shortedURL)
	}

	res := openapi.ImplResponse{
		Code: http.StatusOK,
		Body: responseUrl,
	}

	return res, nil
}
