package main

import (
	"context"
	"errors"
	"os"
	"testing"
)

type mockGeminiClient struct {
	response string
	err      error
}

func (m *mockGeminiClient) Generate(_ context.Context, _ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func TestNewTestDataGeneratorFromEnv_DefaultWhenNoAPIKey(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GEMINI_MODEL", "")

	gen := NewTestDataGeneratorFromEnv()
	if _, ok := gen.(*DefaultTestDataGenerator); !ok {
		t.Fatalf("expected DefaultTestDataGenerator, got %T", gen)
	}
}

func TestDefaultGeneratorPreservesFallbackValues(t *testing.T) {
	gen := &DefaultTestDataGenerator{}
	user, err := gen.GenerateUser(context.Background(), UserGenerationInput{
		FallbackUsername:  "Test-1231",
		FallbackFirstName: "FirstName_example",
		FallbackLastName:  "LastName_example",
		FallbackEmail:     "Test-1231@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "Test-1231" || user.FirstName != "FirstName_example" || user.LastName != "LastName_example" || user.Email != "Test-1231@example.com" {
		t.Fatalf("default user generator did not preserve fallback values: %+v", user)
	}

	comment, err := gen.GenerateComment(context.Background(), CommentGenerationInput{
		FallbackContent:  "search-first",
		FallbackUsername: "Test-1231",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Content != "search-first" || comment.Username != "Test-1231" {
		t.Fatalf("default comment generator did not preserve fallback values: %+v", comment)
	}
}

func TestGeminiGeneratorParsesJSON(t *testing.T) {
	gen := NewGeminiTestDataGenerator(&mockGeminiClient{
		response: `{"username":"llm-user","first_name":"Ada","last_name":"Lovelace","email":"ada@example.com"}`,
	})
	user, err := gen.GenerateUser(context.Background(), UserGenerationInput{
		FallbackUsername:  "fallback-user",
		FallbackFirstName: "Fallback",
		FallbackLastName:  "User",
		FallbackEmail:     "fallback@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "llm-user" || user.FirstName != "Ada" || user.LastName != "Lovelace" || user.Email != "ada@example.com" {
		t.Fatalf("unexpected generated user: %+v", user)
	}
}

func TestGeminiGeneratorParsesFencedJSON(t *testing.T) {
	gen := NewGeminiTestDataGenerator(&mockGeminiClient{
		response: "```json\n{\"content\":\"hello from llm\",\"username\":\"llm-user\"}\n```",
	})
	comment, err := gen.GenerateComment(context.Background(), CommentGenerationInput{
		FallbackContent:  "fallback-comment",
		FallbackUsername: "fallback-user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Content != "hello from llm" || comment.Username != "llm-user" {
		t.Fatalf("unexpected generated comment: %+v", comment)
	}
}

func TestFallbackGeneratorUsesDefaultOnPrimaryFailure(t *testing.T) {
	primary := NewGeminiTestDataGenerator(&mockGeminiClient{err: errors.New("boom")})
	fallback := &DefaultTestDataGenerator{}
	gen := NewFallbackTestDataGenerator(primary, fallback)

	user, err := gen.GenerateUser(context.Background(), UserGenerationInput{
		FallbackUsername:  "fallback-user",
		FallbackFirstName: "Fallback",
		FallbackLastName:  "User",
		FallbackEmail:     "fallback@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "fallback-user" {
		t.Fatalf("expected fallback username, got %q", user.Username)
	}
}

func TestNewTestDataGeneratorFromEnv_DefaultOnClientInitFailure(t *testing.T) {
	prev := os.Getenv("GEMINI_API_KEY")
	defer os.Setenv("GEMINI_API_KEY", prev)
	t.Setenv("GEMINI_API_KEY", " ")
	gen := NewTestDataGeneratorFromEnv()
	if _, ok := gen.(*DefaultTestDataGenerator); !ok {
		t.Fatalf("expected default generator when key is blank, got %T", gen)
	}
}
