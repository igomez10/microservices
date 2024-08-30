package scope

import (
	"context"
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/db"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

// s *ScopeApiService openapi.ScopeApiServicer
type ScopeApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *ScopeApiService) CreateScope(ctx context.Context, newScope openapi.Scope) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ScopeApiService.CreateScope")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	log = log.
		With().
		Str("new_scope", fmt.Sprintf("%+v", newScope)).
		Logger()

	// check scope with name doesnt exist
	params := db.CreateScopeParams{
		Name:        newScope.Name,
		Description: newScope.Description,
	}
	createdScope, err := s.DB.CreateScope(ctx, s.DBConn, params)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to create scope")

		return openapi.ImplResponse{
			Code: http.StatusConflict,
			Body: openapi.Error{
				Code:    http.StatusConflict,
				Message: "scope already exists",
			},
		}, nil
	}

	apiscope := converter.FromDBScopeToAPIScope(createdScope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiscope,
	}, nil
}

func (s *ScopeApiService) DeleteScope(ctx context.Context, scopeID int32) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ScopeApiService.DeleteScope")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	log = log.
		With().
		Int("scope_id", int(scopeID)).
		Logger()

	//verify scope exists
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	if err != nil {
		log.Error().
			Err(err).
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

	apiScope := converter.FromDBScopeToAPIScope(scope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScope,
	}, nil

}

func (s *ScopeApiService) GetScope(ctx context.Context, scopeID int32) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ScopeApiService.GetScope")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	log = log.
		With().
		Int("scope_id", int(scopeID)).
		Logger()

	s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	if err != nil {
		log.Error().
			Err(err).
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

	apiScope := converter.FromDBScopeToAPIScope(scope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScope,
	}, nil
}

func (s *ScopeApiService) ListScopes(ctx context.Context, limit int32, offset int32) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ScopeApiService.ListScopes")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	log = log.
		With().
		Int("limit", int(limit)).
		Int("offset", int(offset)).
		Logger()

	limit = limit % 20
	if limit == 0 {
		limit = 20
	}

	scopes, err := s.DB.ListScopes(ctx, s.DBConn, db.ListScopesParams{

		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Error().
			Err(err).
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
		apiScopes[i] = converter.FromDBScopeToAPIScope(scope)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiScopes,
	}, nil
}

func (s *ScopeApiService) UpdateScope(ctx context.Context, scopeID int32, updatedScope openapi.Scope) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ScopeApiService.UpdateScope")
	defer span.End()
	log := contexthelper.GetLoggerInContext(ctx)
	log = log.
		With().
		Int("scope_id", int(scopeID)).
		Str("updated_scope", fmt.Sprintf("%+v", updatedScope)).
		Logger()

	// get scope from db
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
	if err != nil {
		log.Error().
			Err(err).
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
		Body: converter.FromDBScopeToAPIScope(scope),
	}, nil
}
