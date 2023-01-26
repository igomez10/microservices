package cache

import (
	"context"
	"net/http"
	"socialapp/internal/responseWriter"
	"time"

	"github.com/go-redis/redis/v8"
	newrelic "github.com/newrelic/go-agent"
	nrredis "github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	Client *redis.Client
}

type CacheConfig struct {
	RedisOpts *redis.Options
}

func NewCache(config CacheConfig) *Cache {
	client := redis.NewClient(config.RedisOpts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to connect to redis")
	}

	// New Relic setup for redis
	client.AddHook(nrredis.NewHook(config.RedisOpts))

	c := &Cache{
		Client: client,
	}
	return c
}

var metricRedisCahe = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "socialapp_api_cache",
	Help: "The total number of cache hits",
}, []string{"cache", "status"})

func (c *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if cache should be used
		shouldSearchCache := true
		if r.Header.Get("Cache-Control") == "no-store" {
			metricRedisCahe.WithLabelValues("redis", "skipped").Inc()
			shouldSearchCache = false
		}
		customW := responseWriter.NewCustomResponseWriter(w)
		if r.Method == http.MethodGet && shouldSearchCache {
			txn := newrelic.FromContext(r.Context())
			r = r.WithContext(newrelic.NewContext(r.Context(), txn))

			// attempt to return here, if not found in cache then continue to handler
			key := r.Method + "+" + r.URL.Path + r.URL.RawQuery
			// search in cache
			val, err := c.Client.Get(r.Context(), key).Bytes()
			if err == nil {
				// found in cache
				metricRedisCahe.WithLabelValues("redis", "hit").Inc()
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
					metricRedisCahe.WithLabelValues("redis", "miss").Inc()
					key := r.Method + "+" + r.URL.Path

					if err := c.Client.Set(r.Context(), key, customW.Body, time.Minute*10).Err(); err != nil {
						log.Error().Stack().Err(err).Msg("Failed to set key in redis")
					}
				}
			}()
		}

		next.ServeHTTP(customW, r)

	})
}
