package url

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/contexthelper"
	"github.com/igomez10/microservices/urlshortener/pkg/converter"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
	"github.com/igomez10/microservices/urlshortener/pkg/tracerhelper"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel/trace"
)

// validate URLApiService implements the URLApiServicer interface
var _ server.URLAPIServicer = (*URLApiService)(nil)

const appName = "urlshortener"

// implements the URLApiServicer interface
// s *URLApiService server.URLApiServicer
type URLApiService struct {
	DB     db.Querier
	DBConn db.DBTX

	metrics newrelic.Application
}

type MetricEvent struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
	IsErr bool   `json:"is_err"`
}

func (m *MetricEvent) toMap() map[string]interface{} {
	return map[string]interface{}{
		"alias":  m.Alias,
		"url":    m.Url,
		"is_err": m.IsErr,
	}
}

func (s *URLApiService) CreateUrl(ctx context.Context, newURL server.Url) (server.ImplResponse, error) {
	// ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateUrl")
	span := trace.SpanFromContext(ctx)

	defer span.End()

	event := MetricEvent{
		Alias: newURL.Alias,
		Url:   newURL.Url,
	}
	defer func() {
		s.metrics.RecordCustomEvent("CreateUrl", event.toMap())
	}()

	newURLParams := db.CreateURLParams{
		Alias: newURL.Alias,
		Url:   newURL.Url,
	}
	res, err := s.DB.CreateURL(ctx, s.DBConn, newURLParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return server.ImplResponse{
				Code: http.StatusConflict,
				Body: server.Error{
					Message: fmt.Sprintf("url with alias %q already exists", newURL.Alias),
					Code:    http.StatusConflict,
				},
			}, nil
		}

		// other error
		event.IsErr = true
		return server.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: server.Error{
				Message: fmt.Sprintf("error creating url %q with alias %q", newURL.Url, newURL.Alias),
			},
		}, err
	}

	serverURL := converter.FromDBUrlToAPIUrl(res)
	return server.ImplResponse{
		Code: http.StatusOK,
		Body: serverURL,
	}, nil

}

func (s *URLApiService) DeleteUrl(ctx context.Context, alias string) (server.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "DeleteUrl")
	defer span.End()

	log := contexthelper.GetLoggerInContext(ctx)
	event := MetricEvent{
		Alias: alias,
		IsErr: false,
	}

	defer func() {
		s.metrics.RecordCustomEvent("DeleteUrl", event.toMap())
	}()

	if err := s.DB.DeleteURL(ctx, s.DBConn, alias); err != nil {
		log.Error().Err(err).Msgf("error deleting url %q", alias)
		event.IsErr = true
		return server.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: server.Error{
				Message: "error deleting url",
				Code:    http.StatusInternalServerError,
			},
		}, err
	}
	return server.ImplResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string) (server.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrl")
	defer span.End()

	event := MetricEvent{
		Alias: alias,
		IsErr: false,
	}
	defer func() {
		s.metrics.RecordCustomEvent("GetUrl", event.toMap())
	}()

	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		event.IsErr = true
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
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrlData")
	defer span.End()

	event := MetricEvent{
		Alias: alias,
	}

	defer func() {
		s.metrics.RecordCustomEvent("GetUrlData", event.toMap())
	}()

	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	if err != nil {
		event.IsErr = true
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

	return res, nil
}
