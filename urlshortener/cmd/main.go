package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/url"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Port       int    `long:"port" env:"PORT" default:"8080" description:"HTTP port"`
	HTTPAddr   string `long:"http-addr" env:"HTTP_ADDR" defatult:"" description:"HTTP address"`
	DBURL      string `long:"db-url" env:"DB_URL" default:"postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" description:"Database URL"`
	logLevel   string `long:"log-level" env:"LOG_LEVEL" default:"info" description:"Log level"`
	MetaServer struct {
		Addr string `long:"addr" env:"ADDR" default:"localhost" description:"Meta service address"`
		Port int    `long:"port" env:"PORT" default:"8081" description:"Meta service port"`
	} `group:"Meta service" namespace:"meta" env-namespace:"META"`
}

func main() {
	// Parse flags
	if _, err := flags.Parse(&opts); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Fatal().Err(err).Msg("failed to parse flags")
		}
		os.Exit(0)
	}

	// log config print stack trace
	log.Logger = log.With().Caller().Logger()
	// Set log level zerolog
	if _, err := zerolog.ParseLevel(opts.logLevel); err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	// Connect to database
	dbConn, err := sql.Open("nrpostgres", opts.DBURL)
	if err != nil {
		log.Fatal().Err(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	if dbConn == nil {
		log.Fatal().Msg("db is nil")
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbConn.PingContext(pingCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database, shutting down")
	}

	// Create queries instance for db calls
	queries := db.New()

	// Create URL API service and controller
	URLAPIService := &url.URLApiService{
		DB:     queries,
		DBConn: dbConn,
	}
	URLAPIController := server.NewURLAPIController(URLAPIService)

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start meta service
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", opts.MetaServer.Addr, opts.MetaServer.Port), nil); err != nil {
			log.Fatal().Err(err).Msg("failed to start meta service")
		}
	}()

	// Start HTTP server
	urlRouter := server.NewRouter(URLAPIController)
	addr := fmt.Sprintf("%s:%d", opts.HTTPAddr, opts.Port)
	log.Info().Str("addr", addr).Msg("starting HTTP server")
	if err := http.ListenAndServe(addr, urlRouter); err != nil {
		log.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}
