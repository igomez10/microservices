package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	defaultGeminiModel = "gemini-2.5-flash"
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

type TestDataGenerator interface {
	GenerateUser(ctx context.Context, in UserGenerationInput) (GeneratedUserData, error)
	GenerateComment(ctx context.Context, in CommentGenerationInput) (GeneratedCommentData, error)
}

type DefaultTestDataGenerator struct{}

func (d *DefaultTestDataGenerator) GenerateUser(_ context.Context, in UserGenerationInput) (GeneratedUserData, error) {
	return GeneratedUserData{
		Username:  in.FallbackUsername,
		FirstName: in.FallbackFirstName,
		LastName:  in.FallbackLastName,
		Email:     in.FallbackEmail,
	}, nil
}

func (d *DefaultTestDataGenerator) GenerateComment(_ context.Context, in CommentGenerationInput) (GeneratedCommentData, error) {
	return GeneratedCommentData{
		Content:  in.FallbackContent,
		Username: in.FallbackUsername,
	}, nil
}

type GeminiModelClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type GeminiTestDataGenerator struct {
	client GeminiModelClient
}

func NewGeminiTestDataGenerator(client GeminiModelClient) *GeminiTestDataGenerator {
	return &GeminiTestDataGenerator{client: client}
}

func (g *GeminiTestDataGenerator) GenerateUser(ctx context.Context, in UserGenerationInput) (GeneratedUserData, error) {
	prompt := fmt.Sprintf(`Return ONLY valid compact JSON with keys "username","first_name","last_name","email".
Generate realistic testing values. Keep username URL-safe and <= 64 chars.
Fallback values:
username=%q
first_name=%q
last_name=%q
email=%q`, in.FallbackUsername, in.FallbackFirstName, in.FallbackLastName, in.FallbackEmail)

	resp, err := g.client.Generate(ctx, prompt)
	if err != nil {
		return GeneratedUserData{}, err
	}

	var candidate struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(resp)), &candidate); err != nil {
		return GeneratedUserData{}, fmt.Errorf("invalid gemini user payload: %w", err)
	}

	user := GeneratedUserData{
		Username:  firstNonEmpty(sanitizeIdentifier(candidate.Username), in.FallbackUsername),
		FirstName: firstNonEmpty(strings.TrimSpace(candidate.FirstName), in.FallbackFirstName),
		LastName:  firstNonEmpty(strings.TrimSpace(candidate.LastName), in.FallbackLastName),
		Email:     firstNonEmpty(strings.TrimSpace(candidate.Email), in.FallbackEmail),
	}
	return user, nil
}

func (g *GeminiTestDataGenerator) GenerateComment(ctx context.Context, in CommentGenerationInput) (GeneratedCommentData, error) {
	prompt := fmt.Sprintf(`Return ONLY valid compact JSON with keys "content","username".
Generate one short social comment for test data.
Fallback values:
content=%q
username=%q`, in.FallbackContent, in.FallbackUsername)

	resp, err := g.client.Generate(ctx, prompt)
	if err != nil {
		return GeneratedCommentData{}, err
	}

	var candidate struct {
		Content  string `json:"content"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(resp)), &candidate); err != nil {
		return GeneratedCommentData{}, fmt.Errorf("invalid gemini comment payload: %w", err)
	}

	comment := GeneratedCommentData{
		Content:  firstNonEmpty(strings.TrimSpace(candidate.Content), in.FallbackContent),
		Username: firstNonEmpty(sanitizeIdentifier(candidate.Username), in.FallbackUsername),
	}
	return comment, nil
}

type FallbackTestDataGenerator struct {
	primary  TestDataGenerator
	fallback TestDataGenerator
}

func NewFallbackTestDataGenerator(primary, fallback TestDataGenerator) *FallbackTestDataGenerator {
	return &FallbackTestDataGenerator{primary: primary, fallback: fallback}
}

func (g *FallbackTestDataGenerator) GenerateUser(ctx context.Context, in UserGenerationInput) (GeneratedUserData, error) {
	data, err := g.primary.GenerateUser(ctx, in)
	if err == nil {
		return data, nil
	}
	slogError("gemini user generation failed, using fallback", err)
	return g.fallback.GenerateUser(ctx, in)
}

func (g *FallbackTestDataGenerator) GenerateComment(ctx context.Context, in CommentGenerationInput) (GeneratedCommentData, error) {
	data, err := g.primary.GenerateComment(ctx, in)
	if err == nil {
		return data, nil
	}
	slogError("gemini comment generation failed, using fallback", err)
	return g.fallback.GenerateComment(ctx, in)
}

func NewTestDataGeneratorFromEnv() TestDataGenerator {
	apiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	model := strings.TrimSpace(os.Getenv("GEMINI_MODEL"))
	if apiKey == "" {
		return &DefaultTestDataGenerator{}
	}
	if model == "" {
		model = defaultGeminiModel
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	geminiClient, err := NewGeminiClient(ctx, GeminiClientConfig{
		APIKey: apiKey,
		Model:  model,
	})
	if err != nil {
		slogError("gemini client init failed, using default generator", err)
		return &DefaultTestDataGenerator{}
	}

	return NewFallbackTestDataGenerator(
		NewGeminiTestDataGenerator(geminiClient),
		&DefaultTestDataGenerator{},
	)
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func sanitizeIdentifier(value string) string {
	sanitized := validIdentifierChars.ReplaceAllString(strings.TrimSpace(value), "-")
	if sanitized == "" {
		return sanitized
	}
	if len(sanitized) > 64 {
		return sanitized[:64]
	}
	return sanitized
}

func slogError(msg string, err error) {
	slog.Warn(msg, "error", err)
}

func extractJSONObject(text string) string {
	trimmed := strings.TrimSpace(text)
	if strings.HasPrefix(trimmed, "```") {
		lines := strings.Split(trimmed, "\n")
		if len(lines) >= 3 {
			lines = lines[1:]
			if strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
				lines = lines[:len(lines)-1]
			}
			trimmed = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end > start {
		return trimmed[start : end+1]
	}
	return trimmed
}
