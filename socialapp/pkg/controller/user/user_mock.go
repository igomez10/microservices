package user

import (
	"context"

	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

// s *UserMock openapi.UserApiServicer
type UserMock struct {
	counter          int
	responseToReturn openapi.ImplResponse
	errorToReturn    error
}

func (s *UserMock) CreateUser(_ context.Context, _ openapi.CreateUserRequest) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) DeleteUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) GetUserByUsername(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) GetUserComments(_ context.Context, _ string, _, _ int32) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) ListUsers(_ context.Context, _, _ int32, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) UpdateUser(_ context.Context, _ string, _ openapi.User) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) Welcome(_ context.Context) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) ResetPassword(_ context.Context, _ openapi.ResetPasswordRequest) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) ChangePassword(_ context.Context, _ openapi.ChangePasswordRequest) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) GetUserFollowers(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) GetFollowingUsers(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) FollowUser(_ context.Context, _, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) UnfollowUser(_ context.Context, _, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) GetRolesForUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) UpdateRolesForUser(_ context.Context, _ string, _ []string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}
