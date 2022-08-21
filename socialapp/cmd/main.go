package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"socialapp/pkg/db"
	socialappapi "socialapp/socialappapi/go"

	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Server started")
	dbConn, err := sql.Open("postgres", "postgres://postgres:password@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	if dbConn == nil {
		log.Fatal("db is nil")
	}
	defer dbConn.Close()

	queries := db.New()

	CommentApiService := &CommentService{db: *queries, dbConn: dbConn}
	CommentApiController := socialappapi.NewCommentApiController(CommentApiService)

	UserApiService := &UserApiService{db: *queries, dbConn: dbConn}
	UserApiController := socialappapi.NewUserApiController(UserApiService)

	router := socialappapi.NewRouter(CommentApiController, UserApiController)
	log.Println("Server is listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type CommentService struct {
	db     db.Queries
	dbConn db.DBTX
}

func (s *CommentService) CreateComment(ctx context.Context, username string, comment socialappapi.Comment) (socialappapi.ImplResponse, error) {
	params := db.CreateCommentForUserParams{
		Username: username,
		Content:  comment.Content,
	}

	newComment, err := s.db.CreateCommentForUser(ctx, s.dbConn, params)
	if err != nil {
		return socialappapi.Response(http.StatusNotFound, nil), nil
	}

	return socialappapi.Response(http.StatusOK, newComment), nil
}

func (s *CommentService) GetComment(ctx context.Context, id int32) (socialappapi.ImplResponse, error) {
	commnet, err := s.db.GetComment(ctx, s.dbConn, id)
	if err != nil {
		return socialappapi.Response(http.StatusNotFound, nil), nil
	}

	return socialappapi.Response(http.StatusOK, commnet), nil
}

func (s *CommentService) GetUserComments(ctx context.Context, username string) (socialappapi.ImplResponse, error) {
	commnet, err := s.db.GetUserComments(ctx, s.dbConn, username)
	if err != nil {
		return socialappapi.Response(http.StatusNotFound, nil), nil
	}

	return socialappapi.Response(http.StatusOK, commnet), nil
}

type UserApiService struct {
	db     db.Queries
	dbConn db.DBTX
}

// CreateUser - Create a new user
func (s *UserApiService) CreateUser(ctx context.Context) (socialappapi.ImplResponse, error) {
	// TODO - update CreateUser with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return socialappapi.Response Response(200, User{}) or use other options such as http.Ok ...
	//return socialappapi.Response(200, User{}), nil

	//TODO: Uncomment the next line to return socialappapi.Response Response(0, Error{}) or use other options such as http.Ok ...
	//return socialappapi.Response(0, Error{}), nil

	return socialappapi.Response(http.StatusNotImplemented, nil), errors.New("CreateUser method not implemented")
}

// DeleteUser - Deletes a particular user
func (s *UserApiService) DeleteUser(ctx context.Context, username string) (socialappapi.ImplResponse, error) {
	// TODO - update DeleteUser with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return socialappapi.Response Response(200, User{}) or use other options such as http.Ok ...
	//return socialappapi.Response(200, User{}), nil

	//TODO: Uncomment the next line to return socialappapi.Response Response(0, Error{}) or use other options such as http.Ok ...
	//return socialappapi.Response(0, Error{}), nil

	return socialappapi.Response(http.StatusNotImplemented, nil), errors.New("DeleteUser method not implemented")
}

// GetUserByUsername - Get a particular user by username
func (s *UserApiService) GetUserByUsername(ctx context.Context, username string) (socialappapi.ImplResponse, error) {
	u, err := s.db.GetUserByUsername(ctx, s.dbConn, username)
	if err != nil {
		return socialappapi.Response(http.StatusNotFound, nil), nil
	}

	return socialappapi.Response(http.StatusOK, u), nil
}

// GetUserComments - Gets all comments for a user
func (s *UserApiService) GetUserComments(ctx context.Context, username string) (socialappapi.ImplResponse, error) {
	// TODO - update GetUserComments with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return socialappapi.Response Response(200, Comment{}) or use other options such as http.Ok ...
	//return socialappapi.Response(200, Comment{}), nil

	//TODO: Uncomment the next line to return socialappapi.Response Response(0, Error{}) or use other options such as http.Ok ...
	//return socialappapi.Response(0, Error{}), nil
	commnet, err := s.db.GetUserComments(ctx, s.dbConn, username)
	if err != nil {
		return socialappapi.Response(http.StatusNotFound, nil), nil
	}

	return socialappapi.Response(http.StatusOK, commnet), nil
}

// ListUsers - Returns all the users
func (s *UserApiService) ListUsers(ctx context.Context) (socialappapi.ImplResponse, error) {
	// TODO - update ListUsers with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return socialappapi.Response Response(200, []User{}) or use other options such as http.Ok ...
	//return socialappapi.Response(200, []User{}), nil

	//TODO: Uncomment the next line to return socialappapi.Response Response(0, Error{}) or use other options such as http.Ok ...
	//return socialappapi.Response(0, Error{}), nil

	return socialappapi.Response(http.StatusNotImplemented, nil), errors.New("ListUsers method not implemented")
}

// UpdateUser - Update a user
func (s *UserApiService) UpdateUser(ctx context.Context, username string) (socialappapi.ImplResponse, error) {
	// TODO - update UpdateUser with the required logic for this service method.
	// Add api_user_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return socialappapi.Response Response(200, User{}) or use other options such as http.Ok ...
	//return socialappapi.Response(200, User{}), nil

	//TODO: Uncomment the next line to return socialappapi.Response Response(0, Error{}) or use other options such as http.Ok ...
	//return socialappapi.Response(0, Error{}), nil

	return socialappapi.Response(http.StatusNotImplemented, nil), errors.New("UpdateUser method not implemented")
}
