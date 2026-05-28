package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/igomez10/microservices/socialapp/pkg/llm"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
)

var validIdentifierChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

type UserGenerationInput struct {
	FallbackUsername  string
	FallbackFirstName string
	FallbackLastName  string
	FallbackEmail     string
}

type GeneratedUserData struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
}

type CommentGenerationInput struct {
	FallbackContent  string
	FallbackUsername string
}

type GeneratedCommentData struct {
	Content  string
	Username string
}

type DataGenerator interface {
	GenerateUser(ctx context.Context, in UserGenerationInput) (GeneratedUserData, error)
	GenerateComment(ctx context.Context, in CommentGenerationInput) (GeneratedCommentData, error)
}

type LLMProviderClient struct {
	provider llm.Provider
	model    string
}

func NewLLMProviderClient(provider llm.Provider, model string) *LLMProviderClient {
	return &LLMProviderClient{
		provider: provider,
		model:    model,
	}
}

func (c *LLMProviderClient) Generate(ctx context.Context, prompt string) (string, error) {
	temperature := float32(0.2)
	resp, err := c.provider.Generate(ctx, llm.GenerateRequest{
		Model:            llm.ModelGemini25FlashLite,
		Prompt:           prompt,
		Temperature:      temperature,
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

type LLMTestDataGenerator struct {
	client llm.Provider
}

func StructToMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (g *LLMTestDataGenerator) GenerateUser(ctx context.Context, in UserGenerationInput) (openapi.CreateUserRequest, error) {
	seed := time.Now().UnixNano()
	prompt := fmt.Sprintf(`SEED=%d
Generate realistic testing values. Keep URL-safe and <= 64 chars.
Fallback values:
username=%q
first_name=%q
last_name=%q
email=%q`, seed, in.FallbackUsername, in.FallbackFirstName, in.FallbackLastName, in.FallbackEmail)

	schema, err := StructToMap(openapi.CreateUserRequest{})
	if err != nil {
		return openapi.CreateUserRequest{}, fmt.Errorf("create user request schema: %w", err)
	}

	req := llm.GenerateRequest{
		Model:              llm.ModelGemini25FlashLite,
		Prompt:             prompt,
		Temperature:        0.0,
		ResponseMIMEType:   llm.MimeTypeJSON,
		ResponseJSONSchema: schema,
	}
	resp, err := g.client.Generate(ctx, req)
	if err != nil {
		return openapi.CreateUserRequest{}, err
	}

	parsed := openapi.CreateUserRequest{}
	if err := json.Unmarshal([]byte(resp.Text), &parsed); err != nil {
		return openapi.CreateUserRequest{}, fmt.Errorf("invalid llm user payload: %w", err)
	}

	return parsed, nil
}

func (g *LLMTestDataGenerator) GenerateComment(ctx context.Context, in CommentGenerationInput) (openapi.Create, error) {
	seed := time.Now().UnixNano()
	prompt := fmt.Sprintf(`SEED=%d
Generate realistic testing values. Keep URL-safe and <= 64 chars.
Fallback values:
content=%q
username=%q`, seed, in.FallbackContent, in.FallbackUsername)

	schema, err := StructToMap(openapi.CreateCommentRequest{})
	if err != nil {
		return openapi.CreateCommentRequest{}, fmt.Errorf("create comment request schema: %w", err)
	}

	req := llm.GenerateRequest{
		Model:              llm.ModelGemini25FlashLite,
		Prompt:             prompt,
		Temperature:        0.0,
		ResponseMIMEType:   llm.MimeTypeJSON,
		ResponseJSONSchema: schema,
	}
	resp, err := g.client.Generate(ctx, req)
	if err != nil {
		return openapi.CreateCommentRequest{}, err
	}

	parsed := openapi.CreateCommentRequest{}
	if err := json.Unmarshal([]byte(resp.Text), &parsed); err != nil {
		return openapi.CreateCommentRequest{}, fmt.Errorf("invalid llm comment payload: %w", err)
	}

	return parsed, nil
}
