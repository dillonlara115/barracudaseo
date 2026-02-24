package aicontext

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dillonlara115/barracudaseo/internal/ai/providers"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// ContextType identifies the kind of AI request being assembled
type ContextType string

const (
	ContextBrief      ContextType = "brief"
	ContextArticle    ContextType = "article"
	ContextDiagnostic ContextType = "diagnostic"
	ContextVoice      ContextType = "voice"
	ContextDigest     ContextType = "digest"
	ContextExplain    ContextType = "explain"
)

// Params holds optional parameters for context assembly
type Params struct {
	Keyword     string
	PageURL     string
	DateRange   string // e.g. "30d", "60d", "90d"
	TopN        int    // number of similar pages to include
	ExtraPrompt string // appended to the system prompt
}

// Builder assembles AI prompts from project data
type Builder struct {
	serviceRole *supabase.Client
	logger      *zap.Logger
}

// NewBuilder creates a new context builder
func NewBuilder(serviceRoleClient *supabase.Client, logger *zap.Logger) *Builder {
	return &Builder{
		serviceRole: serviceRoleClient,
		logger:      logger,
	}
}

// CrawledPage represents a crawled page record
type CrawledPage struct {
	URL             string      `json:"url"`
	Title           string      `json:"title"`
	MetaDescription string      `json:"meta_description"`
	Headings        interface{} `json:"headings"`
	BodySummary     string      `json:"body_summary"`
	WordCount       int         `json:"word_count"`
}

// VoiceProfile represents the writing voice for a project
type VoiceProfile struct {
	Tone          string `json:"tone"`
	Structure     string `json:"structure"`
	SentenceStyle string `json:"sentence_style"`
	BrandContext  string `json:"brand_context"`
	AvoidList     string `json:"avoid_list"`
}

// GSCRow represents a row of GSC data
type GSCRow struct {
	Query       string  `json:"query"`
	Page        string  `json:"page"`
	Clicks      int     `json:"clicks"`
	Impressions int     `json:"impressions"`
	CTR         float64 `json:"ctr"`
	Position    float64 `json:"position"`
}

// MemoryEntry represents a project memory fact
type MemoryEntry struct {
	MemoryText    string `json:"memory_text"`
	SourceFeature string `json:"source_feature"`
}

// BuildMessages assembles the system prompt and user content for an AI request
func (b *Builder) BuildMessages(ctx context.Context, projectID string, ctxType ContextType, params Params) ([]providers.Message, error) {
	topN := params.TopN
	if topN == 0 {
		topN = 5
	}

	var prompt strings.Builder
	prompt.WriteString(systemPreamble(ctxType))

	// Load crawled pages (top N by word count as a proxy for relevance without embedding search)
	pages, err := b.loadCrawledPages(ctx, projectID, topN)
	if err != nil {
		b.logger.Warn("Failed to load crawled pages for context", zap.Error(err))
	} else if len(pages) > 0 {
		prompt.WriteString("\n## Crawled Site Content\n")
		for _, p := range pages {
			prompt.WriteString(fmt.Sprintf("### %s\nURL: %s\nMeta: %s\nWord count: %d\n",
				p.Title, p.URL, p.MetaDescription, p.WordCount))
			if p.BodySummary != "" {
				summary := p.BodySummary
				if len(summary) > 500 {
					summary = summary[:500] + "..."
				}
				prompt.WriteString(fmt.Sprintf("Summary: %s\n", summary))
			}
			prompt.WriteString("\n")
		}
	}

	// Load GSC data snapshot
	gscRows, err := b.loadGSCData(ctx, projectID, params.PageURL)
	if err != nil {
		b.logger.Warn("Failed to load GSC data for context", zap.Error(err))
	} else if len(gscRows) > 0 {
		prompt.WriteString("\n## GSC Performance Snapshot\n")
		for _, row := range gscRows {
			prompt.WriteString(fmt.Sprintf("- Query: \"%s\" | Page: %s | Clicks: %d | Impressions: %d | CTR: %.2f%% | Pos: %.1f\n",
				row.Query, row.Page, row.Clicks, row.Impressions, row.CTR*100, row.Position))
		}
		prompt.WriteString("\n")
	}

	// Load voice profile (for content generation contexts)
	if ctxType == ContextArticle || ctxType == ContextBrief {
		voice, err := b.loadVoiceProfile(ctx, projectID)
		if err != nil {
			b.logger.Warn("Failed to load voice profile for context", zap.Error(err))
		} else if voice != nil {
			prompt.WriteString("\n## Writing Voice Profile\n")
			if voice.Tone != "" {
				prompt.WriteString(fmt.Sprintf("**Tone:** %s\n", voice.Tone))
			}
			if voice.Structure != "" {
				prompt.WriteString(fmt.Sprintf("**Structure:** %s\n", voice.Structure))
			}
			if voice.SentenceStyle != "" {
				prompt.WriteString(fmt.Sprintf("**Sentence Style:** %s\n", voice.SentenceStyle))
			}
			if voice.BrandContext != "" {
				prompt.WriteString(fmt.Sprintf("**Brand Context:** %s\n", voice.BrandContext))
			}
			if voice.AvoidList != "" {
				prompt.WriteString(fmt.Sprintf("**Avoid:** %s\n", voice.AvoidList))
			}
			prompt.WriteString("\n")
		}
	}

	// Load project memory
	memories, err := b.loadMemory(ctx, projectID)
	if err != nil {
		b.logger.Warn("Failed to load project memory", zap.Error(err))
	} else if len(memories) > 0 {
		prompt.WriteString("\n## Project Memory\n")
		for _, m := range memories {
			prompt.WriteString(fmt.Sprintf("- [%s] %s\n", m.SourceFeature, m.MemoryText))
		}
		prompt.WriteString("\n")
	}

	if params.ExtraPrompt != "" {
		prompt.WriteString("\n" + params.ExtraPrompt + "\n")
	}

	messages := []providers.Message{
		{Role: "system", Content: prompt.String()},
	}

	return messages, nil
}

