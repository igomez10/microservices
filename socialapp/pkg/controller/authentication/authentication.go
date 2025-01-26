package authentication

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	socialappjwt "github.com/igomez10/microservices/socialapp/internal/jwt"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

var _ openapi.AuthenticationAPIServicer = (*AuthenticationService)(nil)

// s *AuthenticationService openapi.AuthenticationApiServicer
type AuthenticationService struct {
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
