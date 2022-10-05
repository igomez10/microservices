package role

import (
	"context"
	"socialapp/socialappapi/openapi"
)

// s *RoleMock openapi.RoleApiServicer
type RoleMock struct {
	errorToReturn error
}

func (s *RoleMock) CreateRole(_ context.Context, _ openapi.Role) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *RoleMock) DeleteRole(_ context.Context, _ string, _ openapi.Role) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *RoleMock) GetRole(_ context.Context, _ string) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *RoleMock) ListRoles(_ context.Context, _ int32, _ int32) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (s *RoleMock) UpdateRole(_ context.Context, _ string, _ openapi.Role) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}
