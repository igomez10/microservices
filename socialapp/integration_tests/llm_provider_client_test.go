package main

import (
	"context"
	"testing"

	"github.com/igomez10/microservices/socialapp/pkg/llm"
)

func TestLLMProviderClientGenerate(t *testing.T) {
	mock := llm.NewMockProvider()
	mock.EnqueueResponse(llm.GenerateResponse{Text: `{"username":"alice"}`})

	client := NewLLMProviderClient(mock, "gemini-test-model")
	text, err := client.Generate(context.Background(), "return json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != `{"username":"alice"}` {
		t.Fatalf("unexpected text: %q", text)
	}

	reqs := mock.Requests()
	if len(reqs) != 1 {
		t.Fatalf("expected one request, got %d", len(reqs))
	}
	if reqs[0].Model != "gemini-test-model" {
		t.Fatalf("unexpected model: %q", reqs[0].Model)
	}
	if reqs[0].ResponseMIMEType != "application/json" {
		t.Fatalf("expected response mime type application/json, got %q", reqs[0].ResponseMIMEType)
	}
}
