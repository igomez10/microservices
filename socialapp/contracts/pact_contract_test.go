package contracts

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"

	client "github.com/igomez10/microservices/socialapp/client"
)

func TestUserAPI_CreateUserContract(t *testing.T) {
	runConsumerContract(t, "socialapp-client", "socialapp", contractSpec{
		providerState:  "a user can be created",
		description:    "a request to create a user",
		requestMethod:  consumer.Method(http.MethodPost),
		requestPath:    "/v1/users",
		responseStatus: http.StatusOK,
		requestBuilder: func(r *consumer.V2RequestBuilder) {
			r.Header("Content-Type", matchers.String("application/json"))
			r.JSONBody(matchers.StructMatcher{
				"username":   matchers.String("pact-user"),
				"password":   matchers.String("pact-secret"),
				"first_name": matchers.String("Pact"),
				"last_name":  matchers.String("Tester"),
				"email":      matchers.Term("pact@example.com", ".+@.+\\..+"),
			})
		},
		responseBuilder: func(r *consumer.V2ResponseBuilder) {
			r.Header("Content-Type", matchers.String("application/json"))
			r.JSONBody(matchers.StructMatcher{
				"id":         matchers.Identifier(),
				"username":   matchers.String("pact-user"),
				"first_name": matchers.String("Pact"),
				"last_name":  matchers.String("Tester"),
				"email":      matchers.String("pact@example.com"),
				"created_at": matchers.Timestamp(),
			})
		},
		executionHandler: executeCreateUser,
	})
}

func TestCommentAPI_CreateCommentContract(t *testing.T) {
	runConsumerContract(t, "socialapp-client", "socialapp", contractSpec{
		providerState:  "a comment can be posted",
		description:    "a request to create a comment",
		requestMethod:  consumer.Method(http.MethodPost),
		requestPath:    "/v1/comments",
		responseStatus: http.StatusOK,
		requestBuilder: func(r *consumer.V2RequestBuilder) {
			r.Header("Content-Type", matchers.String("application/json"))
			r.JSONBody(matchers.StructMatcher{
				"username": matchers.String("pact-user"),
				"content":  matchers.String("Hello from Pact"),
			})
		},
		responseBuilder: func(r *consumer.V2ResponseBuilder) {
			r.Header("Content-Type", matchers.String("application/json"))
			r.JSONBody(matchers.StructMatcher{
				"id":         matchers.Identifier(),
				"username":   matchers.String("pact-user"),
				"content":    matchers.String("Hello from Pact"),
				"created_at": matchers.Timestamp(),
			})
		},
		executionHandler: executeCreateComment,
	})
}

func TestScopeAPI_CreateScopeContract(t *testing.T) {
	runConsumerContract(t, "socialapp-client", "socialapp", contractSpec{
		providerState:  "a scope can be created",
		description:    "a request to create a scope",
		requestMethod:  consumer.Method(http.MethodPost),
		requestPath:    "/v1/scopes",
		responseStatus: http.StatusOK,
		requestBuilder: func(r *consumer.V2RequestBuilder) {
			r.Header("Content-Type", matchers.String("application/json"))
			r.JSONBody(matchers.StructMatcher{
				"name":        matchers.String("socialapp.read"),
				"description": matchers.String("Read data from Socialapp"),
			})
		},
		responseBuilder: func(r *consumer.V2ResponseBuilder) {
			r.Header("Content-Type", matchers.String("application/json"))
			r.JSONBody(matchers.StructMatcher{
				"id":          matchers.Identifier(),
				"name":        matchers.String("socialapp.read"),
				"description": matchers.String("Read data from Socialapp"),
				"created_at":  matchers.Timestamp(),
			})
		},
		executionHandler: executeCreateScope,
	})
}

func executeCreateUser(config consumer.MockServerConfig) error {
	apiClient := newSocialappClient(config)
	req := client.NewCreateUserRequest("pact-user", "pact-secret", "Pact", "Tester", "pact@example.com")
	_, _, err := apiClient.UserAPI.CreateUser(context.Background()).CreateUserRequest(*req).Execute()
	return err
}

func executeCreateComment(config consumer.MockServerConfig) error {
	apiClient := newSocialappClient(config)
	comment := client.Comment{
		Username: "pact-user",
		Content:  "Hello from Pact",
	}
	_, _, err := apiClient.CommentAPI.CreateComment(context.Background()).Comment(comment).Execute()
	return err
}

func executeCreateScope(config consumer.MockServerConfig) error {
	apiClient := newSocialappClient(config)
	scope := client.Scope{
		Name:        "socialapp.read",
		Description: "Read data from Socialapp",
	}
	_, _, err := apiClient.ScopeAPI.CreateScope(context.Background()).Scope(scope).Execute()
	return err
}

func newSocialappClient(config consumer.MockServerConfig) *client.APIClient {
	apiConfig := client.NewConfiguration()
	apiConfig.Servers = client.ServerConfigurations{
		{URL: fmt.Sprintf("http://%s:%d", config.Host, config.Port)},
	}
	return client.NewAPIClient(apiConfig)
}
