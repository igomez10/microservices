package url

import (
	"context"

	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

// s *URLMock openapi.UserApiServicer
type URLMock struct {
	counter          int
	responseToReturn openapi.ImplResponse
	errorToReturn    error
}

func (s *URLMock) CreateUser(_ context.Context, _ openapi.User) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *URLMock) DeleteUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *URLMock) GetUserByUsername(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *URLMock) GetUserComments(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *URLMock) ListUsers(_ context.Context) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *URLMock) UpdateUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}
