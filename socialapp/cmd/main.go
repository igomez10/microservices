package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"socialapp/pkg/controller/comment"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var (
	appPort  *int = flag.Int("port", 8080, "main port for application")
	metaPort *int = flag.Int("meta_port", 9095, "meta port for metric/service_discovery/etc")
)

func main() {
	flag.Parse()
	log.Info().Msgf("Starting PORT: %d, METAPORT: %d", *appPort, *metaPort)
	go startPrometheusServer(*metaPort)
	dbConn, err := sql.Open("postgres", "postgres://postgres:password@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err)
	}
	defer dbConn.Close()

	if dbConn == nil {
		log.Fatal().Msg("db is nil")
	}
	defer dbConn.Close()

	queries := db.New()

	// Comment service
	CommentApiService := &comment.CommentService{DB: queries, DBConn: dbConn}
	CommentApiController := openapi.NewCommentApiController(CommentApiService)

	// User service
	UserApiService := &user.UserApiService{DB: queries, DBConn: dbConn}
	UserApiController := openapi.NewUserApiController(UserApiService)

	router := openapi.NewRouter(CommentApiController, UserApiController)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *appPort), router); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}

}

func startPrometheusServer(port int) error {

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
