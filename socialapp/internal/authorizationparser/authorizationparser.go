package authorizationparser

import "github.com/getkin/kin-openapi/openapi3"

type EndpointAuthorizations map[string]map[string][]string

const SECURITY_REQUIREMENT_IDENTIFIER = "OAuth2"

func FromOpenAPIToEndpointScopes(doc *openapi3.T) EndpointAuthorizations {
	// response to return
	res := EndpointAuthorizations{}

	// iterate over paths
	for path, item := range doc.Paths.Map() {
		// iterate over possible operations (methods of the path)
		for method, operation := range item.Operations() {
			// An explicit empty security array means the operation is public.
			// Keep an entry with zero scopes so callers can distinguish that
			// deliberate policy from missing path or method metadata.
			if operation.Security != nil && len(*operation.Security) == 0 {
				if _, exist := res[path]; !exist {
					res[path] = map[string][]string{method: {}}
				} else {
					res[path][method] = []string{}
				}
			} else if operation.Security != nil {
				// iterate over security requirements of the operation
				for _, secReq := range *operation.Security {
					requiredScopes := secReq[SECURITY_REQUIREMENT_IDENTIFIER] // []string{"socialapp.users.read", "socialapp.users.write"}
					if _, exist := res[path]; !exist {
						res[path] = map[string][]string{ // "/user"
							method: requiredScopes, // "GET": []string{"socialapp.users.read", "socialapp.users.write"}
						}
					} else {
						res[path][method] = requiredScopes
					}
				}
			} else {
				// if the operation has no security requirements, we add an empty array of scopes
				if _, exist := res[path]; !exist {
					res[path] = map[string][]string{ // "/user"
						method: {},
					}
				} else {
					res[path][method] = []string{}
				}
			}
		}
	}
	return res
}
