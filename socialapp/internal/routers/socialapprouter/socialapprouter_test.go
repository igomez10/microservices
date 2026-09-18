package socialapprouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igomez10/microservices/socialapp/internal/authorizationparser"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

type testRouter struct {
	routes openapi.Routes
}

func (r testRouter) Routes() openapi.Routes {
	return r.routes
}

func (r testRouter) OrderedRoutes() []openapi.Route {
	routes := make([]openapi.Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}
	return routes
}

func anonymousScopes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, contexthelper.SetRequestedScopesInContext(r, map[string]bool{}))
	})
}

func TestNewSocialAppRouterAuthorizationPolicy(t *testing.T) {
	tests := []struct {
		name       string
		route      openapi.Route
		policies   authorizationparser.EndpointAuthorizations
		wantStatus int
		wantError  bool
	}{
		{
			name: "protected route rejects anonymous request",
			route: openapi.Route{
				Name:        "ListUsers",
				Method:      http.MethodGet,
				Pattern:     "/api/v1/users",
				HandlerFunc: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) },
			},
			policies: authorizationparser.EndpointAuthorizations{
				"/v1/users": {http.MethodGet: {"socialapp.users.list"}},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "explicit public route permits anonymous request",
			route: openapi.Route{
				Name:        "CreateUser",
				Method:      http.MethodPost,
				Pattern:     "/api/v1/users",
				HandlerFunc: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) },
			},
			policies: authorizationparser.EndpointAuthorizations{
				"/v1/users": {http.MethodPost: {}},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "missing policy fails closed at startup",
			route: openapi.Route{
				Name:        "ListUsers",
				Method:      http.MethodGet,
				Pattern:     "/api/v1/users",
				HandlerFunc: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) },
			},
			policies:  authorizationparser.EndpointAuthorizations{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, err := NewSocialAppRouter(
				[]func(http.Handler) http.Handler{anonymousScopes},
				[]openapi.Router{testRouter{routes: openapi.Routes{"test": tt.route}}},
				tt.policies,
			)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected missing authorization policy to return an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewSocialAppRouter() error = %v", err)
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.route.Method, tt.route.Pattern, nil)
			router.Router.ServeHTTP(recorder, request)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
		})
	}
}
