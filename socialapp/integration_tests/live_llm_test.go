package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLiveLLMGeneratesUserAndCommentFields(t *testing.T) {
	if os.Getenv("RUN_LIVE_LLM_TESTS") != "1" {
		t.Skip("set RUN_LIVE_LLM_TESTS=1 to run live LLM generation test")
	}

	apiKey := resolveLLMAPIKeyFromEnv()
	if apiKey == "" {
		t.Skip("no LLM API key found in environment")
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = defaultGeminiModel
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := NewGeminiClient(ctx, GeminiClientConfig{
		APIKey: apiKey,
		Model:  model,
	})
	if err != nil {
		t.Fatalf("failed to create gemini client: %v", err)
	}

	generator := NewGeminiTestDataGenerator(client)

	userData, err := generator.GenerateUser(ctx, UserGenerationInput{
		FallbackUsername:  "test-live-user-fallback",
		FallbackFirstName: "FirstName_example",
		FallbackLastName:  "LastName_example",
		FallbackEmail:     "test-live-user-fallback@example.com",
	})
	if err != nil {
		t.Fatalf("failed to generate user fields from llm: %v", err)
	}

	if userData.Username == "" || userData.FirstName == "" || userData.LastName == "" || userData.Email == "" {
		t.Fatalf("expected all generated user fields to be non-empty, got: %+v", userData)
	}

	commentData, err := generator.GenerateComment(ctx, CommentGenerationInput{
		FallbackContent:  "fallback-comment",
		FallbackUsername: userData.Username,
	})
	if err != nil {
		t.Fatalf("failed to generate comment fields from llm: %v", err)
	}

	if commentData.Content == "" || commentData.Username == "" {
		t.Fatalf("expected generated comment fields to be non-empty, got: %+v", commentData)
	}
}
