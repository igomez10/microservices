package role

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/converter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/pkg/snowflake"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// s *RoleApiService openapi.RoleApiServicer

var _ openapi.RoleAPIServicer = (*RoleApiService)(nil)

type RoleApiService struct {
	DB                 dbpgx.Querier
	DBConn             dbpgx.DBTX
	SnowflakeGenerator snowflake.IDGenerator
}

func (s *RoleApiService) CreateRole(ctx context.Context, newRole openapi.Role) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "CreateRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("new_role", fmt.Sprintf("%+v", newRole))

	// check role with name doesn't exist
	roleID, err := s.SnowflakeGenerator.NextID()
	if err != nil {
		logger.Error("Failed to generate snowflake ID for role", "error", err)
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to generate ID",
			},
		}, nil
	}
	params := db.CreateRoleWithIDParams{
		ID:          roleID,
		Name:        newRole.Name,
		Description: newRole.Description,
	}
	createdRole, err := s.DB.CreateRoleWithID(ctx, s.DBConn, params)
	if err != nil {
		logger.Error("Failed to create role", "error", err)

		// Check if it's a duplicate key violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return openapi.ImplResponse{
				Code: http.StatusConflict,
				Body: openapi.Error{
					Code:    http.StatusConflict,
					Message: "Role already exists",
				},
			}, nil
		}

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create role",
			},
		}, nil
	}

	role, err := s.DB.GetRole(ctx, s.DBConn, createdRole.ID)
	if err != nil {
		logger.Error("failed to retrieve created role", "error", err, "role_id", int(createdRole.ID))

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to find created role",
			},
		}, nil
	}

	apiRole := converter.FromDBRoleToAPIRole(role)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRole,
	}, nil
}

func (s *RoleApiService) DeleteRole(ctx context.Context, roleID int64) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "DeleteRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("role_id", int(roleID))
	//verify role exists
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		logger.Error("failed to retrieve role", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	if err := s.DB.DeleteRole(ctx, s.DBConn, role.ID); err != nil {
		logger.Error("failed to retrieve created role", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to delete role",
			},
		}, nil
	}

	apiRole := converter.FromDBRoleToAPIRole(role)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRole,
	}, nil

}

func (s *RoleApiService) GetRole(ctx context.Context, roleID int64) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "GetRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("role_id", int(roleID))

	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			// return 404 if role not found
			return openapi.ImplResponse{
				Code: http.StatusNotFound,
				Body: openapi.Error{
					Code:    http.StatusNotFound,
					Message: "role not found",
				},
			}, nil

		default:
			logger.Error("failed to retrieve role", "error", err)
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Code:    http.StatusInternalServerError,
					Message: "failed to retrieve role",
				},
			}, nil
		}
	}

	apiRole := converter.FromPGXDBRoleToAPIRole(role)

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRole,
	}, nil
}

func (s *RoleApiService) ListRoles(ctx context.Context, limit int32, offset int32) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ListRoles")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("limit", int(limit), "offset", int(offset))

	limit = limit % 20
	if limit == 0 {
		limit = 20
	}

	roles, err := s.DB.ListRoles(ctx, s.DBConn, dbpgx.ListRolesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to retrieve roles", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to list roles",
			},
		}, nil
	}

	apiRoles := make([]openapi.Role, len(roles))
	for i, role := range roles {
		apiRoles[i] = converter.FromPGXDBRoleToAPIRole(role)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRoles,
	}, nil
}

