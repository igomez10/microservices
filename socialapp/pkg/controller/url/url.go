package url

import (
	"context"
	"net/http"
	"socialapp/internal/contexthelper"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
)

// implements the URLApiServicer interface
// s *URLApiService openapi.URLApiServicer
type URLApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *URLApiService) CreateUrl(ctx context.Context, newURL openapi.Url) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	_, err := s.DB.GetURLFromAlias(ctx, s.DBConn, newURL.Alias)
	if err == nil {
		log.Error().Err(err).Msg("url with alias already exists")
		return openapi.ImplResponse{
			Code: http.StatusConflict,
			Body: openapi.Error{
				Message: "url with alias already exists",
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
	return openapi.ImplResponse{
		Code: http.StatusCreated,
		Body: res,
	}, nil

}

func (s *URLApiService) DeleteUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		log.Error().Err(err).Msg("url with alias already exists")
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Message: "url with alias already exists",
			},
		}, nil
	}

	res := openapi.ImplResponse{
		Code: http.StatusOK,
		Headers: map[string][]string{
			"Location": {shortedURL.Url},
		},
	}

	// add location hedaer for redirect in the response
	return res, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)

	// validate we dont have a url with the same alias
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		log.Error().Err(err).Msg("alias does not exist")
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Message: "alias does not exist",
			},
		}, nil
	}

	res := openapi.ImplResponse{
		Code: http.StatusPermanentRedirect,
		Headers: map[string][]string{
			"Location": {shortedURL.Url},
		},
		Body: openapi.Url{
			Alias: shortedURL.Alias,
			Url:   shortedURL.Url,
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
