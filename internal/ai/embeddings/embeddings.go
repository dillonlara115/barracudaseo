package embeddings

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

const (
	EmbeddingModel = "text-embedding-004"
	EmbeddingDim   = 768
)

// Service wraps the Gemini embedding API
type Service struct {
	client *genai.Client
	logger *zap.Logger
}

// NewService creates an embedding service from a Gemini API key
func NewService(apiKey string, logger *zap.Logger) (*Service, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client for embeddings: %w", err)
	}

	return &Service{client: client, logger: logger}, nil
}

// Embed generates a 768-dimension float32 vector for the given text
func (s *Service) Embed(ctx context.Context, text string) ([]float32, error) {
	contents := []*genai.Content{genai.NewContentFromText(text, genai.RoleUser)}

	result, err := s.client.Models.EmbedContent(ctx, EmbeddingModel, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}

	if result == nil || len(result.Embeddings) == 0 || len(result.Embeddings[0].Values) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}

	return result.Embeddings[0].Values, nil
}

// EmbedBatch generates embeddings for multiple texts in a single API call
func (s *Service) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	contents := make([]*genai.Content, len(texts))
	for i, text := range texts {
		contents[i] = genai.NewContentFromText(text, genai.RoleUser)
	}

	result, err := s.client.Models.EmbedContent(ctx, EmbeddingModel, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("batch embedding request failed: %w", err)
	}

	if result == nil || len(result.Embeddings) != len(texts) {
		return nil, fmt.Errorf("unexpected embedding response: got %d embeddings for %d texts",
			len(result.Embeddings), len(texts))
	}

	vectors := make([][]float32, len(texts))
	for i, emb := range result.Embeddings {
		vectors[i] = emb.Values
	}

	return vectors, nil
}

// FormatForPgvector converts a float32 slice to a pgvector-compatible string: [0.1,0.2,...]
func FormatForPgvector(v []float32) string {
	if len(v) == 0 {
		return "[]"
	}
	s := "["
	for i, val := range v {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%g", val)
	}
	s += "]"
	return s
}
