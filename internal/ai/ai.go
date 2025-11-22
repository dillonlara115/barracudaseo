package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dillonlara115/barracuda/internal/ai/providers"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// Message represents a chat message for AI completion
type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// AIClient handles AI operations with provider abstraction
type AIClient struct {
	supabase     *supabase.Client
	serviceRole  *supabase.Client
	logger       *zap.Logger
	defaultModel string
}

// NewAIClient creates a new AI client
func NewAIClient(supabaseClient, serviceRoleClient *supabase.Client, logger *zap.Logger) *AIClient {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini" // Default model
	}

	return &AIClient{
		supabase:     supabaseClient,
		serviceRole:  serviceRoleClient,
		logger:       logger,
		defaultModel: model,
	}
}

// GetAPIKey retrieves the OpenAI API key for a user, falling back to app-wide key
func (c *AIClient) GetAPIKey(ctx context.Context, userID string) (string, error) {
	// First, try to get user's own key
	var result struct {
		OpenAIAPIKey *string `json:"openai_api_key"`
	}

	_, err := c.supabase.From("user_ai_settings").
		Select("openai_api_key", "", false).
		Eq("user_id", userID).
		Single().
		ExecuteTo(&result)

	if err == nil && result.OpenAIAPIKey != nil && *result.OpenAIAPIKey != "" {
		c.logger.Debug("Using user-provided OpenAI API key", zap.String("user_id", userID))
		return *result.OpenAIAPIKey, nil
	}

	// Fallback to app-wide key from environment
	appKey := os.Getenv("OPENAI_API_KEY")
	if appKey == "" {
		return "", fmt.Errorf("no OpenAI API key found (neither user key nor OPENAI_API_KEY env var)")
	}

	c.logger.Debug("Using app-wide OpenAI API key", zap.String("user_id", userID))
	return appKey, nil
}

