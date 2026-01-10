package url

import (
	"context"
	"net/http"
	"strings"

	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/contexthelper"
	"github.com/igomez10/microservices/urlshortener/pkg/converter"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
	"github.com/igomez10/microservices/urlshortener/pkg/tracerhelper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// validate URLApiService implements the URLApiServicer interface
var _ server.URLAPIServicer = (*URLApiService)(nil)

const appName = "urlshortener"

// implements the URLApiServicer interface
// s *URLApiService server.URLApiServicer
type URLApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *URLApiService) CreateUrl(ctx context.Context, newURL server.Url, requestID string) (serverRes server.ImplResponse, err error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateUrl")
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}()

	newURLParams := db.CreateURLParams{
		Alias: newURL.Alias,
		Url:   newURL.Url,
	}

	ctx, dbspan := tracerhelper.GetTracer().Start(ctx, "db.CreateUrl")
	res, err := s.DB.CreateURL(ctx, s.DBConn, newURLParams)
	dbspan.End()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return server.ImplResponse{
				Code: http.StatusConflict,
				Body: server.Error{
					Message: "url with alias already exists",
					Code:    http.StatusConflict,
				},
			}, err
		}

		// other error
		logger := contexthelper.GetLoggerInContext(ctx)
		logger.Error("error creating url", "url", newURL.Url, "alias", newURL.Alias, "err", err)
		return server.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: server.Error{
				Message: "error creating url",
				Code:    http.StatusInternalServerError,
			},
		}, err
	}

	serverURL := converter.FromDBUrlToAPIUrl(res)
	return server.ImplResponse{
		Code: http.StatusOK,
		Body: serverURL,
	}, nil

}

func (s *URLApiService) DeleteUrl(ctx context.Context, alias string, requestID string) (server.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "DeleteUrl")
	defer span.End()

	logger := contexthelper.GetLoggerInContext(ctx)
	ctx, dbspan := tracerhelper.GetTracer().Start(ctx, "db.DeleteUrl")
	if err := s.DB.DeleteURL(ctx, s.DBConn, alias); err != nil {
		logger.Error("error deleting url", "alias", alias, "err", err)
		return server.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: server.Error{
				Message: "error deleting url",
				Code:    http.StatusInternalServerError,
			},
		}, err
	}
	dbspan.End()

	return server.ImplResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *URLApiService) GetUrl(ctx context.Context, alias string, requestID string) (server.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrl")
	defer span.End()

	ctx, dbspan := tracerhelper.GetTracer().Start(ctx, "db.GetURLFromAlias")
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	dbspan.End()
	if err != nil {
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

func (s *URLApiService) GetUrlData(ctx context.Context, alias string, requestID string) (server.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetUrlData")
	defer span.End()

	ctx, dbspan := tracerhelper.GetTracer().Start(ctx, "db.GetURLFromAlias")
	shortedURL, err := s.DB.GetURLFromAlias(ctx, s.DBConn, alias)
	dbspan.End()
	if err != nil {
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

	span.AddEvent("GetUrlData")
	span.SetAttributes(
		attribute.KeyValue{
			Key:   "alias",
			Value: attribute.StringValue(alias),
		},
		attribute.KeyValue{
			Key:   "url",
			Value: attribute.StringValue(shortedURL.Url),
		},
	)
	span.SetStatus(codes.Ok, "GetUrlData")

	return res, nil
}
