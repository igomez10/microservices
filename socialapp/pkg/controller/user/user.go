package user

import (
	"context"
	"net/http"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
	"strings"

	"github.com/rs/zerolog/log"
)

// implements the UserService interface
// s *UserApiService openapi.UserApiServicer
type UserApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *UserApiService) CreateUser(ctx context.Context, user openapi.User) (openapi.ImplResponse, error) {
	// validate we dont have a user with the same username that is not deleted
	noCaseUsername := strings.ToLower(user.Username)
	if _, err := s.DB.GetUserByUsername(ctx, s.DBConn, noCaseUsername); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	// validate we dont have a user with the same email that is not deleted
	noCaseEmail := strings.ToLower(user.Email)
	if _, err := s.DB.GetUserByEmail(ctx, s.DBConn, noCaseEmail); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	params := db.CreateUserParams{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	createdUser, err := s.DB.CreateUser(ctx, s.DBConn, params)
	if err != nil {
		log.Err(err).Msg("Error creating user")
		return openapi.Response(http.StatusInternalServerError, nil), nil
	}

	return openapi.Response(http.StatusOK, createdUser), nil
}

// DeleteUser - Deletes a particular user
func (s *UserApiService) DeleteUser(ctx context.Context, username string) (openapi.ImplResponse, error) {
	if err := s.DB.DeleteUserByUsername(ctx, s.DBConn, username); err != nil {
		log.Err(err).Msg("Error deleting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, nil), nil
}

// GetUserByUsername - Get a particular user by username
func (s *UserApiService) GetUserByUsername(ctx context.Context, username string) (openapi.ImplResponse, error) {
	u, err := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, u), nil
}

// GetUserComments - Gets all comments for a user
func (s *UserApiService) GetUserComments(ctx context.Context, username string) (openapi.ImplResponse, error) {
	commnet, err := s.DB.GetUserComments(ctx, s.DBConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting user comments")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

// ListUsers - Returns all the users
func (s *UserApiService) ListUsers(ctx context.Context) (openapi.ImplResponse, error) {
	commnet, err := s.DB.ListUsers(ctx, s.DBConn)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

func (s *UserApiService) UpdateUser(ctx context.Context, username string, user openapi.User) (openapi.ImplResponse, error) {
	// validate we dont have a user with the same username that is not deleted
	noCaseUsername := strings.ToLower(user.Username)
	if _, err := s.DB.GetUserByUsername(ctx, s.DBConn, noCaseUsername); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	// validate we dont have a user with the same email that is not deleted
	noCaseEmail := strings.ToLower(user.Email)
	if _, err := s.DB.GetUserByEmail(ctx, s.DBConn, noCaseEmail); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	params := db.UpdateUserByUsernameParams{
		OldUsername: username,
		NewUsername: user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
	}

	uUser, err := s.DB.UpdateUserByUsername(ctx, s.DBConn, params)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, uUser), nil
}
