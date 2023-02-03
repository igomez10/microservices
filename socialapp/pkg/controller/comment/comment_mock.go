package comment

import (
	"context"

	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

// s *CommentMock openapi.CommentApiServicer
type CommentMock struct {
	counter          int
	responseToReturn openapi.ImplResponse
	errorToReturn    error
}

func (s *CommentMock) CreateComment(_ context.Context, _ openapi.Comment) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) GetComment(_ context.Context, _ int32) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) GetUserComments(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) GetUserFeed(_ context.Context, _ string) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}
