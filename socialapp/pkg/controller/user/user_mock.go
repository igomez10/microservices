package user

import (
	"context"
	"socialapp/socialappapi/openapi"
)

// s *UserMock openapi.UserApiServicer
type UserMock struct {
	counter          int
	responseToReturn openapi.ImplResponse
	errorToReturn    error
}

func (s *UserMock) CreateUser(_ context.Context, _ openapi.User) (openapi.ImplResponse, error) {
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

func (s *UserMock) GetUserComments(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) ListUsers(_ context.Context) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *UserMock) UpdateUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}
