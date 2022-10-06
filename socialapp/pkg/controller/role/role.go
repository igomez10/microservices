package role

import (
	"context"
	"net/http"
	"socialapp/internal/contexthelper"
	"socialapp/internal/converter"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"

	"github.com/rs/zerolog/log"
)

// s *RoleApiService openapi.RoleApiServicer
type RoleApiService struct {
	DB     db.Querier
	DBConn db.DBTX
}

func (s *RoleApiService) CreateRole(ctx context.Context, newRole openapi.Role) (openapi.ImplResponse, error) {
	// check role with name doesnt exist
	params := db.CreateRoleParams{
		Name:        newRole.Name,
		Description: newRole.Description,
	}
	res, err := s.DB.CreateRole(ctx, s.DBConn, params)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Str("role_id", newRole.Name).
			Msg("failed to create role")

		return openapi.ImplResponse{
			Code: http.StatusConflict,
			Body: openapi.Error{
				Code:    http.StatusConflict,
				Message: "role already exists",
			},
		}, nil
	}

	roleID, err := res.LastInsertId()
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Str("role_name", newRole.Name).
			Msg("failed to retrieve created role")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to retrieve created role id",
			},
		}, nil
	}

	role, err := s.DB.GetRole(ctx, s.DBConn, roleID)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve created role")

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

func (s *RoleApiService) DeleteRole(ctx context.Context, roleID int32) (openapi.ImplResponse, error) {
	//verify role exists
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve role")

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	if err := s.DB.DeleteRole(ctx, s.DBConn, role.ID); err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve created role")

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

func (s *RoleApiService) GetRole(ctx context.Context, roleID int32) (openapi.ImplResponse, error) {
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve created role")

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	// get scopes for role
	dbScopes, err := s.DB.ListRoleScopes(ctx, s.DBConn, role.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve scopes for role")

		return openapi.ImplResponse{
			Code: http.StatusInternalServerError,
			Body: openapi.Error{
				Code:    http.StatusInternalServerError,
				Message: "failed to retrieve scopes for role",
			},
		}, nil
	}

	apiRole := converter.FromDBRoleToAPIRole(role)
	for i := range dbScopes {
		apiRole.Scopes = append(apiRole.Scopes, converter.FromDBScopeToAPIScope(dbScopes[i]))
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRole,
	}, nil
}

func (s *RoleApiService) ListRoles(ctx context.Context, limit int32, offset int32) (openapi.ImplResponse, error) {
	roles, err := s.DB.ListRoles(ctx, s.DBConn, db.ListRolesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Msg("failed to retrieve roles")

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
		apiRoles[i] = converter.FromDBRoleToAPIRole(role)
	}

	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRoles,
	}, nil
}

func (s *RoleApiService) UpdateRole(ctx context.Context, roleID int32, newRole openapi.Role) (openapi.ImplResponse, error) {
	// get role from db
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve role")

		return openapi.ImplResponse{
			Code: http.StatusNotFound,
			Body: openapi.Error{
				Code:    http.StatusNotFound,
				Message: "role not found",
			},
		}, nil
	}

	params := db.UpdateRoleParams{
		ID:          role.ID,
		Name:        newRole.Name,
		Description: newRole.Description,
	}

	// update role
	_, err = s.DB.UpdateRole(ctx, s.DBConn, params)
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to update role")

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
		log.Error().
			Err(err).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(ctx)).
			Int("role_id", int(roleID)).
			Msg("failed to retrieve updated role")
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
		Body: converter.FromDBRoleToAPIRole(role),
	}, nil
}
