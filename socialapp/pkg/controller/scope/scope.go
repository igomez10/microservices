package scope

import (
	"context"
	"net/http"
	"socialapp/internal/converter"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"

	"github.com/rs/zerolog/log"
)

// s *ScopeApiService openapi.ScopeApiServicer
type ScopeApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *ScopeApiService) CreateScope(ctx context.Context, newScope openapi.Scope) (openapi.ImplResponse, error) {
	// check scope with name doesnt exist
	params := db.CreateScopeParams{
		Name:        newScope.Name,
		Description: newScope.Description,
	}
	res, err := s.DB.CreateScope(ctx, s.DBConn, params)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Str("scope_id", newScope.Name).
			Msg("failed to create scope")

		return openapi.ImplResponse{
			Code: http.StatusConflict,
			Body: openapi.Error{
				Code:    http.StatusConflict,
				Message: "scope already exists",
			},
		}, nil
	}

	scopeID, err := res.LastInsertId()
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Str("scope_name", newScope.Name).
			Msg("failed to retrieve created scope")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to retrieve created scope id",
			},
		}, nil
	}

	scope, err := s.DB.GetScope(ctx, s.DBConn, scopeID)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to retrieve created scope")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to find created scope",
			},
		}, nil
	}

	apiscope := converter.FromDBScopeToApiScope(scope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiscope,
	}, nil
}

func (s *ScopeApiService) DeleteScope(ctx context.Context, scopeID int32) (openapi.ImplResponse, error) {
	//verify scope exists
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to retrieve scope")

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "scope not found",
			},
		}, nil
	}

	deleteErr := s.DB.DeleteScope(ctx, s.DBConn, scope.ID)
	if err != nil {
		log.Error().
			Err(deleteErr).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to retrieve created scope")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to delete scope",
			},
		}, nil
	}

	apiScope := converter.FromDBScopeToApiScope(scope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScope,
	}, nil

}

func (s *ScopeApiService) GetScope(ctx context.Context, scopeID int32) (openapi.ImplResponse, error) {
	s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to retrieve created scope")

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "scope not found",
			},
		}, nil
	}

	apiScope := converter.FromDBScopeToApiScope(scope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScope,
	}, nil
}

func (s *ScopeApiService) ListScopes(ctx context.Context, limit int32, offset int32) (openapi.ImplResponse, error) {
	scopes, err := s.DB.ListScopes(ctx, s.DBConn, db.ListScopesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Msg("failed to retrieve scopes")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to list scopes",
			},
		}, nil
	}

	apiScopes := make([]openapi.Scope, len(scopes))
	for i, scope := range scopes {
		apiScopes[i] = converter.FromDBScopeToApiScope(scope)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScopes,
	}, nil
}

func (s *ScopeApiService) UpdateScope(ctx context.Context, scopeID int32, updatedScope openapi.Scope) (openapi.ImplResponse, error) {
	// get scope from db
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to retrieve scope")

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "scope not found",
			},
		}, nil
	}

	params := db.UpdateScopeParams{
		ID:          scope.ID,
		Name:        updatedScope.Name,
		Description: updatedScope.Description,
	}

	// update scope
	_, err = s.DB.UpdateScope(ctx, s.DBConn, params)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to update scope")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to update scope",
			},
		}, nil
	}

	// get scope again
	scope, err = s.DB.GetScope(ctx, s.DBConn, scope.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Int("scope_id", int(scopeID)).
			Msg("failed to retrieve updated scope")
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to find updated scope",
			},
		}, nil
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: converter.FromDBScopeToApiScope(scope),
	}, nil
}

func (s *ScopeApiService) ListScopesForRole(ctx context.Context, roleID int32) (openapi.ImplResponse, error) {
	scopes, err := s.DB.ListRoleScopes(ctx, s.DBConn, int64(roleID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
			Msg("failed to retrieve scopes")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to list scopes",
			},
		}, nil
	}

	apiScopes := make([]openapi.Scope, len(scopes))
	for i, scope := range scopes {
		apiScopes[i] = converter.FromDBScopeToApiScope(scope)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScopes,
	}, nil
}
