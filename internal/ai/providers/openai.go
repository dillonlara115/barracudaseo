package providers

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// OpenAIProvider implements AIProvider using OpenAI API
type OpenAIProvider struct {
	client *openai.Client
	logger *zap.Logger
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, logger *zap.Logger) *OpenAIProvider {
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{
		client: client,
		logger: logger,
	}
}

// Completion creates a chat completion
func (p *OpenAIProvider) Completion(ctx context.Context, messages []Message) (string, error) {
	// Convert messages to OpenAI format
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Create completion request
	req := openai.ChatCompletionRequest{
		Model:    "gpt-4o-mini", // Default model
		Messages: openaiMessages,
	}

	// Call API
	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		p.logger.Error("OpenAI API error", zap.Error(err))
		return "", fmt.Errorf("openai API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}



