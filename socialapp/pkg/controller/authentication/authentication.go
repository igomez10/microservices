package authentication

import (
	"context"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
)

// s *AuthenticationService openapi.AuthenticationApiServicer
type AuthenticationService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *AuthenticationService) GetAccessToken(ctx context.Context) (openapi.ImplResponse, error) {
	panic("not implemented") // TODO: Implement
}
