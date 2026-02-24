package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/ai/providers"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// Service handles post-interaction memory extraction and storage
type Service struct {
	provider    *providers.GeminiProvider
	serviceRole *supabase.Client
	logger      *zap.Logger
}

// NewService creates a memory extraction service
func NewService(provider *providers.GeminiProvider, serviceRoleClient *supabase.Client, logger *zap.Logger) *Service {
	return &Service{
		provider:    provider,
		serviceRole: serviceRoleClient,
		logger:      logger,
	}
}

type memoryRecord struct {
	ID            string `json:"id"`
	ProjectID     string `json:"project_id"`
	MemoryText    string `json:"memory_text"`
	SourceFeature string `json:"source_feature"`
	CreatedAt     string `json:"created_at"`
}

// ExtractAndStore runs a lightweight AI call to extract key facts from an interaction,
// then persists them to the project_memory table.
func (s *Service) ExtractAndStore(ctx context.Context, projectID string, feature string, interactionText string) error {
	if interactionText == "" {
		return nil
	}

	// Truncate long interactions to keep the extraction prompt cheap
	if len(interactionText) > 4000 {
		interactionText = interactionText[:4000]
	}

	messages := []providers.Message{
		{
			Role:    "system",
			Content: "You extract key facts from SEO tool interactions. Return 1-5 concise bullet points of project-specific facts worth remembering for future AI calls. Each fact should be a single sentence. If there are no meaningful facts to extract, return NONE.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Feature: %s\n\nInteraction:\n%s", feature, interactionText),
		},
	}

	response, err := s.provider.Completion(ctx, messages)
	if err != nil {
		return fmt.Errorf("memory extraction failed: %w", err)
	}

	if response == "" || response == "NONE" {
		return nil
	}

	record := memoryRecord{
		ID:            uuid.New().String(),
		ProjectID:     projectID,
		MemoryText:    response,
		SourceFeature: feature,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal memory record: %w", err)
	}

	_, _, err = s.serviceRole.From("project_memory").Insert(json.RawMessage(data), false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to insert memory: %w", err)
	}

	s.logger.Debug("Stored project memory",
		zap.String("project_id", projectID),
		zap.String("feature", feature),
	)

	return nil
}