func systemPreamble(ctxType ContextType) string {
	switch ctxType {
	case ContextBrief:
		return "You are an SEO content strategist. Generate a structured, data-backed content brief. Use the site data and GSC metrics provided to ground every recommendation.\n"
	case ContextArticle:
		return "You are a professional content writer. Write a full, publish-ready article draft. Match the writing voice profile exactly. Include internal link suggestions inline.\n"
	case ContextDiagnostic:
		return "You are an SEO technical analyst. Diagnose the issue using the data provided. Rank likely causes by probability and provide a recommended next action.\n"
	case ContextVoice:
		return "You are a brand voice analyst. Analyze the provided site content and extract five structured voice components: Tone, Structure, Sentence Style, Brand Context, and Avoid List.\n"
	case ContextDigest:
		return "You are an SEO performance analyst. Summarize the week's biggest GSC movers, new keyword appearances, and provide one prioritized action for the coming week.\n"
	case ContextExplain:
		return "You are an SEO analyst. Explain concisely why this keyword opportunity is underperforming and what is likely suppressing performance.\n"
	default:
		return "You are an SEO expert assistant. Answer based on the project data provided.\n"
	}
}

func (b *Builder) loadCrawledPages(ctx context.Context, projectID string, limit int) ([]CrawledPage, error) {
	var pages []CrawledPage
	data, _, err := b.serviceRole.From("crawled_pages").
		Select("url,title,meta_description,headings,body_summary,word_count", "", false).
		Eq("project_id", projectID).
		Order("word_count", nil).
		Limit(limit, "").
		Execute()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &pages); err != nil {
		return nil, fmt.Errorf("failed to parse crawled pages: %w", err)
	}
	return pages, nil
}

func (b *Builder) loadGSCData(ctx context.Context, projectID string, pageURL string) ([]GSCRow, error) {
	rowType := "query"
	if pageURL != "" {
		rowType = "page"
	}

	data, _, err := b.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", rowType).
		Execute()
	if err != nil {
		return nil, err
	}

	var rawRows []map[string]interface{}
	if err := json.Unmarshal(data, &rawRows); err != nil {
		return nil, fmt.Errorf("failed to parse GSC rows: %w", err)
	}

	var rows []GSCRow
	for _, raw := range rawRows {
		metrics, _ := raw["metrics"].(map[string]interface{})
		if metrics == nil {
			continue
		}

		dimVal, _ := raw["dimension_value"].(string)
		rt, _ := raw["row_type"].(string)

		row := GSCRow{
			Clicks:      int(getFloat64(metrics["clicks"])),
			Impressions: int(getFloat64(metrics["impressions"])),
			CTR:         getFloat64(metrics["ctr"]),
			Position:    getFloat64(metrics["position"]),
		}

		if rt == "page" {
			row.Page = dimVal
		} else {
			row.Query = dimVal
		}

		if pageURL != "" && row.Page != pageURL {
			continue
		}

		rows = append(rows, row)
		if len(rows) >= 20 {
			break
		}
	}

	return rows, nil
}

func getFloat64(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	default:
		return 0
	}
}

func (b *Builder) loadVoiceProfile(ctx context.Context, projectID string) (*VoiceProfile, error) {
	var profiles []VoiceProfile
	data, _, err := b.serviceRole.From("writing_voice").
		Select("tone,structure,sentence_style,brand_context,avoid_list", "", false).
		Eq("project_id", projectID).
		Limit(1, "").
		Execute()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, fmt.Errorf("failed to parse voice profile: %w", err)
	}
	if len(profiles) == 0 {
		return nil, nil
	}
	return &profiles[0], nil
}

func (b *Builder) loadMemory(ctx context.Context, projectID string) ([]MemoryEntry, error) {
	var entries []MemoryEntry
	data, _, err := b.serviceRole.From("project_memory").
		Select("memory_text,source_feature", "", false).
		Eq("project_id", projectID).
		Order("created_at", nil).
		Limit(20, "").
		Execute()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse project memory: %w", err)
	}
	return entries, nil
}
