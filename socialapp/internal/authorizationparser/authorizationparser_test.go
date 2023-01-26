package authorizationparser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
)

func TestParseAuthorization(t *testing.T) {
	specBytes, err := ioutil.ReadAll(bytes.NewBuffer(testAPI))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := openapi3.NewLoader().LoadFromData(specBytes)

	if err != nil {
		t.Fatal(err)
	}

	res := FromOpenAPIToEndpointScopes(doc)
	fmt.Printf("%+v", res)

	if diff := cmp.Diff(expectedEndpointAuthorizationsForTestApi, res); diff != "" {
		t.Error(diff)
	}
}

// FromOpenAPIToEndpointScopes parses the OpenAPI specification and returns a map of endpoint to scopes
// that are required to access the endpoint. If the endpoint does not require any scopes, the map will
// not contain the endpoint. If the endpoint requires multiple scopes, the map will contain the endpoint
// with a list of scopes [array], see example.
// example reponse:
//
//	{
//		"/users": {
//			"GET": ["socialapp.users.list", "socialapp.users.list"],
//			"POST": ["socialapp.users.create"]
//		},
//		"/users/{username}": {
//			"GET": ["socialapp.users.get"],
//			"PUT": ["socialapp.users.update"],
//			"DELETE": ["socialapp.users.delete"]
//		}
//	}

var expectedEndpointAuthorizationsForTestApi = EndpointAuthorizations{
	"/users": {
		"GET":  []string{"socialapp.users.list"},
		"POST": []string{"noauth"},
	},
	"/users/{username}": {
		"GET":    []string{"socialapp.users.read"},
		"PUT":    []string{"socialapp.users.update"},
		"DELETE": []string{"socialapp.users.delete"},
	},
}

// example openapi.yaml
var testAPI []byte = []byte(`
openapi: 3.0.0
info:
  version: "1.0.0"
  title: "Socialapp"
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT
  description: Socialapp is a generic social network.
  contact:
    name: "Socialapp"
    url: "https://microservices.onrender.com"
    email: "ignacio.gomez.arboleda@gmail.com"
servers:
  - url: https://microservices.onrender.com
  - url: http://localhost:8080
  - url: http://localhost:8085

tags:
  - name: User
    description: User management
  - name: Comment
    description: Comment management
  - name: Authentication
    description: Authentication management
  - name: Following
    description: Following management

paths:
  /users:
    get:
      summary: "List users"
      description: "List all users in the system (paginated)"
      operationId: listUsers
      security:
        - OAuth2: [socialapp.users.list]
      tags:
        - User
      parameters:
        - name: limit
          in: query
          description: "Maximum number of users to return"
          required: false
          schema:
            type: integer
            format: int32
            default: 20
        - name: offset
          in: query
          description: "Pagination offset"
          required: false
          schema:
            type: integer
            format: int32
            default: 0
      responses:
        "200":
          description: List of all the users
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/User"
    post:
      summary: Create user
      description: Create a new user in the system
      operationId: createUser
      tags:
        - User
      requestBody:
        description: Create a new user
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/CreateUserRequest"
        required: true
      responses:
        "200":
          description: User was created successfully
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/User"
  /users/{username}:
    delete:
      summary: Deletes a particular user
      description: Deletes a particular user by username
      operationId: deleteUser
      tags:
        - User
      security:
        - OAuth2: [socialapp.users.delete]
      parameters:
        - name: username
          in: path
          description: username of the user
          required: true
          schema:
            type: string
          example: "johndoe"
      responses:
        "200":
          description: User was deleted succesfully
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/User"
    get:
      summary: Get a particular user by username
      description: Get a particular user by username
      operationId: getUserByUsername
      tags:
        - User
      security:
        - OAuth2: [socialapp.users.read]
      parameters:
        - name: username
          in: path
          description: username of the user
          required: true
          schema:
            type: string
          example: "johndoe"
      responses:
        "200":
          description: Details about a user by ID
          headers:
            x-next:
              description: A link to the next page of responses
              schema:
                type: string
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/User"
    put:
      summary: Update a user
      description: Update a user by username
      operationId: updateUser
      tags:
        - User
      security:
        - OAuth2: [socialapp.users.update]
      parameters:
        - name: username
          in: path
          description: username of the user
          required: true
          schema:
            type: string
          example: "johndoe"
      requestBody:
        description: Update a user
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/User"
      responses:
        "200":
          description: User was updated successfully
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/User"

components:
  schemas:
    CreateUserRequest:
      type: object
      required:
        - username
        - first_name
        - last_name
        - email
        - password
      example:
        username: "johndoe"
        first_name: "John"
        last_name: "Doe"
        password: "Secure123!"
        email: "johndoe@mail.com"
      properties:
        id:
          type: integer
          format: int64
        username:
          type: string
        password:
          type: string
        first_name:
          type: string
        last_name:
          type: string
        email:
          type: string
        created_at:
          type: string
          format: date-time
    User:
      type: object
      required:
        - username
        - first_name
        - last_name
        - email
      example:
        username: "johndoe"
        first_name: "John"
        last_name: "Doe"
        email: "johndoe@mail.com"
      properties:
        id:
          type: integer
          format: int64
        username:
          type: string
        first_name:
          type: string
        last_name:
          type: string
        email:
          type: string
        created_at:
          type: string
          format: date-time

    Comment:
      type: object
      required:
        - content
        - username
      example:
        content: "This is a comment"
        username: "johndoe"
      properties:
        id:
          type: integer
          format: int64
        content:
          type: string
        like_count:
          type: integer
          format: int64
        created_at:
          type: string
          format: date-time
        username:
          type: string
`)
