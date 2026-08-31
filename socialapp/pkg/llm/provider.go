package llm

import "context"

const DefaultModel = DefaultGeminiModel

type Provider interface {
	Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
}

// Models enum
type Model string

const (
	ModelGemini25FlashLite Model = "gemini-2.5-flash-lite"
)

type MimeType string

const (
	MimeTypeJSON MimeType = "application/json"
)

type GenerateRequest struct {
	Model       Model
	Prompt      string
	Temperature float32

	// Structured output controls. When both are provided, Gemini is asked
	// to return output constrained by the provided schema.
	ResponseMIMEType   MimeType
	ResponseJSONSchema map[string]any
}

type GenerateResponse struct {
	Text string
}
