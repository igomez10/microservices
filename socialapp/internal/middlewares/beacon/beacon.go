package beacon

import (
	"net/http"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/responseWriter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

// Beacon is used to configure a centralized logging platform for the application.
type Beacon struct {
}

// Middleware returns a handler that logs all the request/response information for
// the centralized logging platform
// Middleware function for handling HTTP requests and logging their responses.
func (b *Beacon) Middleware(next http.Handler) http.Handler {
	beaconMeter := otel.GetMeterProvider().Meter("beacon")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.beacon")
		defer span.End()

		r = r.WithContext(ctx)

		customW := responseWriter.NewCustomResponseWriter(w)
		logger := contexthelper.GetLoggerInContext(r.Context())
		r = r.WithContext(contexthelper.SetLoggerInContext(r.Context(), logger))
		startTime := time.Now()

		// Handle the request
		next.ServeHTTP(customW, r)

		attributes := attribute.NewSet(
			attribute.String("pattern", contexthelper.GetRequestPatternInContext(r.Context())),
			attribute.String("method", r.Method),
			attribute.String("user_agent", r.UserAgent()),
			attribute.String("request_host", r.Host),
			attribute.Int("status_code", customW.StatusCode),
		)
		hist, err := beaconMeter.Float64Histogram("http_server_response_duration",
			metric.WithExplicitBucketBoundaries([]float64{
				0.001,
				0.002,
				0.003,
				0.004,
				0.005,
				0.006,
				0.007,
				0.008,
				0.009,
				0.010,
				0.015,
				0.020,
				0.025,
				0.030,
				0.040,
				0.050,
				0.075,
				0.100,
				0.150,
				0.200,
				0.250,
				0.300,
				0.350,
				0.400,
				0.500,
				0.750,
				1.000,
				2.000,
				5.000,
				10.000,
			}...))
		if err != nil {
			log.Error().Err(err).Msg("failed to create histogram in beacon")
		} else {
			hist.Record(ctx, time.Since(startTime).Seconds(), metric.WithAttributeSet(attributes))
		}

		span.SetAttributes([]attribute.KeyValue{
			{
				Key:   "http.method",
				Value: attribute.StringValue(r.Method),
			},
			{
				Key:   "http.url",
				Value: attribute.StringValue(r.URL.String()),
			},
			{
				Key:   "http.status_code",
				Value: attribute.IntValue(customW.StatusCode),
			},
			{
				Key:   "http.user_agent",
				Value: attribute.StringValue(r.UserAgent()),
			},
			{
				Key:   "http.request_id",
				Value: attribute.StringValue(r.Header.Get("X-Request-ID")),
			},
		}...)

		if customW.StatusCode >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
		} else {
			span.SetStatus(codes.Ok, "HTTP success")
		}

		// Log the response
		statusCode := customW.StatusCode
		logEvent := createLogEvent(r, statusCode, startTime, &logger)
		logEvent.Msgf("finished request")
	})
}

// Create a log event based on the request and response.
func createLogEvent(r *http.Request, statusCode int, startTime time.Time, logger *zerolog.Logger) *zerolog.Event {
	requestPattern := contexthelper.GetRequestPatternInContext(r.Context())
	latency := time.Since(startTime).Milliseconds()

	logevent := logger.WithLevel(zerolog.InfoLevel).
		Str("path", r.URL.Path).
		Str("query", r.URL.RawQuery).
		Str("pattern", requestPattern).
		Str("method", r.Method).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent()).
		Str("referer", r.Referer()).
		Str("request_host", r.Host).
		Int("status_code", statusCode).
		Int64("latency_ms", latency)

	for k, v := range r.Header {
		if k == "Authorization" {
			continue
		}

		logevent.Strs(k, v)
	}

	switch statusCode {
	case http.StatusUnauthorized:
		username, _ := contexthelper.GetUsernameInContext(r.Context())
		return logevent.
			Str("authorization", r.Header.Get("Authorization")).
			Str("username", username).
			Str("error", "Unauthorized")
	default:
		return logevent
	}
}
