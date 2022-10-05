package role

import (
	"context"
	"net/http"
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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

	apiRole := FromDBRoleToApiRole(role)
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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

	apiRole := FromDBRoleToApiRole(role)
	return openapi.ImplResponse{
		Code: http.StatusOK,
		Body: apiRole,
	}, nil

}

func (s *RoleApiService) GetRole(ctx context.Context, roleID int32) (openapi.ImplResponse, error) {
	s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	role, err := s.DB.GetRole(ctx, s.DBConn, int64(roleID))
	if err != nil {
		log.Error().
			Err(err).
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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

	apiRole := FromDBRoleToApiRole(role)
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
		apiRoles[i] = FromDBRoleToApiRole(role)
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
			Str("X-Request-ID", ctx.Value("X-Request-ID").(string)).
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
		Body: FromDBRoleToApiRole(role),
	}, nil
}

func FromDBRoleToApiRole(dbRole db.Role) openapi.Role {
	apiRole := openapi.Role{
		Id:          dbRole.ID,
		Name:        dbRole.Name,
		Description: dbRole.Description,
		CreatedAt:   dbRole.CreatedAt,
	}
	return apiRole
}
