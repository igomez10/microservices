package url

import (
	"context"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
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

	// urlClient is the Client for the urlshortener service
	Client *urlClient.APIClient

	// feature flags
	UseURLShortenerService bool
}

func (s *URLApiService) CreateUrl(ctx context.Context, newURL openapi.Url) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

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

	openapiURL := converter.FromDBUrlToAPIUrl(res)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: openapiURL,
	}, nil

}

func (s *URLApiService) DeleteUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	if err := s.DB.DeleteURL(ctx, s.DBConn, alias); err != nil {
		log.Error().Err(err).Msg("error deleting url")
		return openapi.ImplResponse{}, err
	}
	return openapi.ImplResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	var shortURL string
	if s.UseURLShortenerService {
		// use the urlshortener service
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

		shortURL = u.Url
	} else {
		shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
		if err != nil {
			log.Error().Err(err).Msg("alias does not exist")
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

	res := openapi.ImplResponse{
		Code: http.StatusPermanentRedirect,
		Headers: map[string][]string{
			"Location": {shortURL},
		},
	}

	// add location hedaer for redirect in the response
	return res, nil
}

func (s *URLApiService) GetUrlData(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		log.Error().Err(err).Msg("alias does not exist")
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Message: "alias does not exist",
				Code:    http.StatusNotFound,
			},
		}, nil
	}

	res := openapi.ImplResponse{
		Code: http.StatusOK,
		Body: openapi.Url{
			Alias: shortedURL.Alias,
			Url:   shortedURL.Url,
		},
	}

	// add location hedaer for redirect in the response
	return res, nil
}
