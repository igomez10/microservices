package main

import (
	"database/sql"
	"net/http"
	"socialapp/pkg/controller/comment"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	go startPrometheusServer()
	log.Printf("Server started")
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

	CommentApiService := &comment.CommentService{DB: queries, DBConn: dbConn}
	CommentApiController := openapi.NewCommentApiController(CommentApiService)

	UserApiService := &user.UserApiService{DB: queries, DBConn: dbConn}
	UserApiController := openapi.NewUserApiController(UserApiService)

	router := openapi.NewRouter(CommentApiController, UserApiController)
	log.Debug().Msg("Server is listening on port :8080")
	log.Fatal().Err(http.ListenAndServe(":8080", router))
}

func startPrometheusServer() error {
	addr := ":9095"
	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
