package scope

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/snowflake"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// parseAPIID parses an int64 ID that arrived on the API surface as a string
// (see openapi.yaml: int64 fields are serialized as JSON strings to avoid
// JS precision loss). Returns a 400 ApiError payload on parse failure.
func parseAPIID(s string) (int64, *openapi.Error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, &openapi.Error{
			Code:    http.StatusBadRequest,
			Message: "invalid id: must be a numeric string",
		}
	}
	return id, nil
}

// s *ScopeApiService openapi.ScopeApiServicer
var _ openapi.ScopeAPIServicer = (*ScopeApiService)(nil)

type ScopeApiService struct {
	DB                 db.Querier
	DBConn             db.DBTX
	SnowflakeGenerator snowflake.IDGenerator
}

func (s *ScopeApiService) CreateScope(ctx context.Context, newScope openapi.Scope) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateScope")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("new_scope", fmt.Sprintf("%+v", newScope))

	// Generate snowflake ID for the scope
	scopeID, err := s.SnowflakeGenerator.NextID()
	if err != nil {
		logger.Error("Error generating snowflake ID for scope", "error", err)
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to generate ID",
			},
		}, nil
	}

	params := db.CreateScopeWithIDParams{
		ID:          scopeID,
		Name:        newScope.Name,
		Description: newScope.Description,
	}
	createdScope, err := s.DB.CreateScopeWithID(ctx, s.DBConn, params)
	if err != nil {
		logger.Error("failed to create scope", "error", err)

		// Check if it's a duplicate key violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return openapi.ImplResponse{
				Code: http.StatusConflict,
				Body: openapi.Error{
					Code:    http.StatusConflict,
					Message: "Scope already exists",
				},
			}, nil
		}

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create scope",
			},
		}, nil
	}

	apiscope := converter.FromDBScopeToAPIScope(createdScope)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiscope,
	}, nil
}

func (s *ScopeApiService) DeleteScope(ctx context.Context, scopeIDStr string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "DeleteScope")
	defer span.End()
	scopeID, parseErr := parseAPIID(scopeIDStr)
	if parseErr != nil {
		return openapi.Response(int(parseErr.Code), *parseErr), nil
	}
	logger := contexthelper.GetLoggerInContext(ctx).With("scope_id", scopeID)

	//verify scope exists
	scope, err := s.DB.GetScope(ctx, s.DBConn, scopeID)
	if err != nil {
		logger.Error("failed to retrieve scope", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "scope not found",
			},
		}, nil
	}

	deleteErr := s.DB.DeleteScope(ctx, s.DBConn, scope.ID)
	if deleteErr != nil {
		logger.Error("failed to retrieve created scope", "error", deleteErr)

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

func (s *ScopeApiService) GetScope(ctx context.Context, scopeIDStr string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetScope")
	defer span.End()
	scopeID, parseErr := parseAPIID(scopeIDStr)
	if parseErr != nil {
		return openapi.Response(int(parseErr.Code), *parseErr), nil
	}
	logger := contexthelper.GetLoggerInContext(ctx).With("scope_id", scopeID)

	scope, err := s.DB.GetScope(ctx, s.DBConn, scopeID)
	if err != nil {
		logger.Debug("failed to retrieve scope", "error", err)

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
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ListScopes")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("limit", int(limit), "offset", int(offset))

	limit = limit % 20
	if limit == 0 {
		limit = 20
	}

	scopes, err := s.DB.ListScopes(ctx, s.DBConn, db.ListScopesParams{

		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to retrieve scopes", "error", err)

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

func (s *ScopeApiService) UpdateScope(ctx context.Context, scopeIDStr string, updatedScope openapi.Scope) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "UpdateScope")
	defer span.End()
	scopeID, parseErr := parseAPIID(scopeIDStr)
	if parseErr != nil {
		return openapi.Response(int(parseErr.Code), *parseErr), nil
	}
	logger := contexthelper.GetLoggerInContext(ctx).With(
		"scope_id", scopeID,
		"updated_scope", fmt.Sprintf("%+v", updatedScope),
	)

	// get scope from db
	scope, err := s.DB.GetScope(ctx, s.DBConn, scopeID)
	if err != nil {
		logger.Error("failed to retrieve scope", "error", err)

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
		logger.Error("failed to update scope", "error", err)

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
		logger.Error("failed to retrieve updated scope", "error", err)
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
