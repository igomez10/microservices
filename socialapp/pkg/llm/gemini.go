package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const DefaultGeminiModel = "gemini-2.5-flash-lite"

type GeminiConfig struct {
	APIKey string
	Model  string
}

type GeminiProvider struct {
	client *genai.Client
	model  string
}

func NewGeminiProvider(ctx context.Context, cfg GeminiConfig) (*GeminiProvider, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("missing gemini api key")
	}

	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = DefaultGeminiModel
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendVertexAI,
	})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	return &GeminiProvider{
		client: client,
		model:  model,
	}, nil
}

func (p *GeminiProvider) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	prompt := strings.TrimSpace(req.Prompt)
	if prompt == "" {
		return GenerateResponse{}, errors.New("prompt is required")
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = p.model
	}

	cfg := &genai.GenerateContentConfig{}
	if req.Temperature != 0 {
		cfg.Temperature = &req.Temperature
	}

	cfg.ResponseMIMEType = string(req.ResponseMIMEType)
	if req.ResponseJSONSchema != nil {
		cfg.ResponseJsonSchema = req.ResponseJSONSchema
	}

	resp, err := p.client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(prompt),
		cfg,
	)
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("generate gemini content: %w", err)
	}

	text := strings.TrimSpace(resp.Text())
	if text == "" {
		return GenerateResponse{}, errors.New("gemini returned empty response")
	}

	return GenerateResponse{Text: text}, nil
}