// CreateChatCompletion creates a chat completion using the appropriate provider
func (c *AIClient) CreateChatCompletion(ctx context.Context, userID string, model string, messages []Message) (string, error) {
	if model == "" {
		model = c.defaultModel
	}

	// Get API key
	apiKey, err := c.GetAPIKey(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get API key: %w", err)
	}

	// Create OpenAI provider
	provider := providers.NewOpenAIProvider(apiKey, c.logger)

	// Convert messages to provider format
	providerMessages := make([]providers.Message, len(messages))
	for i, msg := range messages {
		providerMessages[i] = providers.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Call provider
	response, err := provider.Completion(ctx, providerMessages)
	if err != nil {
		return "", fmt.Errorf("AI completion failed: %w", err)
	}

	return response, nil
}

// GenerateIssueInsight generates an AI insight for a specific issue
func (c *AIClient) GenerateIssueInsight(ctx context.Context, userID string, issue map[string]interface{}, page map[string]interface{}, gscData map[string]interface{}) (string, error) {
	// Build prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("You are an SEO technical analyst. The user is auditing a website.\n\n")
	promptBuilder.WriteString("Here is the issue:\n")
	promptBuilder.WriteString(fmt.Sprintf("Message: %s\n", getString(issue, "message")))
	promptBuilder.WriteString(fmt.Sprintf("Type: %s\n", getString(issue, "type")))
	promptBuilder.WriteString(fmt.Sprintf("Severity: %s\n", getString(issue, "severity")))
	
	if val := getString(issue, "value"); val != "" {
		promptBuilder.WriteString(fmt.Sprintf("Value: %s\n", val))
	}
	if rec := getString(issue, "recommendation"); rec != "" {
		promptBuilder.WriteString(fmt.Sprintf("Current Recommendation: %s\n", rec))
	}

	promptBuilder.WriteString("\nHere is metadata about the page:\n")
	promptBuilder.WriteString(fmt.Sprintf("URL: %s\n", getString(page, "url")))
	promptBuilder.WriteString(fmt.Sprintf("Title: %s\n", getString(page, "title")))
	promptBuilder.WriteString(fmt.Sprintf("Meta Description: %s\n", getString(page, "meta_description")))
	promptBuilder.WriteString(fmt.Sprintf("Status Code: %v\n", getValue(page, "status_code")))
	promptBuilder.WriteString(fmt.Sprintf("H1: %s\n", getString(page, "h1")))

	// Count links if available
	if data, ok := page["data"].(map[string]interface{}); ok {
		if internalLinks, ok := data["internal_links"].([]interface{}); ok {
			promptBuilder.WriteString(fmt.Sprintf("Internal Links: %d\n", len(internalLinks)))
		}
		if externalLinks, ok := data["external_links"].([]interface{}); ok {
			promptBuilder.WriteString(fmt.Sprintf("External Links: %d\n", len(externalLinks)))
		}
	}

	if gscData != nil && len(gscData) > 0 {
		promptBuilder.WriteString("\nGSC Metrics:\n")
		for k, v := range gscData {
			promptBuilder.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
	}

	promptBuilder.WriteString("\nProvide:\n")
	promptBuilder.WriteString("1. Why this issue matters\n")
	promptBuilder.WriteString("2. How to fix it\n")
	promptBuilder.WriteString("3. Impact on SEO (low/medium/high)\n")
	promptBuilder.WriteString("4. Example fix or improved snippet\n")
	promptBuilder.WriteString("5. Additional considerations\n")

	messages := []Message{
		{Role: "user", Content: promptBuilder.String()},
	}

	return c.CreateChatCompletion(ctx, userID, "", messages)
}

// GenerateCrawlSummary generates an AI summary for an entire crawl
func (c *AIClient) GenerateCrawlSummary(ctx context.Context, userID string, crawlData map[string]interface{}) (string, error) {
	// Build prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("Act as a senior SEO strategist.\n\n")
	promptBuilder.WriteString("Here is the crawl data:\n")
	promptBuilder.WriteString(fmt.Sprintf("- Total pages: %v\n", getValue(crawlData, "total_pages")))
	promptBuilder.WriteString(fmt.Sprintf("- Total issues: %v\n", getValue(crawlData, "total_issues")))

	// Issues breakdown
	if issuesByType, ok := crawlData["issues_by_type"].(map[string]interface{}); ok {
		promptBuilder.WriteString("- Issues by type:\n")
		for issueType, count := range issuesByType {
			promptBuilder.WriteString(fmt.Sprintf("  %s: %v\n", issueType, count))
		}
	}

	if issuesBySeverity, ok := crawlData["issues_by_severity"].(map[string]interface{}); ok {
		promptBuilder.WriteString("- Issues by severity:\n")
		for severity, count := range issuesBySeverity {
			promptBuilder.WriteString(fmt.Sprintf("  %s: %v\n", severity, count))
		}
	}

	if slowPages, ok := crawlData["slow_pages"].([]interface{}); ok && len(slowPages) > 0 {
		promptBuilder.WriteString(fmt.Sprintf("- Slow pages (>3s): %d\n", len(slowPages)))
	}

	if redirectChains, ok := crawlData["redirect_chains"].(int); ok && redirectChains > 0 {
		promptBuilder.WriteString(fmt.Sprintf("- Pages with redirect chains: %d\n", redirectChains))
	}

	if metadataIssues, ok := crawlData["metadata_issues"].(int); ok && metadataIssues > 0 {
		promptBuilder.WriteString(fmt.Sprintf("- Pages missing metadata: %d\n", metadataIssues))
	}

	if gscSummary, ok := crawlData["gsc_summary"].(string); ok && gscSummary != "" {
		promptBuilder.WriteString(fmt.Sprintf("\nGSC metrics summary: %s\n", gscSummary))
	}

	promptBuilder.WriteString("\nProvide:\n")
	promptBuilder.WriteString("1. High-level executive summary\n")
	promptBuilder.WriteString("2. Top 5 fixes to prioritize\n")
	promptBuilder.WriteString("3. Major technical blockers\n")
	promptBuilder.WriteString("4. Quick wins\n")
	promptBuilder.WriteString("5. Opportunities for content improvement\n")
	promptBuilder.WriteString("6. Potential ranking improvements if fixes are implemented\n")
	promptBuilder.WriteString("7. Simple action plan for the next 7 days\n")

	messages := []Message{
		{Role: "user", Content: promptBuilder.String()},
	}

	return c.CreateChatCompletion(ctx, userID, "", messages)
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func getValue(m map[string]interface{}, key string) interface{} {
	if val, ok := m[key]; ok {
		return val
	}
	return nil
}

