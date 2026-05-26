package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type GeminiClientConfig struct {
	APIKey string
	Model  string
}

type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, cfg GeminiClientConfig) (*GeminiClient, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("missing GOOGLE_CLOUD_API_KEY")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		cfg.Model = defaultGeminiModel
	}
	c, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendVertexAI,
	})
	if err != nil {
		return nil, err
	}
	return &GeminiClient{
		client: c,
		model:  cfg.Model,
	}, nil
}

func (c *GeminiClient) Generate(ctx context.Context, prompt string) (string, error) {
	temperature := float32(0.2)
	resp, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature: &temperature,
			ResponseJsonSchema: ,
		},
	)
	if err != nil {
		return "", err
	}

	// Decode from the canonical response structure to avoid relying on SDK helper methods.
	raw, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	var wire struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(raw, &wire); err != nil {
		return "", err
	}
	if len(wire.Candidates) == 0 || len(wire.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty response")
	}

	text := strings.TrimSpace(wire.Candidates[0].Content.Parts[0].Text)
	if text == "" {
		return "", fmt.Errorf("gemini returned blank text")
	}
	return text, nil
}
