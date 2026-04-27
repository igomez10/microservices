package converter

import (
	"strconv"

	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	db "github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

// formatID renders an int64 database identifier as the JSON-string form used
// on the API surface (see openapi.yaml — int64 fields are serialized as
// strings to avoid precision loss in JS clients).
func formatID(id int64) string {
	return strconv.FormatInt(id, 10)
}

func FromDBCmtToAPICmt(comment db.Comment, user db.User) openapi.Comment {
	cmnt := openapi.Comment{
		Id:        formatID(comment.ID),
		Content:   comment.Content,
		LikeCount: formatID(comment.LikeCount),
		CreatedAt: comment.CreatedAt.Time,
		Username:  user.Username,
	}

	return cmnt
}
func FromPGXDBRoleToAPIRole(dbRole dbpgx.Role) openapi.Role {
	apiRole := openapi.Role{
		Id:          formatID(dbRole.ID),
		Name:        dbRole.Name,
		Description: dbRole.Description,
		CreatedAt:   dbRole.CreatedAt.Time,
	}

	return apiRole
}

func FromDBRoleToAPIRole(dbRole db.Role) openapi.Role {
	apiRole := openapi.Role{
		Id:          formatID(dbRole.ID),
		Name:        dbRole.Name,
		Description: dbRole.Description,
		CreatedAt:   dbRole.CreatedAt.Time,
	}

	return apiRole
}

func FromDBUserToAPIUser(u db.User) openapi.User {
	apiUser := openapi.User{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Time,
	}

	return apiUser
}

func FromDBScopeToAPIScope(dbScope db.Scope) openapi.Scope {
	apiScope := openapi.Scope{
		Id:          formatID(dbScope.ID),
		Name:        dbScope.Name,
		Description: dbScope.Description,
		CreatedAt:   dbScope.CreatedAt.Time,
	}

	return apiScope
}

func FromDBUrlToAPIUrl(dbUrl db.Url) openapi.Url {
	apiUrl := openapi.Url{
		Alias:     dbUrl.Alias,
		Url:       dbUrl.Url,
		CreatedAt: dbUrl.CreatedAt.Time,
		UpdatedAt: dbUrl.UpdatedAt.Time,
	}

	return apiUrl
}
