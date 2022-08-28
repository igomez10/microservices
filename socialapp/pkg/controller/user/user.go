package user

import (
	"context"
	"net/http"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

// implements the UserServicer interface
// s *UserApiService openapi.UserApiServicer
type UserApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

var (
	createUserLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "socialapp_user_operation_create_user_latency",
		Help: "The latency of create user operations",
	})
	getUserByUsernameLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "socialapp_user_operation_get_use_by_username_latency",
		Help: "The latency of get user operations",
	})
	listUsersLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "socialapp_user_operation_list_users_latency",
		Help: "The latency of list users operations",
	})
	updateUsersLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "socialapp_user_operation_update_user_latency",
		Help: "The latency of update user operations",
	})
	deleteUsersLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "socialapp_user_operation_delete_user_latency",
		Help: "The latency of delete user operations",
	})
	getUserCommentsLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "socialapp_user_operation_get_user_comments_latency",
		Help: "The latency of get user comments operations",
	})
)

func (s *UserApiService) CreateUser(ctx context.Context, user openapi.User) (openapi.ImplResponse, error) {
	start := time.Now()
	defer createUserLatency.Set(float64(time.Since(start).Milliseconds()))
	// validate we dont have a user with the same username that is not deleted
	if _, err := s.DB.GetUserByUsername(ctx, s.DBConn, user.Username); err == nil {
		return openapi.Response(http.StatusConflict, nil), nil
	}

	// validate we dont have a user with the same email that is not deleted
	if _, err := s.DB.GetUserByEmail(ctx, s.DBConn, user.Email); err == nil {
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
	start := time.Now()
	defer getUserByUsernameLatency.Set(float64(time.Since(start).Nanoseconds()))
	u, err := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, u), nil
}

// GetUserComments - Gets all comments for a user
func (s *UserApiService) GetUserComments(ctx context.Context, username string) (openapi.ImplResponse, error) {
	start := time.Now()
	defer getUserCommentsLatency.Set(float64(time.Since(start).Nanoseconds()))
	commnet, err := s.DB.GetUserComments(ctx, s.DBConn, username)
	if err != nil {
		log.Err(err).Msg("Error getting user comments")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

// ListUsers - Returns all the users
func (s *UserApiService) ListUsers(ctx context.Context) (openapi.ImplResponse, error) {
	start := time.Now()
	defer listUsersLatency.Set(float64(time.Since(start).Nanoseconds()))
	commnet, err := s.DB.ListUsers(ctx, s.DBConn)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, commnet), nil
}

func (s *UserApiService) UpdateUser(ctx context.Context, existingUsername string, newUserData openapi.User) (openapi.ImplResponse, error) {
	// get the user to update
	start := time.Now()
	defer updateUsersLatency.Set(float64(time.Since(start).Nanoseconds()))
	existingUser, err := s.DB.GetUserByUsername(ctx, s.DBConn, existingUsername)
	if err != nil {
		log.Err(err).Str("username", existingUsername).Msg("Error getting user")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	if newUserData.Username != "" && newUserData.Username != existingUser.Username {
		// validate we dont have a user with the same username that is not deleted
		noCaseUsername := strings.ToLower(newUserData.Username)
		if _, err := s.DB.GetUserByUsername(ctx, s.DBConn, noCaseUsername); err == nil {
			log.Error().Msg("Username already exists")
			return openapi.Response(http.StatusConflict, nil), nil
		}
		existingUser.Username = newUserData.Username
	}

	if newUserData.Email != "" && newUserData.Email != existingUser.Email {
		// validate we dont have a user with the same email that is not deleted
		noCaseEmail := strings.ToLower(newUserData.Email)
		if _, err := s.DB.GetUserByEmail(ctx, s.DBConn, noCaseEmail); err == nil {
			return openapi.Response(http.StatusConflict, nil), nil
		}
		existingUser.Email = newUserData.Email
	}

	if newUserData.FirstName != "" {
		existingUser.FirstName = newUserData.FirstName
	}
	if newUserData.LastName != "" {
		existingUser.LastName = newUserData.LastName
	}

	params := db.UpdateUserByUsernameParams{
		OldUsername: existingUsername,
		NewUsername: newUserData.Username,
		FirstName:   newUserData.FirstName,
		LastName:    newUserData.LastName,
		Email:       newUserData.Email,
	}

	log.Debug().Msgf("UpdateUserByUsernameParams: \n%+v\n", params)
	uUser, err := s.DB.UpdateUserByUsername(ctx, s.DBConn, params)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return openapi.Response(http.StatusNotFound, nil), nil
	}

	return openapi.Response(http.StatusOK, uUser), nil
}