func (s *RoleApiService) UpdateRole(ctx context.Context, roleID int64, newRole openapi.Role) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "UpdateRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With(
		"role_id", int(roleID),
		"new_role", fmt.Sprintf("%+v", newRole),
	)

	// get role from db
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			// return 404 if role not found
			return openapi.ImplResponse{
				Code: http.StatusNotFound,
				Body: openapi.Error{
					Code:    http.StatusNotFound,
					Message: "role not found",
				},
			}, nil
		default:
			logger.Error("failed to retrieve role", "error", err)

			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Code:    http.StatusInternalServerError,
					Message: "failed to retrieve role",
				},
			}, nil
		}
	}

	params := db.UpdateRoleParams{
		ID:          role.ID,
		Name:        newRole.Name,
		Description: newRole.Description,
	}

	// update role
	if err := s.DB.UpdateRole(ctx, s.DBConn, params); err != nil {
		logger.Error("failed to update role", "error", err, "role_id", int(roleID))

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to update role",
			},
		}, nil
	}

	// get role again
	role, err = s.DB.GetRole(ctx, s.DBConn, role.ID)
	if err != nil {
		logger.Error("failed to retrieve updated role", "error", err)
		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to find updated role",
			},
		}, nil
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: converter.FromPGXDBRoleToAPIRole(role),
	}, nil
}

func (s *RoleApiService) AddScopeToRole(ctx context.Context, roleID int64, scopes []string) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "AddScopeToRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("scopes", scopes, "role_id", roleID)

	// get role from db
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		logger.Error("failed to retrieve role", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	// verify all scopes exist
	dbScopes := []db.Scope{}
	for _, scope := range scopes {
		dbSc, err := s.DB.GetScopeByName(ctx, s.DBConn, scope)
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
		dbScopes = append(dbScopes, dbSc)
	}

	// add scopes to role
	for _, sc := range dbScopes {
		// Generate snowflake ID for the role-scope association
		roleScopeID, err := s.SnowflakeGenerator.NextID()
		if err != nil {
			logger.Error("Error generating snowflake ID for role-scope", "error", err)
			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Code:    http.StatusInternalServerError,
					Message: "Failed to generate ID",
				},
			}, nil
		}

		_, err = s.DB.CreateRoleScopeWithID(ctx, s.DBConn, db.CreateRoleScopeWithIDParams{
			ID:      roleScopeID,
			RoleID:  role.ID,
			ScopeID: sc.ID,
		})
		if err != nil {
			logger.Error("failed to add scope to role", "error", err)

			// Check if it's a duplicate key violation
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return openapi.ImplResponse{
					Code: http.StatusConflict,
					Body: openapi.Error{
						Code:    http.StatusConflict,
						Message: "Scope already added to role",
					},
				}, nil
			}

			return openapi.ImplResponse{
				Code: http.StatusInternalServerError,
				Body: openapi.Error{
					Code:    http.StatusInternalServerError,
					Message: "failed to add scope to role",
				},
			}, nil
		}
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: nil,
	}, nil
}

func (s *RoleApiService) ListScopesForRole(ctx context.Context, roleID int64, limit int32, offset int32) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "ListScopesForRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With(
		"role_id", int(roleID),
		"limit", int(limit),
		"offset", int(offset),
	)

	// get role from db
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		logger.Error("failed to retrieve role", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	limit = limit % 20
	if limit == 0 {
		limit = 20
	}
	// get role scopes from db
	scopes, err := s.DB.ListRoleScopes(ctx, s.DBConn, db.ListRoleScopesParams{
		ID:     role.ID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to retrieve role scopes", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to retrieve role scopes",
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

func (s *RoleApiService) RemoveScopeFromRole(ctx context.Context, roleID int64, scopeID int64) (openapi.ImplResponse, error) {
	ctx, span := tracerhelper.GetTracer().Start(ctx, "RemoveScopeFromRole")
	defer span.End()
	logger := contexthelper.GetLoggerInContext(ctx).With("role_id", int(roleID), "scope_id", int(scopeID))

	// verify role exists
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		logger.Error("failed to retrieve role", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	// verify scope exists
	scope, err := s.DB.GetScope(ctx, s.DBConn, int64(scopeID))
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

	// remove scope from role
	params := db.DeleteRoleScopeParams{
		RoleID:  role.ID,
		ScopeID: scope.ID,
	}
	if err := s.DB.DeleteRoleScope(ctx, s.DBConn, params); err != nil {
		logger.Error("failed to remove scope from role", "error", err)

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to remove scope from role",
			},
		}, nil
	}

	return openapi.ImplResponse{
		Code: http.StatusNoContent,
	}, nil
}
