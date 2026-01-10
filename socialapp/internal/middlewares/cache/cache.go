package cache

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/igomez10/microservices/socialapp/internal/responseWriter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RedisClient interface for testing
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type Cache struct {
	Client RedisClient
}

type CacheConfig struct {
	RedisOpts *redis.Options
}

func NewCache(config CacheConfig) *Cache {
	client := redis.NewClient(config.RedisOpts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("Failed to connect to redis", "error", err)
		os.Exit(1)
	}

	c := &Cache{
		Client: client,
	}
	return c
}

var metricRedisCache = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "socialapp_api_cache",
	Help: "The total number of cache hits",
}, []string{"cache", "status"})

func (c *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.cache")
		defer span.End()

		r = r.WithContext(ctx)
		// check if cache should be used
		shouldSearchCache := true
		if r.Header.Get("Cache-Control") == "no-store" {
			metricRedisCache.WithLabelValues("redis", "skipped").Inc()
			shouldSearchCache = false
		}
		customW := responseWriter.NewCustomResponseWriter(w)
		if r.Method == http.MethodGet && shouldSearchCache {

			// attempt to return here, if not found in cache then continue to handler
			key := r.Method + "+" + r.URL.Path
			if r.URL.RawQuery != "" {
				key += "?" + r.URL.RawQuery
			}
			// search in cache
			val, err := c.Client.Get(r.Context(), key).Bytes()
			if err == nil {
				// found in cache
				metricRedisCache.WithLabelValues("redis", "hit").Inc()
				customW.Header().Set("X-Cache", "HIT")
				customW.Header().Set("Content-Type", "application/json")
				customW.Header().Set("Cache-Control", "public, max-age=3600")
				customW.Write(val)

				return
			}

			// not found in cache
			defer func() {
				// if response is 200, cache response
				if customW.StatusCode == http.StatusOK {
					// count this response as a miss because it was a 200 and we didnt had the value in the cache
					metricRedisCache.WithLabelValues("redis", "miss").Inc()
					key := r.Method + "+" + r.URL.Path
					if r.URL.RawQuery != "" {
						key += "?" + r.URL.RawQuery
					}

					if err := c.Client.Set(r.Context(), key, customW.Body, time.Minute*10).Err(); err != nil {
						slog.Error("Failed to set key in redis", "error", err)
					}
				}
			}()
		}

		next.ServeHTTP(customW, r)

	})
}
