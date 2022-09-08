package authentication

import (
	"context"
	"socialapp/socialappapi/openapi"
)

// s *AuthenticationMock openapi.AuthenticationApiServicer
type AuthenticationMock struct {
	counter          int
	responseToReturn openapi.ImplResponse
	errorToReturn    error
}

func (s *AuthenticationMock) GetAccessToken(_ context.Context) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}
