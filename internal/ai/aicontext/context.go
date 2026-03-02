package aicontext

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/ai/embeddings"
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
	serviceRole  *supabase.Client
	embedService *embeddings.Service
	logger       *zap.Logger
}

// NewBuilder creates a new context builder.
// embedService may be nil; similarity search will fall back to word-count ordering.
func NewBuilder(serviceRoleClient *supabase.Client, embedService *embeddings.Service, logger *zap.Logger) *Builder {
	return &Builder{
		serviceRole:  serviceRoleClient,
		embedService: embedService,
		logger:       logger,
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
	Similarity      float64     `json:"similarity,omitempty"`
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
	Query       string   `json:"query"`
	Page        string   `json:"page"`
	Clicks      int      `json:"clicks"`
	Impressions int      `json:"impressions"`
	CTR         float64  `json:"ctr"`
	Position    float64  `json:"position"`
	TopQueries  []string `json:"top_queries,omitempty"` // For page rows: queries driving this page
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

	// Build the search query for semantic similarity from available params
	searchQuery := params.Keyword
	if searchQuery == "" && params.PageURL != "" {
		searchQuery = params.PageURL
	}

	// Try pgvector similarity search when we have a query and embedding service
	var pages []CrawledPage
	var usedSimilarity bool
	if searchQuery != "" && b.embedService != nil {
		similar, err := b.loadSimilarPages(ctx, projectID, searchQuery, topN)
		if err != nil {
			b.logger.Warn("Similarity search failed, falling back to word-count ordering",
				zap.String("query", searchQuery), zap.Error(err))
		} else if len(similar) > 0 {
			pages = similar
			usedSimilarity = true
		}
	}

	// Fallback: top N pages by word count (longest content first)
	if len(pages) == 0 {
		fallback, err := b.loadCrawledPages(ctx, projectID, topN)
		if err != nil {
			b.logger.Warn("Failed to load crawled pages for context", zap.Error(err))
		} else {
			pages = fallback
		}
	}

	if len(pages) > 0 {
		if usedSimilarity {
			prompt.WriteString(fmt.Sprintf("\n## Crawled Site Content (top %d pages by relevance to \"%s\")\n", len(pages), searchQuery))
		} else {
			prompt.WriteString("\n## Crawled Site Content\n")
		}
		for _, p := range pages {
			header := fmt.Sprintf("### %s\nURL: %s\nMeta: %s\nWord count: %d\n",
				p.Title, p.URL, p.MetaDescription, p.WordCount)
			if usedSimilarity && p.Similarity > 0 {
				header = fmt.Sprintf("### %s\nURL: %s\nMeta: %s\nWord count: %d\nRelevance: %.0f%%\n",
					p.Title, p.URL, p.MetaDescription, p.WordCount, p.Similarity*100)
			}
			prompt.WriteString(header)
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
	gscRows, keywordFound, err := b.loadGSCData(ctx, projectID, params.PageURL, params.Keyword)
	if err != nil {
		b.logger.Warn("Failed to load GSC data for context", zap.Error(err))
	} else if len(gscRows) > 0 {
		prompt.WriteString("\n## GSC Performance Snapshot\n")
		for _, row := range gscRows {
			prompt.WriteString(fmt.Sprintf("- Query: \"%s\" | Page: %s | Clicks: %d | Impressions: %d | CTR: %.2f%% | Pos: %.1f\n",
				row.Query, row.Page, row.Clicks, row.Impressions, row.CTR*100, row.Position))
			if len(row.TopQueries) > 0 {
				prompt.WriteString(fmt.Sprintf("  Top queries for this page: %s\n", strings.Join(row.TopQueries, ", ")))
			}
		}
		prompt.WriteString("\n")
		if ctxType == ContextExplain && params.Keyword != "" && !keywordFound {
			prompt.WriteString(fmt.Sprintf("Note: The keyword \"%s\" was not found in the cached GSC snapshot. It may have very low impressions, be outside the sync date range, or not yet be synced. Suggest the user re-sync GSC or verify the keyword appears in Search Console for the selected period.\n", params.Keyword))
		}
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

// loadCrawledPages returns the top N pages by word count (descending) as a fallback
// when no keyword/page is available for semantic search.
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

// loadSimilarPages uses pgvector cosine similarity to find the most relevant
// crawled pages for a given search query. It embeds the query text, then calls
// the match_crawled_pages RPC function in Supabase.
func (b *Builder) loadSimilarPages(ctx context.Context, projectID string, query string, limit int) ([]CrawledPage, error) {
	if b.embedService == nil {
		return nil, fmt.Errorf("embedding service not available")
	}

	vec, err := b.embedService.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	vecStr := embeddings.FormatForPgvector(vec)

	rpcBody := map[string]interface{}{
		"query_embedding":  vecStr,
		"match_project_id": projectID,
		"match_count":      limit,
		"match_threshold":  0.0,
	}

	result := b.serviceRole.Rpc("match_crawled_pages", "", rpcBody)
	if result == "" {
		return nil, fmt.Errorf("empty response from match_crawled_pages RPC")
	}

	var pages []CrawledPage
	if err := json.Unmarshal([]byte(result), &pages); err != nil {
		return nil, fmt.Errorf("failed to parse similarity results: %w", err)
	}

	return pages, nil
}

func (b *Builder) loadGSCData(ctx context.Context, projectID string, pageURL string, keyword string) ([]GSCRow, bool, error) {
	rowType := "query"
	if pageURL != "" {
		rowType = "page"
	}

	snapshotID := b.loadLatestSnapshotID(ctx, projectID)
	if snapshotID == "" {
		return nil, false, nil
	}

	// For Explain flow: ensure the requested keyword is included in the snapshot.
	var keywordRow *GSCRow
	if keyword != "" && rowType == "query" {
		kwData, _, err := b.serviceRole.From("gsc_performance_rows").
			Select("*", "", false).
			Eq("snapshot_id", snapshotID).
			Eq("project_id", projectID).
			Eq("row_type", "query").
			Eq("dimension_value", keyword).
			Limit(1, "").
			Execute()
		if err == nil {
			var kwRaw []map[string]interface{}
			if json.Unmarshal(kwData, &kwRaw) == nil && len(kwRaw) > 0 {
				keywordRow = rawToGSCRow(kwRaw[0], "query")
			}
		}
	}

	data, _, err := b.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("snapshot_id", snapshotID).
		Eq("project_id", projectID).
		Eq("row_type", rowType).
		Execute()
	if err != nil {
		return nil, false, err
	}

	var rawRows []map[string]interface{}
	if err := json.Unmarshal(data, &rawRows); err != nil {
		return nil, false, fmt.Errorf("failed to parse GSC rows: %w", err)
	}

	seen := make(map[string]bool)
	var rows []GSCRow

	if keywordRow != nil {
		rows = append(rows, *keywordRow)
		seen[keywordRow.Query] = true
	}

	for _, raw := range rawRows {
		rt, _ := raw["row_type"].(string)
		row := rawToGSCRow(raw, rt)
		if row == nil {
			continue
		}

		if pageURL != "" && row.Page != pageURL {
			continue
		}

		key := row.Query
		if key == "" {
			key = row.Page
		}
		if seen[key] {
			continue
		}
		seen[key] = true

		rows = append(rows, *row)
		if len(rows) >= 20 {
			break
		}
	}

	return rows, keywordRow != nil, nil
}

// loadLatestSnapshotID returns the ID of the most recent GSC snapshot for the project.
func (b *Builder) loadLatestSnapshotID(ctx context.Context, projectID string) string {
	data, _, err := b.serviceRole.From("gsc_performance_snapshots").
		Select("id,captured_on", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		return ""
	}
	var snapshots []map[string]interface{}
	if json.Unmarshal(data, &snapshots) != nil || len(snapshots) == 0 {
		return ""
	}
	type snapWithTime struct {
		id   string
		time time.Time
	}
	var withTime []snapWithTime
	for _, s := range snapshots {
		id, _ := s["id"].(string)
		if id == "" {
			id = fmt.Sprintf("%v", s["id"])
		}
		captured, _ := s["captured_on"].(string)
		t, _ := time.Parse("2006-01-02", captured)
		withTime = append(withTime, snapWithTime{id: id, time: t})
	}
	sort.Slice(withTime, func(i, j int) bool {
		return withTime[i].time.After(withTime[j].time)
	})
	return withTime[0].id
}

func rawToGSCRow(raw map[string]interface{}, rt string) *GSCRow {
	metrics, _ := raw["metrics"].(map[string]interface{})
	if metrics == nil {
		return nil
	}
	dimVal, _ := raw["dimension_value"].(string)
	row := &GSCRow{
		Clicks:      int(getFloat64(metrics["clicks"])),
		Impressions: int(getFloat64(metrics["impressions"])),
		CTR:         getFloat64(metrics["ctr"]),
		Position:    getFloat64(metrics["position"]),
	}
	if rt == "page" {
		row.Page = dimVal
		// Parse top_queries for page rows (helps Diagnose flow with keyword cannibalization)
		if tq, ok := raw["top_queries"].([]interface{}); ok && len(tq) > 0 {
			maxQ := 5
			if len(tq) < maxQ {
				maxQ = len(tq)
			}
			for i := 0; i < maxQ; i++ {
				if qm, ok := tq[i].(map[string]interface{}); ok {
					if q, ok := qm["query"].(string); ok && q != "" {
						row.TopQueries = append(row.TopQueries, q)
					}
				}
			}
		}
	} else {
		row.Query = dimVal
	}
	return row
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
