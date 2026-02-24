package providers

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

// GeminiProvider implements the AIProvider interface using Google Gemini
type GeminiProvider struct {
	client *genai.Client
	model  string
	logger *zap.Logger
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey string, model string, logger *zap.Logger) (*GeminiProvider, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiProvider{
		client: client,
		model:  model,
		logger: logger,
	}, nil
}

// Completion generates a chat completion using Gemini
func (p *GeminiProvider) Completion(ctx context.Context, messages []Message) (string, error) {
	var systemInstruction string
	var contents []*genai.Content

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			systemInstruction = msg.Content
		case "user":
			contents = append(contents, genai.NewContentFromText(msg.Content, genai.RoleUser))
		case "assistant":
			contents = append(contents, genai.NewContentFromText(msg.Content, genai.RoleModel))
		}
	}

	config := &genai.GenerateContentConfig{}
	if systemInstruction != "" {
		config.SystemInstruction = genai.NewContentFromText(systemInstruction, genai.RoleUser)
	}

	result, err := p.client.Models.GenerateContent(ctx, p.model, contents, config)
	if err != nil {
		return "", fmt.Errorf("Gemini completion failed: %w", err)
	}

	return extractGeminiText(result), nil
}

// StreamCompletion streams a chat completion using Gemini, sending chunks to the callback
func (p *GeminiProvider) StreamCompletion(ctx context.Context, messages []Message, onChunk func(text string) error) error {
	var systemInstruction string
	var contents []*genai.Content

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			systemInstruction = msg.Content
		case "user":
			contents = append(contents, genai.NewContentFromText(msg.Content, genai.RoleUser))
		case "assistant":
			contents = append(contents, genai.NewContentFromText(msg.Content, genai.RoleModel))
		}
	}

	config := &genai.GenerateContentConfig{}
	if systemInstruction != "" {
		config.SystemInstruction = genai.NewContentFromText(systemInstruction, genai.RoleUser)
	}

	for result, err := range p.client.Models.GenerateContentStream(ctx, p.model, contents, config) {
		if err != nil {
			return fmt.Errorf("Gemini stream error: %w", err)
		}
		text := extractGeminiText(result)
		if text != "" {
			if err := onChunk(text); err != nil {
				return err
			}
		}
	}

	return nil
}

// UsageFromResponse extracts token usage from the last Gemini response
func (p *GeminiProvider) UsageFromResponse(result *genai.GenerateContentResponse) (inputTokens, outputTokens int32) {
	if result != nil && result.UsageMetadata != nil {
		return result.UsageMetadata.PromptTokenCount, result.UsageMetadata.CandidatesTokenCount
	}
	return 0, 0
}

func extractGeminiText(result *genai.GenerateContentResponse) string {
	if result == nil || len(result.Candidates) == 0 {
		return ""
	}

	var parts []string
	for _, candidate := range result.Candidates {
		if candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				parts = append(parts, part.Text)
			}
		}
	}
	return strings.Join(parts, "")
}
