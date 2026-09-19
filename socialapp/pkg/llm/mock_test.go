package llm

import (
	"context"
	"errors"
	"testing"
)

func TestMockProviderReturnsQueuedResponse(t *testing.T) {
	mock := NewMockProvider()
	mock.EnqueueResponse(GenerateResponse{Text: "hello"})

	resp, err := mock.Generate(context.Background(), GenerateRequest{Prompt: "p1", Model: "m1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "hello" {
		t.Fatalf("expected response text hello, got %q", resp.Text)
	}

	reqs := mock.Requests()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 request, got %d", len(reqs))
	}
	if reqs[0].Prompt != "p1" || reqs[0].Model != "m1" {
		t.Fatalf("unexpected captured request: %+v", reqs[0])
	}
}

func TestMockProviderReturnsQueuedError(t *testing.T) {
	mock := NewMockProvider()
	mock.EnqueueError(errors.New("boom"))

	_, err := mock.Generate(context.Background(), GenerateRequest{Prompt: "p2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "boom" {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestMockProviderExhaustedQueue(t *testing.T) {
	mock := NewMockProvider()

	_, err := mock.Generate(context.Background(), GenerateRequest{Prompt: "p3"})
	if err == nil {
		t.Fatal("expected exhaustion error, got nil")
	}
}
