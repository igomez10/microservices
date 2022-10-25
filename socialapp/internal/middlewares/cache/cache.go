package cache

import (
	"net/http"
	"socialapp/internal/responseWriter"
	"time"

	"github.com/go-redis/redis"
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
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to connect to redis")
	}

	c := &Cache{
		Client: client,
	}
	return c
}

func (c *Cache) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if cache should be used
		shouldSearchCache := true
		if r.Header.Get("Cache-Control") == "no-store" {
			shouldSearchCache = false
		}
		customW := responseWriter.NewCustomResponseWriter(w)
		if r.Method == http.MethodGet && shouldSearchCache {
			// attempt to return here, if not found in cache, continue to handler
			key := r.Method + "+" + r.URL.Path
			val, err := c.Client.Get(key).Result()
			if err != nil {
				log.Error().Stack().Err(err).Msg("Failed to get key from redis")
			} else {
				log.Info().Msgf("Found key %q in redis", key)
				customW.Write([]byte(val))
				customW.Header().Add("Cache-Control", "public, max-age=3600")
				return
			}

			defer func() {
				// if response is 200, cache response
				if customW.StatusCode == http.StatusOK {
					log.Info().Msgf("Caching response for key %q", r.Method+"+"+r.URL.Path)
					key := r.Method + "+" + r.URL.Path

					err := c.Client.Set(key, customW.Body, time.Minute*10).Err()
					if err != nil {
						log.Error().Stack().Err(err).Msg("Failed to set key in redis")
					}
				}
			}()
		}

		next.ServeHTTP(customW, r)

	})
}
