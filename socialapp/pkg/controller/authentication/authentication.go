package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/pkg/db"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

// s *AuthenticationService openapi.AuthenticationApiServicer
type AuthenticationService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *AuthenticationService) GetAccessToken(ctx context.Context) (openapi.ImplResponse, error) {
	log := contexthelper.GetLoggerInContext(ctx)
	username, ok := contexthelper.GetUsernameInContext(ctx)
	if !ok {
		log.Error().Str("username", username).Msg("username not found in context")
		return openapi.ImplResponse{
			Code: http.StatusUnauthorized,
			Body: openapi.Error{
				Code:    http.StatusUnauthorized,
				Message: fmt.Errorf("failed to resolve username").Error(),
			},
		}, nil
	}

	requestedScopes, ok := contexthelper.GetRequestedScopesInContext(ctx)
	if !ok {
		log.Error().Interface("scopes", requestedScopes).Msg("scopes not found in context")
		return openapi.ImplResponse{
			Code: http.StatusUnauthorized,
			Body: openapi.Error{
				Code:    http.StatusUnauthorized,
				Message: fmt.Errorf("failed to resolve scopes").Error(),
			},
		}, nil
	}

	usr, err := s.DB.GetUserByUsername(ctx, s.DBConn, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: fmt.Errorf("failed to get user").Error(),
			},
		}, nil
	}

	uuidToken, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate token")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Errorf("failed to generate token").Error(),
			},
		}, nil
	}

	token := sha256.Sum256([]byte(uuidToken.String()))
	tokenString := base64.URLEncoding.EncodeToString(token[:])
	validUntil := time.Now().UTC().Add(30 * 24 * time.Hour)

	createdToken, err := s.DB.CreateToken(ctx, s.DBConn, db.CreateTokenParams{
		Token:      tokenString,
		UserID:     usr.ID,
		ValidUntil: validUntil,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Errorf("failed to create token").Error(),
			},
		}, nil
	}

	for scopeName := range requestedScopes {
		// get id of the scope
		scope, err := s.DB.GetScopeByName(ctx, s.DBConn, scopeName)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get scope")
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Code:    http.StatusInternalServerError,
					Message: fmt.Errorf("failed to get scope").Error(),
				},
			}, nil
		}

		// create token scope
		_, err = s.DB.CreateTokenToScope(ctx, s.DBConn, db.CreateTokenToScopeParams{
			TokenID: createdToken.ID,
			ScopeID: scope.ID,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create token scope")
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Code:    http.StatusInternalServerError,
					Message: fmt.Errorf("failed to create token scope").Error(),
				},
			}, nil
		}
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Errorf("failed to create token").Error(),
			},
		}, nil
	}

	assignedScopes := make([]string, 0, len(requestedScopes))
	for scopeName := range requestedScopes {
		assignedScopes = append(assignedScopes, scopeName)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: openapi.AccessToken{
			AccessToken: tokenString,
			Scopes:      assignedScopes,
			ExpiresIn:   int32(time.Until(validUntil).Seconds()),
			TokenType:   "Bearer",
		},
	}, nil
}
