package authorizationparser

import "github.com/getkin/kin-openapi/openapi3"

type EndpointAuthorizations map[string]map[string][]string

const SECURITY_REQUIREMENT_IDENTIFIER = "OAuth2"

func FromOpenAPIToEndpointScopes(doc *openapi3.T) EndpointAuthorizations {
	// response to return
	res := EndpointAuthorizations{}

	// iterate over paths
	for path, item := range doc.Paths {
		// iterate over possible operations (methods of the path)
		for method, operation := range item.Operations() {
			// check if the operation has security requirements
			if operation.Security != nil {
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
						method: {"noauth"}, // "GET": []string{}
					}
				} else {
					res[path][method] = []string{"noauth"}
				}
			}
		}
	}

	return res
}
