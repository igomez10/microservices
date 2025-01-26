package authentication

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	socialappjwt "github.com/igomez10/microservices/socialapp/internal/jwt"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ openapi.AuthenticationAPIServicer = (*AuthenticationService)(nil)

// s *AuthenticationService openapi.AuthenticationApiServicer
type AuthenticationService struct {
	DB        db.Querier
	DBConn    *pgxpool.Pool
	JWTSecret string
}

func (s *AuthenticationService) GetAccessToken(ctx context.Context) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetAccessToken")
	defer span.End()
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

	jwtEnabled := true
	if jwtEnabled {
		scopes := make([]string, 0, len(requestedScopes))
		for scope := range requestedScopes {
			scopes = append(scopes, scope)
		}

		newtoken := socialappjwt.SocialAPPToken{
			Username:  username,
			Scopes:    scopes,
			Audience:  jwt.ClaimStrings{"socialapp"},
			Issuer:    "socialapp",
			Expires:   jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "socialapp",
		}

		stringToken, err := socialappjwt.TokenToString(&newtoken, []byte(s.JWTSecret))
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

		return openapi.ImplResponse{
			Code: http.StatusOK,
			Body: openapi.AccessToken{
				AccessToken: stringToken,
				Scopes:      scopes,
				ExpiresIn:   int32(time.Until(newtoken.Expires.Time).Seconds()),
				TokenType:   "Bearer",
			},
		}, nil
	}

	// begin transaction
	tx, err := s.DBConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
		}, nil
	}
	defer tx.Rollback(ctx)

	usr, err := s.DB.GetUserByUsername(ctx, tx, username)
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

	createdToken, err := s.DB.CreateToken(ctx, tx, db.CreateTokenParams{
		Token:  tokenString,
		UserID: usr.ID,
		ValidUntil: pgtype.Timestamp{
			Time:  validUntil,
			Valid: true,
		},
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
		scope, err := s.DB.GetScopeByName(ctx, tx, scopeName)
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
		_, err = s.DB.CreateTokenToScope(ctx, tx, db.CreateTokenToScopeParams{
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

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to commit transaction")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
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
