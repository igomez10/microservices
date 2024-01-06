package url

import (
	"context"

	"github.com/igomez10/microservices/urlshortener/generated/server"
)

// validate UrlMock implements the URLAPIServicer interface
var _ server.URLAPIServicer = (*URLMock)(nil)

// s *URLMock server.URLAPIServicer
type URLMock struct {
	counter          int
	responseToReturn server.ImplResponse
	errorToReturn    error
}

func (s *URLMock) CreateUrl(_ context.Context, _ server.Url) (server.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *URLMock) DeleteUrl(_ context.Context, _ string) (server.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *URLMock) GetUrl(_ context.Context, _ string) (server.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *URLMock) GetUrlData(_ context.Context, _ string) (server.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}
