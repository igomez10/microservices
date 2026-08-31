package llm

import (
	"context"
	"errors"
	"sync"
)

// MockProvider is a deterministic in-memory Provider implementation for unit tests.
type MockProvider struct {
	mu       sync.Mutex
	queue    []mockResult
	requests []GenerateRequest
}

type mockResult struct {
	resp GenerateResponse
	err  error
}

// NewMockProvider creates a new empty mock provider.
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// EnqueueResponse schedules a successful response for the next Generate call.
func (m *MockProvider) EnqueueResponse(resp GenerateResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queue = append(m.queue, mockResult{resp: resp})
}

// EnqueueError schedules an error for the next Generate call.
func (m *MockProvider) EnqueueError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queue = append(m.queue, mockResult{err: err})
}

// Requests returns a copy of all Generate requests received so far.
func (m *MockProvider) Requests() []GenerateRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]GenerateRequest, len(m.requests))
	copy(out, m.requests)
	return out
}

// Generate implements Provider. It pops from the queued results in FIFO order.
func (m *MockProvider) Generate(_ context.Context, req GenerateRequest) (GenerateResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requests = append(m.requests, req)
	if len(m.queue) == 0 {
		return GenerateResponse{}, errors.New("mock provider has no queued result")
	}

	next := m.queue[0]
	m.queue = m.queue[1:]
	if next.err != nil {
		return GenerateResponse{}, next.err
	}
	return next.resp, nil
}
