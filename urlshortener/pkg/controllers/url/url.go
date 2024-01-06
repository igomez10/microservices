package url

import (
	"context"
	"net/http"

	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/contexthelper"
	"github.com/igomez10/microservices/urlshortener/pkg/converter"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
)

// validate URLApiService implements the URLApiServicer interface
var _ server.URLAPIServicer = (*URLApiService)(nil)

// implements the URLApiServicer interface
// s *URLApiService server.URLApiServicer
type URLApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *URLApiService) CreateUrl(ctx context.Context, newURL server.Url) (server.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	_, err := s.DB.GetURLFromAlias(ctx, s.DBConn, newURL.Alias)
	if err == nil {
		log.Error().Err(err).Msg("url with alias already exists")
		return server.ImplResponse{
			Code: http.StatusConflict,
			Body: server.Error{
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
		return server.ImplResponse{}, err
	}

	serverURL := converter.FromDBUrlToAPIUrl(res)
	return server.ImplResponse{
		Code: http.StatusOK,
		Body: serverURL,
	}, nil

}

func (s *URLApiService) DeleteUrl(ctx context.Context, alias string) (server.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	if err := s.DB.DeleteURL(ctx, s.DBConn, alias); err != nil {
		log.Error().Err(err).Msg("error deleting url")
		return server.ImplResponse{}, err
	}
	return server.ImplResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string) (server.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		log.Error().Err(err).Msg("alias does not exist")
		return server.ImplResponse{
			Code: http.StatusNotFound,
			Body: server.Error{
				Message: "alias does not exist",
				Code:    http.StatusNotFound,
			},
		}, nil
	}

	res := server.ImplResponse{
		Code: http.StatusPermanentRedirect,
		Headers: map[string][]string{
			"Location": {shortedURL.Url},
		},
	}

	// add location hedaer for redirect in the response
	return res, nil
}

func (s *URLApiService) GetUrlData(ctx context.Context, alias string) (server.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		log.Error().Err(err).Msg("alias does not exist")
		return server.ImplResponse{
			Code: http.StatusNotFound,
			Body: server.Error{
				Message: "alias does not exist",
				Code:    http.StatusNotFound,
			},
		}, nil
	}

	res := server.ImplResponse{
		Code: http.StatusOK,
		Body: server.Url{
			Alias: shortedURL.Alias,
			Url:   shortedURL.Url,
		},
	}

	// add location hedaer for redirect in the response
	return res, nil
}
