package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func main() {
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

	CommentApiService := &CommentService{db: *queries, dbConn: dbConn}
	CommentApiController := openapi.NewCommentApiController(CommentApiService)

	UserApiService := &UserApiService{db: *queries, dbConn: dbConn}
	UserApiController := openapi.NewUserApiController(UserApiService)

	router := openapi.NewRouter(CommentApiController, UserApiController)
	log.Debug().Msg("Server is listening on port :8080")
	log.Fatal().Err(http.ListenAndServe(":8080", router))
}

type CommentService struct {
	db     db.Queries
	dbConn db.DBTX
}

func (s *CommentService) CreateComment(ctx context.Context, username string, comment openapi.Comment) (openapi.ImplResponse, error) {
	params := db.CreateCommentForUserParams{
		Username: username,
		Content:  comment.Content,
	}

	newComment, err := s.db.CreateCommentForUser(ctx, s.dbConn, params)
	if err != nil {
		log.Err(err).Msg("Error creating comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, newComment), nil
}

func (s *CommentService) GetComment(ctx context.Context, id int32) (openapi.ImplResponse, error) {
	commnet, err := s.db.GetComment(ctx, s.dbConn, id)
	if err != nil {
		log.Err(err).Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

func (s *CommentService) GetUserComments(ctx context.Context, username string) (openapi.ImplResponse, error) {
	commnet, err := s.db.GetUserComments(ctx, s.dbConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting comment")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

type UserApiService struct {
	db     db.Queries
	dbConn db.DBTX
}

func (s *UserApiService) CreateUser(ctx context.Context, user openapi.User) (openapi.ImplResponse, error) {

	// validate we dont have a user with the same username that is not deleted
	if _, err := s.db.GetUserByUsername(ctx, s.dbConn, user.Username); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	// validate we dont have a user with the same email that is not deleted
	if _, err := s.db.GetUserByEmail(ctx, s.dbConn, user.Email); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	params := db.CreateUserParams{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	createdUser, err := s.db.CreateUser(ctx, s.dbConn, params)
	if err != nil {
		log.Err(err).Msg("Error creating user")
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	return openapi.Response(http.StatusOK, createdUser), nil
}

// s *UserApiService openapi.UserApiServicer

// DeleteUser - Deletes a particular user
func (s *UserApiService) DeleteUser(ctx context.Context, username string) (openapi.ImplResponse, error) {
	if err := s.db.DeleteUserByUsername(ctx, s.dbConn, username); err != nil {
		log.Err(err).Msg("Error deleting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil
}

// GetUserByUsername - Get a particular user by username
func (s *UserApiService) GetUserByUsername(ctx context.Context, username string) (openapi.ImplResponse, error) {
	u, err := s.db.GetUserByUsername(ctx, s.dbConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, u), nil
}

// GetUserComments - Gets all comments for a user
func (s *UserApiService) GetUserComments(ctx context.Context, username string) (openapi.ImplResponse, error) {
	commnet, err := s.db.GetUserComments(ctx, s.dbConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting user comments")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

// ListUsers - Returns all the users
func (s *UserApiService) ListUsers(ctx context.Context) (openapi.ImplResponse, error) {
	commnet, err := s.db.ListUsers(ctx, s.dbConn)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

// UpdateUser - Update a user
func (s *UserApiService) UpdateUser(ctx context.Context, username string) (openapi.ImplResponse, error) {
	// TODO - update UpdateUser with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return openapi.Response Response(200, User{}) or use other options such as http.Ok ...
	//return openapi.Response(200, User{}), nil

	//TODO: Uncomment the next line to return openapi.Response Response(0, Error{}) or use other options such as http.Ok ...
	//return openapi.Response(0, Error{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("UpdateUser method not implemented")
}
