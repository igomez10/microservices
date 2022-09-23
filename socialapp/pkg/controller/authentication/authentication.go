package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// s *AuthenticationService openapi.AuthenticationApiServicer
type AuthenticationService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *AuthenticationService) GetAccessToken(ctx context.Context) (openapi.ImplResponse, error) {
	username := ctx.Value("username")
	if username == nil {
		return openapi.ImplResponse{
			Code: http.StatusUnauthorized,
			Body: openapi.Error{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			},
		}, nil
	}

	usr, err := s.DB.GetUserByUsername(ctx, s.DBConn, username.(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: []byte(err.Error()),
		}, nil
	}

	uuidToken, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: []byte(err.Error()),
		}, nil
	}

	token := sha256.Sum256([]byte(uuidToken.String()))
	tokenString := base64.URLEncoding.EncodeToString(token[:])
	validUntil := time.Now().Add(time.Hour * 6)

	_, err = s.DB.CreateToken(ctx, s.DBConn, db.CreateTokenParams{
		Token:      tokenString,
		UserID:     usr.ID,
		ValidUntil: validUntil,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: []byte(err.Error()),
		}, nil
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: openapi.AccessToken{
			AccessToken: tokenString,
			Scopes:      []string{"read", "write"},
			ExpiresAt:   validUntil,
		},
	}, nil
}
