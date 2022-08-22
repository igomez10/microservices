package comment

import (
	"context"
	"socialapp/socialappapi/openapi"
)

// s *CommentMock openapi.UserApiServicer
type CommentMock struct {
	counter          int
	responseToReturn openapi.ImplResponse
	errorToReturn    error
}

func (s *CommentMock) CreateUser(_ context.Context, _ openapi.User) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) DeleteUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) GetUserByUsername(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) GetUserComments(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) ListUsers(_ context.Context) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) UpdateUser(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}
