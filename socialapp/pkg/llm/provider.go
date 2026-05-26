package llm

import "context"

type Provider interface {
	Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
}

type GenerateRequest struct {
	Model       string
	Prompt      string
	Temperature *float32

	// Structured output controls. When both are provided, Gemini is asked
	// to return output constrained by the provided schema.
	ResponseMIMEType   string
	ResponseJSONSchema map[string]any
}

type GenerateResponse struct {
	Text string
}
