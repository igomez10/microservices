package scope

import (
	"context"
	"socialapp/socialappapi/openapi"
)

// s *ScopeMock openapi.ScopeApiServicer
type ScopeMock struct {
	errorToReturn error
}

func (s *ScopeMock) CreateScope(_ context.Context, _ openapi.Scope) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ScopeMock) DeleteScope(_ context.Context, _ int32) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ScopeMock) GetScope(_ context.Context, _ int32) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ScopeMock) ListScopes(_ context.Context, _ int32, _ int32) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *ScopeMock) UpdateScope(_ context.Context, _ int32, _ openapi.Scope) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}
