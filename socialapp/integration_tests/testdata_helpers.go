package main

import (
	"context"

	"github.com/igomez10/microservices/socialapp/client"
)

var testDataGenerator = NewTestDataGeneratorFromEnv()

func generatedUser(
	ctx context.Context,
	fallbackUsername string,
	password string,
	fallbackFirstName string,
	fallbackLastName string,
	fallbackEmail string,
) (string, *client.CreateUserRequest, error) {
	data, err := testDataGenerator.GenerateUser(ctx, UserGenerationInput{
		FallbackUsername:  fallbackUsername,
		FallbackFirstName: fallbackFirstName,
		FallbackLastName:  fallbackLastName,
		FallbackEmail:     fallbackEmail,
	})
	if err != nil {
		return "", nil, err
	}
	req := client.NewCreateUserRequest(data.Username, password, data.FirstName, data.LastName, data.Email)
	return data.Username, req, nil
}

func generatedComment(ctx context.Context, fallbackContent, fallbackUsername string) (*client.CreateCommentRequest, error) {
	data, err := testDataGenerator.GenerateComment(ctx, CommentGenerationInput{
		FallbackContent:  fallbackContent,
		FallbackUsername: fallbackUsername,
	})
	if err != nil {
		return nil, err
	}
	return client.NewCreateCommentRequest(data.Content, data.Username), nil
}
