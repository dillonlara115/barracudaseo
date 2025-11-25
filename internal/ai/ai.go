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
	
	// Add word count if available
	if wordCount := getValue(page, "word_count"); wordCount != nil {
		promptBuilder.WriteString(fmt.Sprintf("Word Count: %v\n", wordCount))
	}

	// Extract headings and content structure from data field
	if data, ok := page["data"].(map[string]interface{}); ok {
		if headings, ok := data["headings"].([]interface{}); ok && len(headings) > 0 {
			promptBuilder.WriteString("\nPage Structure:\n")
			// Include first few headings to understand page content
			maxHeadings := 10
			if len(headings) < maxHeadings {
				maxHeadings = len(headings)
			}
			for i := 0; i < maxHeadings; i++ {
				if heading, ok := headings[i].(map[string]interface{}); ok {
					level := getValue(heading, "level")
					text := getString(heading, "text")
					if text != "" {
						promptBuilder.WriteString(fmt.Sprintf("  H%v: %s\n", level, text))
					}
				}
			}
		}
		
		// Count links if available
		if internalLinks, ok := data["internal_links"].([]interface{}); ok {
			promptBuilder.WriteString(fmt.Sprintf("\nInternal Links: %d\n", len(internalLinks)))
		}
		if externalLinks, ok := data["external_links"].([]interface{}); ok {
			promptBuilder.WriteString(fmt.Sprintf("External Links: %d\n", len(externalLinks)))
		}
	}

	if gscData != nil && len(gscData) > 0 {
		promptBuilder.WriteString("\nGoogle Search Console Performance Data:\n")
		if impressions, ok := gscData["impressions"].(float64); ok && impressions > 0 {
			promptBuilder.WriteString(fmt.Sprintf("- Monthly Impressions: %.0f\n", impressions))
		}
		if clicks, ok := gscData["clicks"].(float64); ok && clicks > 0 {
			promptBuilder.WriteString(fmt.Sprintf("- Monthly Clicks: %.0f\n", clicks))
		}
		if ctr, ok := gscData["ctr"].(float64); ok && ctr > 0 {
			ctrPercent := ctr * 100
			promptBuilder.WriteString(fmt.Sprintf("- Click-Through Rate (CTR): %.2f%%\n", ctrPercent))
		}
		if position, ok := gscData["position"].(float64); ok && position > 0 {
			promptBuilder.WriteString(fmt.Sprintf("- Average Position: %.1f\n", position))
		}
		if topQueries, ok := gscData["top_queries"].([]string); ok && len(topQueries) > 0 {
			promptBuilder.WriteString("- Top Search Queries: ")
			for i, query := range topQueries {
				if i > 0 {
					promptBuilder.WriteString(", ")
				}
				promptBuilder.WriteString(query)
			}
			promptBuilder.WriteString("\n")
		}
		promptBuilder.WriteString("\n**IMPORTANT: Use this traffic data to inform your recommendation. If the page has high impressions but low CTR, mention optimization opportunities. If it ranks well (#1-10), emphasize the importance of fixing issues to maintain rankings. If it has high traffic, prioritize this fix.**\n")
	}

	// Determine issue type and generate focused solution
	issueType := getString(issue, "type")
	
	// For content-generation issues, format response with recommendation first, then insight
	switch issueType {
	case "missing_meta_description":
		promptBuilder.WriteString("\n**CRITICAL: You MUST format your response EXACTLY as shown below. Do not include any other sections or explanations.**\n\n")
		if gscData != nil && len(gscData) > 0 {
			if impressions, ok := gscData["impressions"].(float64); ok && impressions > 0 {
				promptBuilder.WriteString(fmt.Sprintf("**IMPORTANT CONTEXT: This page receives %.0f monthly impressions in Google Search. Adding a meta description could significantly improve click-through rate.**\n\n", impressions))
			}
			if ctr, ok := gscData["ctr"].(float64); ok && ctr > 0 && ctr < 0.02 {
				promptBuilder.WriteString(fmt.Sprintf("**OPTIMIZATION OPPORTUNITY: Current CTR is %.2f%%, which is below average. A compelling meta description could improve this.**\n\n", ctr*100))
			}
		}
		promptBuilder.WriteString("Based on the page URL, title, H1, and content structure, create a compelling meta description that:\n")
		promptBuilder.WriteString("- Is 150-160 characters long\n")
		promptBuilder.WriteString("- Accurately describes the page content\n")
		promptBuilder.WriteString("- Includes a call-to-action when appropriate\n")
		promptBuilder.WriteString("- Is optimized for search results\n")
		if gscData != nil {
			if topQueries, ok := gscData["top_queries"].([]string); ok && len(topQueries) > 0 {
				promptBuilder.WriteString("- Incorporates relevant keywords from top search queries naturally\n")
			}
		}
		promptBuilder.WriteString("\n**Your response MUST start with RECOMMENDATION: followed by the meta description text (no HTML tags, no quotes, just the text). Then on a new line, write INSIGHT: followed by a brief 1-2 sentence explanation that references traffic data if available.**\n\n")
		promptBuilder.WriteString("Example format:\n")
		promptBuilder.WriteString("RECOMMENDATION: Discover the latest updates and features from TransferForge. Stay informed and enhance your experience with our ongoing improvements!\n\n")
		promptBuilder.WriteString("INSIGHT: This meta description effectively summarizes the page content while including a clear call-to-action that encourages user engagement.\n")
	case "missing_title":
		promptBuilder.WriteString("\n**CRITICAL: You MUST format your response EXACTLY as shown below. Do not include any other sections or explanations.**\n\n")
		if gscData != nil && len(gscData) > 0 {
			if impressions, ok := gscData["impressions"].(float64); ok && impressions > 0 {
				promptBuilder.WriteString(fmt.Sprintf("**IMPORTANT CONTEXT: This page receives %.0f monthly impressions in Google Search. Adding a title tag is critical for search visibility.**\n\n", impressions))
			}
		}
		promptBuilder.WriteString("Based on the page URL, H1, and content structure, create a descriptive title that:\n")
		promptBuilder.WriteString("- Is 50-60 characters long\n")
		promptBuilder.WriteString("- Is keyword-rich and descriptive\n")
		promptBuilder.WriteString("- Accurately represents the page content\n")
		if gscData != nil {
			if topQueries, ok := gscData["top_queries"].([]string); ok && len(topQueries) > 0 {
				promptBuilder.WriteString("- Incorporates relevant keywords from top search queries naturally\n")
			}
		}
		promptBuilder.WriteString("\n**Your response MUST start with RECOMMENDATION: followed by the title text (no HTML tags, no quotes, just the text). Then on a new line, write INSIGHT: followed by a brief 1-2 sentence explanation that references traffic data if available.**\n\n")
		promptBuilder.WriteString("Example format:\n")
		promptBuilder.WriteString("RECOMMENDATION: Beta Launch - TransferForge: Features & Updates\n\n")
		promptBuilder.WriteString("INSIGHT: This title effectively describes the page content while staying within optimal length limits for search engine display.\n")
	case "missing_h1", "empty_h1":
		promptBuilder.WriteString("\n**CRITICAL: You MUST format your response EXACTLY as shown below. Do not include any other sections or explanations.**\n\n")
		promptBuilder.WriteString("Based on the page URL, title, and content structure, create an H1 that:\n")
		promptBuilder.WriteString("- Accurately describes the page's primary topic\n")
		promptBuilder.WriteString("- Is descriptive and keyword-rich\n")
		promptBuilder.WriteString("- Matches the page content\n\n")
		promptBuilder.WriteString("**Your response MUST start with RECOMMENDATION: followed by the H1 text (no HTML tags, no quotes, just the text). Then on a new line, write INSIGHT: followed by a brief 1-2 sentence explanation.**\n\n")
		promptBuilder.WriteString("Example format:\n")
		promptBuilder.WriteString("RECOMMENDATION: Beta Launch - TransferForge Features & Updates\n\n")
		promptBuilder.WriteString("INSIGHT: This H1 clearly communicates the page's main topic and aligns with SEO best practices for heading structure.\n")
	case "long_title", "short_title":
		currentTitle := getString(page, "title")
		promptBuilder.WriteString("\n**CRITICAL: You MUST format your response EXACTLY as shown below. Do not include any other sections or explanations.**\n\n")
		promptBuilder.WriteString(fmt.Sprintf("Current title: \"%s\"\n", currentTitle))
		if gscData != nil && len(gscData) > 0 {
			if impressions, ok := gscData["impressions"].(float64); ok && impressions > 0 {
				promptBuilder.WriteString(fmt.Sprintf("**IMPORTANT CONTEXT: This page receives %.0f monthly impressions. Optimizing the title could improve click-through rate.**\n\n", impressions))
			}
			if ctr, ok := gscData["ctr"].(float64); ok && ctr > 0 && ctr < 0.02 {
				promptBuilder.WriteString(fmt.Sprintf("**OPTIMIZATION OPPORTUNITY: Current CTR is %.2f%%. An optimized title could improve this.**\n\n", ctr*100))
			}
		}
		promptBuilder.WriteString("Based on the current title, page URL, H1, and content, create an optimized title that:\n")
		promptBuilder.WriteString("- Is exactly 50-60 characters long\n")
		promptBuilder.WriteString("- Improves upon the current title\n")
		promptBuilder.WriteString("- Is optimized for SEO and user experience\n")
		if gscData != nil {
			if topQueries, ok := gscData["top_queries"].([]string); ok && len(topQueries) > 0 {
				promptBuilder.WriteString("- Incorporates relevant keywords from top search queries naturally\n")
			}
		}
		promptBuilder.WriteString("\n**Your response MUST start with RECOMMENDATION: followed by the improved title text (no HTML tags, no quotes, just the text). Then on a new line, write INSIGHT: followed by a brief 1-2 sentence explanation that references traffic data if available.**\n\n")
		promptBuilder.WriteString("Example format:\n")
		promptBuilder.WriteString("RECOMMENDATION: Exciting Beta Launch - TransferForge: Features & Updates Revealed\n\n")
		promptBuilder.WriteString("INSIGHT: This improved title expands the original to 55 characters, adding descriptive keywords that better convey the page content and improve search visibility.\n")
	case "long_meta_description", "short_meta_description":
		currentMetaDesc := getString(page, "meta_description")
		promptBuilder.WriteString("\n**CRITICAL: You MUST format your response EXACTLY as shown below. Do not include any other sections or explanations.**\n\n")
		promptBuilder.WriteString(fmt.Sprintf("Current meta description: \"%s\"\n", currentMetaDesc))
		if gscData != nil && len(gscData) > 0 {
			if impressions, ok := gscData["impressions"].(float64); ok && impressions > 0 {
				promptBuilder.WriteString(fmt.Sprintf("**IMPORTANT CONTEXT: This page receives %.0f monthly impressions. Optimizing the meta description could improve CTR.**\n\n", impressions))
			}
			if ctr, ok := gscData["ctr"].(float64); ok && ctr > 0 && ctr < 0.02 {
				promptBuilder.WriteString(fmt.Sprintf("**OPTIMIZATION OPPORTUNITY: Current CTR is %.2f%%. An optimized meta description could improve this.**\n\n", ctr*100))
			}
		}
		promptBuilder.WriteString("Based on the current meta description, page URL, title, H1, and content, create an optimized meta description that:\n")
		promptBuilder.WriteString("- Is exactly 150-160 characters long\n")
		promptBuilder.WriteString("- Improves upon the current meta description\n")
		promptBuilder.WriteString("- Is compelling and includes a call-to-action when appropriate\n")
		if gscData != nil {
			if topQueries, ok := gscData["top_queries"].([]string); ok && len(topQueries) > 0 {
				promptBuilder.WriteString("- Incorporates relevant keywords from top search queries naturally\n")
			}
		}
		promptBuilder.WriteString("\n**Your response MUST start with RECOMMENDATION: followed by the improved meta description text (no HTML tags, no quotes, just the text). Then on a new line, write INSIGHT: followed by a brief 1-2 sentence explanation that references traffic data if available.**\n\n")
		promptBuilder.WriteString("Example format:\n")
		promptBuilder.WriteString("RECOMMENDATION: Discover the latest updates and features from TransferForge. Stay informed and enhance your experience!\n\n")
		promptBuilder.WriteString("INSIGHT: This improved meta description optimizes length while maintaining clarity and adding a compelling call-to-action.\n")
	case "missing_image_alt":
		promptBuilder.WriteString("\n**CRITICAL: You MUST format your response EXACTLY as shown below. Do not include any other sections or explanations.**\n\n")
		promptBuilder.WriteString("Based on the page content and context, provide descriptive alt text suggestions.\n")
		promptBuilder.WriteString("If specific image context is not available, provide general guidelines for writing good alt text.\n\n")
		promptBuilder.WriteString("**Your response MUST start with RECOMMENDATION: followed by the alt text or guidelines. Then on a new line, write INSIGHT: followed by a brief explanation if needed.**\n\n")
		promptBuilder.WriteString("Example format:\n")
		promptBuilder.WriteString("RECOMMENDATION: Descriptive alt text guidelines: Use concise, specific descriptions that convey the image's purpose and content.\n\n")
		promptBuilder.WriteString("INSIGHT: Good alt text improves accessibility and helps search engines understand image content.\n")
	default:
		// For other issue types, provide a brief explanation with solution
		promptBuilder.WriteString("\n**TASK: Provide a concise solution for this issue.**\n")
		promptBuilder.WriteString("Focus on:\n")
		promptBuilder.WriteString("1. Brief explanation of why it matters (1-2 sentences)\n")
		promptBuilder.WriteString("2. Specific steps to fix it\n")
		promptBuilder.WriteString("3. Example fix or code snippet if applicable\n\n")
		promptBuilder.WriteString("**Be concise and actionable.**\n")
	}

	// Add system message for content-generation issues to enforce format
	var messages []Message
	// issueType is already declared above, so we reuse it here
	
	// For content-generation issues, add a system message to enforce strict formatting
	if issueType == "missing_meta_description" || issueType == "missing_title" || 
	   issueType == "missing_h1" || issueType == "empty_h1" ||
	   issueType == "long_title" || issueType == "short_title" ||
	   issueType == "long_meta_description" || issueType == "short_meta_description" ||
	   issueType == "missing_image_alt" {
		messages = []Message{
			{
				Role: "system",
				Content: "You are an SEO expert. When asked to generate recommendations, you MUST format your response EXACTLY as: RECOMMENDATION: [text]\n\nINSIGHT: [brief explanation]. Do NOT include numbered sections, headers, or any other formatting. Only output the RECOMMENDATION and INSIGHT sections.",
			},
			{Role: "user", Content: promptBuilder.String()},
		}
	} else {
		messages = []Message{
			{Role: "user", Content: promptBuilder.String()},
		}
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

	if gscSummary, ok := crawlData["gsc_summary"].(map[string]interface{}); ok && len(gscSummary) > 0 {
		promptBuilder.WriteString("\nGoogle Search Console Performance Summary:\n")
		if impressions, ok := gscSummary["total_impressions"].(float64); ok && impressions > 0 {
			promptBuilder.WriteString(fmt.Sprintf("- Total Monthly Impressions: %.0f\n", impressions))
		}
		if clicks, ok := gscSummary["total_clicks"].(float64); ok && clicks > 0 {
			promptBuilder.WriteString(fmt.Sprintf("- Total Monthly Clicks: %.0f\n", clicks))
		}
		if ctr, ok := gscSummary["average_ctr"].(float64); ok && ctr > 0 {
			ctrPercent := ctr * 100
			promptBuilder.WriteString(fmt.Sprintf("- Average CTR: %.2f%%\n", ctrPercent))
		}
		if position, ok := gscSummary["average_position"].(float64); ok && position > 0 {
			promptBuilder.WriteString(fmt.Sprintf("- Average Position: %.1f\n", position))
		}
		if capturedOn, ok := gscSummary["captured_on"].(string); ok {
			promptBuilder.WriteString(fmt.Sprintf("- Data Period: %s\n", capturedOn))
		}
		promptBuilder.WriteString("\n**Use this search performance data to prioritize recommendations. Pages with high impressions but low CTR are prime candidates for optimization. Issues affecting high-traffic pages should be prioritized.**\n")
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

