package converter

import (
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
)

func FromDBCmtToAPICmt(comment db.Comment, user db.User) openapi.Comment {
	cmnt := openapi.Comment{
		Id:        comment.ID,
		Content:   comment.Content,
		LikeCount: int64(comment.LikeCount),
		CreatedAt: comment.CreatedAt,
		Username:  user.Username,
	}

	return cmnt
}

func FromDBRoleToAPIRole(dbRole db.Role) openapi.Role {
	apiRole := openapi.Role{
		Id:          dbRole.ID,
		Name:        dbRole.Name,
		Description: dbRole.Description,
		CreatedAt:   dbRole.CreatedAt,
	}

	return apiRole
}

func FromDBUserToAPIUser(u db.User) openapi.User {
	apiUser := openapi.User{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}

	return apiUser
}

func FromDBScopeToApiScope(dbScope db.Scope) openapi.Scope {
	apiScope := openapi.Scope{
		Id:          dbScope.ID,
		Name:        dbScope.Name,
		Description: dbScope.Description,
		CreatedAt:   dbScope.CreatedAt,
	}

	return apiScope
}
