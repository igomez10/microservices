package comment

import (
	"context"
	"time"

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

func (s *CommentMock) GetComment(_ context.Context, _ int64) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) SearchComments(_ context.Context, _ string, _ time.Time, _ time.Time) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}

func (s *CommentMock) GetUserFeed(_ context.Context) (openapi.ImplResponse, error) {
	s.counter++
	return s.responseToReturn, s.errorToReturn
}
