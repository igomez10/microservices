package url

import (
	"context"
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	urlClient "github.com/igomez10/microservices/urlshortener/generated/clients/go/client"
)

// implements the URLApiServicer interface
// s *URLApiService openapi.URLApiServicer
type URLApiService struct {
	// urlClient is the Client for the urlshortener service
	Client *urlClient.APIClient
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
	logger := contexthelper.GetLoggerInContext(ctx)
	var openapiURL openapi.Url
	newURLRequest := urlClient.NewURL(newURL.Url, newURL.Alias)
	u, createRes, err := s.Client.URLAPI.CreateUrl(ctx).
		URL(*newURLRequest).
		XRequestID(contexthelper.GetRequestIDInContext(ctx)).
		Execute()
	if err != nil {
		logger.Error("unexpected error creating url", "error", err, "url", newURL.Url, "alias", newURL.Alias)
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Message: "error creating url",
				Code:    http.StatusInternalServerError,
			},
		}, nil
	}

	switch createRes.StatusCode {
	case http.StatusConflict:
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

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: openapiURL,
	}, nil

}

// DeleteUrl deletes a url
func (s *URLApiService) DeleteUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "DeleteUrl")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)

	// use the urlshortener service
	res, err := s.Client.URLAPI.
		DeleteUrl(ctx, alias).
		XRequestID(contexthelper.GetRequestIDInContext(ctx)).
		Execute()
	if err != nil || res.StatusCode != http.StatusOK {
		logger.Error("error deleting url", "error", err)
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Message: "error deleting url",
				Code:    http.StatusInternalServerError,
			},
		}, err
	}
	return openapi.ImplResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrl")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)

	var shortURL string
	// use the urlshortener service
	u, res, err := s.Client.URLAPI.
		GetUrlData(ctx, alias).
		XRequestID(contexthelper.GetRequestIDInContext(ctx)).
		Execute()
	if err != nil || res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			logger.Debug("alias does not exist", "error", err)
			return openapi.ImplResponse{
				Code: http.StatusNotFound,
				Body: openapi.Error{
					Message: "alias does not exist",
					Code:    http.StatusNotFound,
				},
			}, nil
		default:
			logger.Error("error getting url from urlshortener service", "error", err, "status_code", res.StatusCode)
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

	// add location header for redirect in the response
	apiRes := openapi.ImplResponse{
		Code: http.StatusPermanentRedirect,
		Headers: map[string][]string{
			"Location": {shortURL},
		},
	}

	return apiRes, nil
}

func (s *URLApiService) GetUrlData(ctx context.Context, alias string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrlData")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx)

	var responseUrl openapi.Url
	u, apiRes, err := s.Client.URLAPI.
		GetUrlData(ctx, alias).
		XRequestID(contexthelper.GetRequestIDInContext(ctx)).
		Execute()
	if err != nil || apiRes.StatusCode != http.StatusOK {
		logger.Error("error getting url from urlshortener service", "error", err)
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

	res := openapi.ImplResponse{
		Code: http.StatusOK,
		Body: responseUrl,
	}

	return res, nil
}
