package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/ai/aicontext"
	"github.com/dillonlara115/barracudaseo/internal/ai/embeddings"
	"github.com/dillonlara115/barracudaseo/internal/ai/memory"
	"github.com/dillonlara115/barracudaseo/internal/ai/providers"
	"github.com/dillonlara115/barracudaseo/internal/ai/stream"
	"github.com/dillonlara115/barracudaseo/internal/ai/usage"
	"github.com/dillonlara115/barracudaseo/internal/sitecrawl"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AISuiteServices holds shared dependencies for all AI suite handlers
type AISuiteServices struct {
	FlashProvider  *providers.GeminiProvider
	LiteProvider   *providers.GeminiProvider
	EmbedService   *embeddings.Service
	ContextBuilder *aicontext.Builder
	MemoryService  *memory.Service
	UsageTracker   *usage.Tracker
	StreamHandler  *stream.Handler
	SiteCrawler    *sitecrawl.Service
}

// InitAISuite initializes all AI suite services. Returns nil if GEMINI_API_KEY is not set.
func (s *Server) InitAISuite() *AISuiteServices {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		s.logger.Warn("GEMINI_API_KEY not set — AI suite features disabled")
		return nil
	}

	flashProvider, err := providers.NewGeminiProvider(apiKey, "gemini-2.5-flash", s.logger)
	if err != nil {
		s.logger.Error("Failed to init Gemini Flash provider", zap.Error(err))
		return nil
	}

	liteProvider, err := providers.NewGeminiProvider(apiKey, "gemini-2.5-flash-lite", s.logger)
	if err != nil {
		s.logger.Error("Failed to init Gemini Flash-Lite provider", zap.Error(err))
		return nil
	}

	embedService, err := embeddings.NewService(apiKey, s.logger)
	if err != nil {
		s.logger.Error("Failed to init embedding service", zap.Error(err))
		return nil
	}

	ctxBuilder := aicontext.NewBuilder(s.serviceRole, embedService, s.logger)
	usageTracker := usage.NewTracker(s.serviceRole, s.logger)
	memoryService := memory.NewService(liteProvider, s.serviceRole, s.logger)
	streamHandler := stream.NewHandler(s.logger)
	siteCrawler := sitecrawl.NewService(s.serviceRole, embedService, s.logger)

	return &AISuiteServices{
		FlashProvider:  flashProvider,
		LiteProvider:   liteProvider,
		EmbedService:   embedService,
		ContextBuilder: ctxBuilder,
		MemoryService:  memoryService,
		UsageTracker:   usageTracker,
		StreamHandler:  streamHandler,
		SiteCrawler:    siteCrawler,
	}
}

// ---- Writing Voice Handlers ----

func (s *Server) handleGenerateVoice(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req struct {
			ProjectID string `json:"project_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if ok, err := s.verifyProjectAccess(userID, req.ProjectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		messages, err := ai.ContextBuilder.BuildMessages(r.Context(), req.ProjectID, aicontext.ContextVoice, aicontext.Params{TopN: 10})
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to build context")
			return
		}

		messages = append(messages, providers.Message{
			Role:    "user",
			Content: "Analyze the site content provided and generate the five voice components. Format each as a clearly labeled section: TONE:, STRUCTURE:, SENTENCE_STYLE:, BRAND_CONTEXT:, AVOID_LIST:. Each section should be a detailed paragraph.",
		})

		fullText, err := ai.StreamHandler.StreamResponse(r.Context(), w, ai.FlashProvider, messages)
		if err != nil {
			s.logger.Error("Voice generation stream failed", zap.Error(err))
			return
		}

		// Parse and store the voice profile
		voice := parseVoiceComponents(fullText)
		voice["id"] = uuid.New().String()
		voice["project_id"] = req.ProjectID
		voice["generated_at"] = time.Now().UTC().Format(time.RFC3339)
		voice["last_edited_at"] = voice["generated_at"]

		data, _ := json.Marshal(voice)
		_, _, err = s.serviceRole.From("writing_voice").
			Upsert(json.RawMessage(data), "project_id", "", "").
			Execute()
		if err != nil {
			s.logger.Error("Failed to store voice profile", zap.Error(err))
		}

		// Extract memory
		go ai.MemoryService.ExtractAndStore(context.Background(), req.ProjectID, "voice", fullText)
	}
}

func (s *Server) handleGetVoice(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		s.respondError(w, http.StatusBadRequest, "project_id required")
		return
	}

	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	data, _, err := s.serviceRole.From("writing_voice").
		Select("*", "", false).
		Eq("project_id", projectID).
		Limit(1, "").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch voice profile (table may not exist yet)", zap.Error(err))
		s.respondJSON(w, http.StatusOK, nil)
		return
	}

	var profiles []map[string]interface{}
	json.Unmarshal(data, &profiles)
	if len(profiles) == 0 {
		s.respondJSON(w, http.StatusOK, nil)
		return
	}

	s.respondJSON(w, http.StatusOK, profiles[0])
}

func (s *Server) handleUpdateVoice(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	projectID, _ := req["project_id"].(string)
	if projectID == "" {
		s.respondError(w, http.StatusBadRequest, "project_id required")
		return
	}

	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	req["last_edited_at"] = time.Now().UTC().Format(time.RFC3339)
	data, _ := json.Marshal(req)

	_, _, err := s.serviceRole.From("writing_voice").
		Update(json.RawMessage(data), "", "").
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to update voice profile")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ---- Content Brief Handlers ----

func (s *Server) handleGenerateBrief(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req struct {
			ProjectID string `json:"project_id"`
			Keyword   string `json:"keyword"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if ok, err := s.verifyProjectAccess(userID, req.ProjectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		// Check usage limit
		status, err := ai.UsageTracker.CheckLimit(r.Context(), req.ProjectID, "brief")
		if err == nil && !status.Allowed {
			s.respondError(w, http.StatusTooManyRequests, status.Message)
			return
		}

		messages, err := ai.ContextBuilder.BuildMessages(r.Context(), req.ProjectID, aicontext.ContextBrief, aicontext.Params{
			Keyword: req.Keyword,
			TopN:    5,
		})
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to build context")
			return
		}

		messages = append(messages, providers.Message{
			Role: "user",
			Content: fmt.Sprintf(`Generate a comprehensive content brief for the keyword: "%s"

Format the brief as a well-structured markdown document with these sections:

## Target Keyword
The primary keyword and search intent (informational / navigational / commercial / transactional).

## Secondary Keywords
A bullet list of 3-5 related keywords to weave into the content.

## Recommended Title
A compelling, SEO-optimized title tag (under 60 characters).

## Meta Description
A click-worthy meta description (under 160 characters).

## Content Angle
One sentence describing the unique angle or hook for this piece.

## Recommended Word Count
A target word count with brief justification.

## Article Outline
A numbered list of H2 headings, each with a 1-2 sentence description of what to cover.

## Internal Link Opportunities
Based on the crawled site content provided, suggest specific pages to link to with recommended anchor text. Format as a bullet list: [anchor text](url).

Ground every recommendation in the GSC data and crawled site content provided. Be specific, not generic.`, req.Keyword),
		})

		fullText, err := ai.StreamHandler.StreamResponse(r.Context(), w, ai.FlashProvider, messages)
		if err != nil {
			s.logger.Error("Brief generation stream failed", zap.Error(err))
			return
		}

		brief := map[string]interface{}{
			"id":         uuid.New().String(),
			"project_id": req.ProjectID,
			"keyword":    req.Keyword,
			"brief_data": fullText,
			"status":     "draft",
			"created_by": userID,
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		}

		data, _ := json.Marshal(brief)
		s.serviceRole.From("content_briefs").Insert(json.RawMessage(data), false, "", "", "").Execute()

		// Log usage
		go ai.UsageTracker.Log(context.Background(), req.ProjectID, userID, "brief", "gemini-2.5-flash", 0, 0)
		go ai.MemoryService.ExtractAndStore(context.Background(), req.ProjectID, "brief", fullText)
	}
}

func (s *Server) handleListBriefs(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	data, _, err := s.serviceRole.From("content_briefs").
		Select("*", "", false).
		Eq("project_id", projectID).
		Order("created_at", nil).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch briefs (table may not exist yet)", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var briefs []map[string]interface{}
	json.Unmarshal(data, &briefs)
	s.respondJSON(w, http.StatusOK, briefs)
}

func (s *Server) handleUpdateBrief(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	briefID, _ := req["id"].(string)
	projectID, _ := req["project_id"].(string)
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	req["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	data, _ := json.Marshal(req)

	_, _, err := s.serviceRole.From("content_briefs").
		Update(json.RawMessage(data), "", "").
		Eq("id", briefID).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to update brief")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ---- Article Writer Handlers ----

func (s *Server) handleGenerateArticle(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req struct {
			ProjectID string `json:"project_id"`
			BriefID   string `json:"brief_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if ok, err := s.verifyProjectAccess(userID, req.ProjectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		status, err := ai.UsageTracker.CheckLimit(r.Context(), req.ProjectID, "article")
		if err == nil && !status.Allowed {
			s.respondError(w, http.StatusTooManyRequests, status.Message)
			return
		}

		// Load the brief
		briefData, _, err := s.serviceRole.From("content_briefs").
			Select("*", "", false).
			Eq("id", req.BriefID).
			Limit(1, "").
			Execute()
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to load brief")
			return
		}

		var briefs []map[string]interface{}
		json.Unmarshal(briefData, &briefs)
		if len(briefs) == 0 {
			s.respondError(w, http.StatusNotFound, "Brief not found")
			return
		}

		brief := briefs[0]
		keyword, _ := brief["keyword"].(string)

		messages, err := ai.ContextBuilder.BuildMessages(r.Context(), req.ProjectID, aicontext.ContextArticle, aicontext.Params{
			Keyword: keyword,
			TopN:    5,
		})
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to build context")
			return
		}

		briefJSON, _ := json.Marshal(brief["brief_data"])
		messages = append(messages, providers.Message{
			Role: "user",
			Content: fmt.Sprintf(`Write a full, publish-ready article based on this content brief:

%s

Write in markdown format. Include all H2 sections from the outline. Suggest internal links inline using [anchor text](url) format. Match the writing voice profile exactly.`, string(briefJSON)),
		})

		fullText, err := ai.StreamHandler.StreamResponse(r.Context(), w, ai.FlashProvider, messages)
		if err != nil {
			s.logger.Error("Article generation stream failed", zap.Error(err))
			return
		}

		wordCount := len(strings.Fields(fullText))

		article := map[string]interface{}{
			"id":         uuid.New().String(),
			"project_id": req.ProjectID,
			"brief_id":   req.BriefID,
			"content":    fullText,
			"word_count": wordCount,
			"status":     "draft",
			"created_by": userID,
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		}

		data, _ := json.Marshal(article)
		s.serviceRole.From("articles").Insert(json.RawMessage(data), false, "", "", "").Execute()

		go ai.UsageTracker.Log(context.Background(), req.ProjectID, userID, "article", "gemini-2.5-flash", 0, 0)
		go ai.MemoryService.ExtractAndStore(context.Background(), req.ProjectID, "article", fullText)
	}
}

func (s *Server) handleListArticles(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	data, _, err := s.serviceRole.From("articles").
		Select("*", "", false).
		Eq("project_id", projectID).
		Order("created_at", nil).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch articles (table may not exist yet)", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var articles []map[string]interface{}
	json.Unmarshal(data, &articles)
	s.respondJSON(w, http.StatusOK, articles)
}

func (s *Server) handleUpdateArticle(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	articleID, _ := req["id"].(string)
	projectID, _ := req["project_id"].(string)
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	req["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	data, _ := json.Marshal(req)

	_, _, err := s.serviceRole.From("articles").
		Update(json.RawMessage(data), "", "").
		Eq("id", articleID).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to update article")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ---- GSC Intelligence Dashboard Handlers ----

// gscRowToFlat extracts dimension_value and flattens nested metrics into top-level fields.
// gsc_performance_rows stores: row_type, dimension_value (text), metrics (jsonb).
func gscRowToFlat(row map[string]interface{}) map[string]interface{} {
	flat := map[string]interface{}{}

	rowType, _ := row["row_type"].(string)
	dimVal, _ := row["dimension_value"].(string)

	switch rowType {
	case "page":
		flat["page"] = dimVal
		flat["query"] = ""
	default:
		flat["query"] = dimVal
		flat["page"] = ""
	}

	if metrics, ok := row["metrics"].(map[string]interface{}); ok {
		for k, v := range metrics {
			flat[k] = v
		}
	}

	return flat
}

func (s *Server) handleQuickWins(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	snapshotID := s.loadLatestGSCSnapshotID(projectID)
	if snapshotID == "" {
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	// Load GSC rows with query dimension data (latest snapshot only)
	data, _, err := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("snapshot_id", snapshotID).
		Eq("project_id", projectID).
		Eq("row_type", "query").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch GSC performance rows for quick wins", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var rawRows []map[string]interface{}
	if err := json.Unmarshal(data, &rawRows); err != nil {
		s.logger.Warn("Failed to parse GSC performance rows for quick wins", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	// Filter to keywords ranking 4-20 with at least 10 impressions, compute opportunity score.
	// The opportunity score naturally down-ranks low-traffic queries, so a permissive impression
	// floor avoids excluding legitimate opportunities on smaller sites.
	var results []map[string]interface{}
	for _, raw := range rawRows {
		flat := gscRowToFlat(raw)
		position := getFloat(flat["position"])
		impressions := getFloat(flat["impressions"])
		ctr := getFloat(flat["ctr"])

		if position >= 4 && position <= 20 && impressions >= 10 {
			if position > 0 {
				flat["opportunity_score"] = impressions * (1 - ctr) * (1 / position)
			}
			results = append(results, flat)
		}
	}

	// Sort by opportunity score descending
	sort.Slice(results, func(i, j int) bool {
		si, _ := results[i]["opportunity_score"].(float64)
		sj, _ := results[j]["opportunity_score"].(float64)
		return si > sj
	})

	if len(results) > 50 {
		results = results[:50]
	}

	s.respondJSON(w, http.StatusOK, results)
}

func (s *Server) handleDecliningPages(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	snapshotID := s.loadLatestGSCSnapshotID(projectID)
	if snapshotID == "" {
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	// Load GSC rows with page dimension data (latest snapshot only)
	data, _, err := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("snapshot_id", snapshotID).
		Eq("project_id", projectID).
		Eq("row_type", "page").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch GSC performance rows for declining pages", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var rawRows []map[string]interface{}
	if err := json.Unmarshal(data, &rawRows); err != nil {
		s.logger.Warn("Failed to parse GSC performance rows for declining pages", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var results []map[string]interface{}
	for _, raw := range rawRows {
		flat := gscRowToFlat(raw)
		results = append(results, flat)
	}

	// Sort by clicks descending
	sort.Slice(results, func(i, j int) bool {
		ci := getFloat(results[i]["clicks"])
		cj := getFloat(results[j]["clicks"])
		return ci > cj
	})

	if len(results) > 50 {
		results = results[:50]
	}

	s.respondJSON(w, http.StatusOK, results)
}

func (s *Server) handleKeywordSuggestions(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	snapshotID := s.loadLatestGSCSnapshotID(projectID)
	if snapshotID == "" {
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	// Load GSC query rows (latest snapshot only)
	data, _, err := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("snapshot_id", snapshotID).
		Eq("project_id", projectID).
		Eq("row_type", "query").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch GSC rows for keyword suggestions", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var rawRows []map[string]interface{}
	json.Unmarshal(data, &rawRows)

	// Load crawled page URLs to detect content gaps
	pageData, _, _ := s.serviceRole.From("crawled_pages").
		Select("url,title", "", false).
		Eq("project_id", projectID).
		Execute()

	crawledURLs := make(map[string]bool)
	var crawledPages []struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	}
	json.Unmarshal(pageData, &crawledPages)
	for _, p := range crawledPages {
		crawledURLs[p.URL] = true
	}

	// Also load existing brief keywords to exclude them
	briefData, _, _ := s.serviceRole.From("content_briefs").
		Select("keyword", "", false).
		Eq("project_id", projectID).
		Execute()

	existingKeywords := make(map[string]bool)
	var briefs []struct {
		Keyword string `json:"keyword"`
	}
	json.Unmarshal(briefData, &briefs)
	for _, b := range briefs {
		existingKeywords[strings.ToLower(b.Keyword)] = true
	}

	// Load GSC page rows to check which queries have ranking pages (latest snapshot only)
	pageRowData, _, _ := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("snapshot_id", snapshotID).
		Eq("project_id", projectID).
		Eq("row_type", "page").
		Execute()

	var pageRows []map[string]interface{}
	json.Unmarshal(pageRowData, &pageRows)

	rankedPages := make(map[string]bool)
	for _, row := range pageRows {
		dimVal, _ := row["dimension_value"].(string)
		if dimVal != "" {
			rankedPages[dimVal] = true
		}
	}

	var suggestions []map[string]interface{}
	for _, raw := range rawRows {
		flat := gscRowToFlat(raw)
		query, _ := flat["query"].(string)
		position, _ := flat["position"].(float64)
		impressions, _ := flat["impressions"].(float64)
		ctr, _ := flat["ctr"].(float64)
		clicks, _ := flat["clicks"].(float64)

		if query == "" || existingKeywords[strings.ToLower(query)] {
			continue
		}

		// Skip very short or branded-looking queries
		if len(query) < 3 {
			continue
		}

		var reason string
		var score float64

		if position >= 5 && position <= 20 && impressions >= 50 {
			// Quick Win: ranking but underperforming
			score = impressions * (1 - ctr) * (1 / position)
			reason = "quick_win"
		} else if impressions >= 100 && position > 20 {
			// Content Gap: high impressions but ranking poorly
			score = impressions * (1 / position)
			reason = "content_gap"
		} else if impressions >= 200 && clicks < 5 {
			// High visibility, low engagement
			score = impressions * 0.5
			reason = "low_engagement"
		} else {
			continue
		}

		suggestions = append(suggestions, map[string]interface{}{
			"query":       query,
			"position":    position,
			"impressions": impressions,
			"ctr":         ctr,
			"clicks":      clicks,
			"score":       score,
			"reason":      reason,
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		si, _ := suggestions[i]["score"].(float64)
		sj, _ := suggestions[j]["score"].(float64)
		return si > sj
	})

	if len(suggestions) > 15 {
		suggestions = suggestions[:15]
	}

	s.respondJSON(w, http.StatusOK, suggestions)
}

func (s *Server) handleExplainOpportunity(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req struct {
			ProjectID string `json:"project_id"`
			Query     string `json:"query"`
			Page      string `json:"page"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if ok, err := s.verifyProjectAccess(userID, req.ProjectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		messages, err := ai.ContextBuilder.BuildMessages(r.Context(), req.ProjectID, aicontext.ContextExplain, aicontext.Params{
			Keyword: req.Query,
			PageURL: req.Page,
		})
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to build context")
			return
		}

		messages = append(messages, providers.Message{
			Role:    "user",
			Content: fmt.Sprintf(`Explain why the keyword "%s" on page %s is underperforming. What is likely suppressing CTR? Provide a concise diagnosis and recommended next action.`, req.Query, req.Page),
		})

		fullText, err := ai.StreamHandler.StreamResponse(r.Context(), w, ai.FlashProvider, messages)
		if err != nil {
			s.logger.Error("Explain opportunity stream failed", zap.Error(err))
		}
		_ = fullText

		go ai.UsageTracker.Log(context.Background(), req.ProjectID, userID, "explain", "gemini-2.5-flash", 0, 0)
	}
}

func (s *Server) handleDiagnoseDecline(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req struct {
			ProjectID string `json:"project_id"`
			PageURL   string `json:"page_url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if ok, err := s.verifyProjectAccess(userID, req.ProjectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		messages, err := ai.ContextBuilder.BuildMessages(r.Context(), req.ProjectID, aicontext.ContextDiagnostic, aicontext.Params{
			PageURL: req.PageURL,
		})
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to build context")
			return
		}

		messages = append(messages, providers.Message{
			Role: "user",
			Content: fmt.Sprintf(`Diagnose the decline for page: %s

Rank likely causes by probability:
1. Content staleness
2. Keyword cannibalization
3. Algorithm sensitivity
4. Lost backlinks
5. Technical issues

Include a recommended next action.`, req.PageURL),
		})

		fullText, err := ai.StreamHandler.StreamResponse(r.Context(), w, ai.FlashProvider, messages)
		if err != nil {
			s.logger.Error("Diagnose decline stream failed", zap.Error(err))
		}
		_ = fullText

		go ai.UsageTracker.Log(context.Background(), req.ProjectID, userID, "diagnose", "gemini-2.5-flash", 0, 0)
	}
}

// ---- Content Gaps Handler ----

// handleContentGaps identifies topics where GSC shows meaningful impressions but
// the site has no page with strong topical coverage. For each high-impression query,
// it checks pgvector cosine similarity against crawled pages. Queries whose best
// matching page falls below a similarity threshold are flagged as content gaps.
func (s *Server) handleContentGaps(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		projectID := r.URL.Query().Get("project_id")
		if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		if ai == nil || ai.EmbedService == nil {
			s.respondError(w, http.StatusServiceUnavailable, "AI features not available")
			return
		}

		snapshotID := s.loadLatestGSCSnapshotID(projectID)
		if snapshotID == "" {
			s.respondJSON(w, http.StatusOK, []interface{}{})
			return
		}

		data, _, err := s.serviceRole.From("gsc_performance_rows").
			Select("*", "", false).
			Eq("snapshot_id", snapshotID).
			Eq("project_id", projectID).
			Eq("row_type", "query").
			Execute()
		if err != nil {
			s.logger.Warn("Failed to fetch GSC rows for content gaps", zap.Error(err))
			s.respondJSON(w, http.StatusOK, []interface{}{})
			return
		}

		var rawRows []map[string]interface{}
		json.Unmarshal(data, &rawRows)

		// Filter to queries with meaningful impressions
		type candidate struct {
			Query       string
			Impressions float64
			Clicks      float64
			CTR         float64
			Position    float64
		}
		var candidates []candidate
		for _, raw := range rawRows {
			flat := gscRowToFlat(raw)
			query, _ := flat["query"].(string)
			impressions, _ := flat["impressions"].(float64)
			clicks, _ := flat["clicks"].(float64)
			ctr, _ := flat["ctr"].(float64)
			position, _ := flat["position"].(float64)

			if query == "" || len(query) < 3 || impressions < 50 {
				continue
			}
			candidates = append(candidates, candidate{
				Query:       query,
				Impressions: impressions,
				Clicks:      clicks,
				CTR:         ctr,
				Position:    position,
			})
		}

		// Sort by impressions descending and cap to avoid excessive embedding calls
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Impressions > candidates[j].Impressions
		})
		if len(candidates) > 100 {
			candidates = candidates[:100]
		}

		// Batch embed all candidate queries
		queryTexts := make([]string, len(candidates))
		for i, c := range candidates {
			queryTexts[i] = c.Query
		}

		queryVectors, err := ai.EmbedService.EmbedBatch(r.Context(), queryTexts)
		if err != nil {
			s.logger.Error("Failed to embed queries for content gaps", zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to analyze content gaps")
			return
		}

		// For each query, check its best match against crawled pages via RPC
		const similarityThreshold = 0.45
		type contentGap struct {
			Query          string  `json:"query"`
			Impressions    float64 `json:"impressions"`
			Clicks         float64 `json:"clicks"`
			CTR            float64 `json:"ctr"`
			Position       float64 `json:"position"`
			BestMatchURL   string  `json:"best_match_url"`
			BestMatchTitle string  `json:"best_match_title"`
			Similarity     float64 `json:"similarity"`
			GapScore       float64 `json:"gap_score"`
		}

		var gaps []contentGap
		for i, c := range candidates {
			if i >= len(queryVectors) {
				break
			}
			vecStr := embeddings.FormatForPgvector(queryVectors[i])
			rpcBody := map[string]interface{}{
				"query_embedding":  vecStr,
				"match_project_id": projectID,
				"match_count":      1,
				"match_threshold":  0.0,
			}

			result := s.serviceRole.Rpc("match_crawled_pages", "", rpcBody)

			var bestMatch float64
			var matchURL, matchTitle string
			if result != "" {
				var matches []struct {
					URL        string  `json:"url"`
					Title      string  `json:"title"`
					Similarity float64 `json:"similarity"`
				}
				if json.Unmarshal([]byte(result), &matches) == nil && len(matches) > 0 {
					bestMatch = matches[0].Similarity
					matchURL = matches[0].URL
					matchTitle = matches[0].Title
				}
			}

			if bestMatch < similarityThreshold {
				gapMagnitude := similarityThreshold - bestMatch
				gaps = append(gaps, contentGap{
					Query:          c.Query,
					Impressions:    c.Impressions,
					Clicks:         c.Clicks,
					CTR:            c.CTR,
					Position:       c.Position,
					BestMatchURL:   matchURL,
					BestMatchTitle: matchTitle,
					Similarity:     bestMatch,
					GapScore:       c.Impressions * gapMagnitude,
				})
			}
		}

		sort.Slice(gaps, func(i, j int) bool {
			return gaps[i].GapScore > gaps[j].GapScore
		})
		if len(gaps) > 30 {
			gaps = gaps[:30]
		}

		s.respondJSON(w, http.StatusOK, gaps)
	}
}

// ---- Internal Link Suggester Handlers ----

// extractInternalLinksFromData extracts internal_links from a page's data jsonb field.
func extractInternalLinksFromData(dataVal interface{}) []string {
	if dataVal == nil {
		return nil
	}
	var dataField map[string]interface{}
	switch v := dataVal.(type) {
	case map[string]interface{}:
		dataField = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(jsonBytes, &dataField); err != nil {
			return nil
		}
	}
	if dataField == nil {
		return nil
	}
	linksVal := dataField["internal_links"]
	if linksVal == nil {
		return nil
	}
	if linkSlice, ok := linksVal.([]interface{}); ok {
		var out []string
		for _, l := range linkSlice {
			if linkStr, ok := l.(string); ok && linkStr != "" {
				out = append(out, linkStr)
			}
		}
		return out
	}
	if linkSlice, ok := linksVal.([]string); ok {
		return linkSlice
	}
	return nil
}

// latestCrawlID returns the most recent succeeded crawl ID for a project, or "" if none.
func (s *Server) latestCrawlID(projectID string) string {
	crawlData, _, err := s.serviceRole.
		From("crawls").
		Select("id", "", false).
		Eq("project_id", projectID).
		Eq("status", "succeeded").
		Order("started_at", nil).
		Limit(1, "").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to find latest crawl", zap.Error(err))
		return ""
	}
	var crawls []map[string]interface{}
	if err := json.Unmarshal(crawlData, &crawls); err != nil || len(crawls) == 0 {
		return ""
	}
	id, _ := crawls[0]["id"].(string)
	return id
}

func (s *Server) handleInternalLinkSuggestions(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	pageURL := r.URL.Query().Get("page_url")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	crawlID := s.latestCrawlID(projectID)
	if crawlID == "" {
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	query := s.serviceRole.From("pages").
		Select("url,title,meta_description,data", "", false).
		Eq("crawl_id", crawlID).
		Limit(50, "")

	if pageURL != "" {
		query = query.Neq("url", pageURL)
	}

	data, _, err := query.Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch pages for link suggestions", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var pages []map[string]interface{}
	json.Unmarshal(data, &pages)

	// Build link counts: how many other pages link to each page (internal_links is in data jsonb)
	inboundCounts := make(map[string]int)
	for _, p := range pages {
		internalLinks := extractInternalLinksFromData(p["data"])
		for _, linkURL := range internalLinks {
			if linkURL != "" {
				inboundCounts[linkURL]++
			}
		}
	}

	// Return pages sorted by fewest inbound links (best candidates for more internal links)
	type suggestion struct {
		URL             string `json:"url"`
		Title           string `json:"title"`
		MetaDescription string `json:"meta_description"`
		InboundLinks    int    `json:"inbound_links"`
	}

	var suggestions []suggestion
	for _, p := range pages {
		url, _ := p["url"].(string)
		title, _ := p["title"].(string)
		meta, _ := p["meta_description"].(string)
		suggestions = append(suggestions, suggestion{
			URL:             url,
			Title:           title,
			MetaDescription: meta,
			InboundLinks:    inboundCounts[url],
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].InboundLinks < suggestions[j].InboundLinks
	})

	if len(suggestions) > 20 {
		suggestions = suggestions[:20]
	}

	s.respondJSON(w, http.StatusOK, suggestions)
}

func (s *Server) handleOrphanedPages(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	crawlID := s.latestCrawlID(projectID)
	if crawlID == "" {
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	data, _, err := s.serviceRole.From("pages").
		Select("url,title,data", "", false).
		Eq("crawl_id", crawlID).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch pages for orphan check", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var pages []map[string]interface{}
	json.Unmarshal(data, &pages)

	// Build set of all internally linked URLs (internal_links is in data jsonb)
	linkedURLs := make(map[string]bool)
	for _, page := range pages {
		for _, linkURL := range extractInternalLinksFromData(page["data"]) {
			if linkURL != "" {
				linkedURLs[linkURL] = true
			}
		}
	}

	// Filter to orphaned pages (not linked to by any other page)
	var orphaned []map[string]interface{}
	for _, page := range pages {
		url, _ := page["url"].(string)
		if url != "" && !linkedURLs[url] {
			title, _ := page["title"].(string)
			orphaned = append(orphaned, map[string]interface{}{
				"url":   url,
				"title": title,
			})
		}
	}

	s.respondJSON(w, http.StatusOK, orphaned)
}

// ---- Site Crawl Handler ----

func (s *Server) handleSiteCrawl(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var req struct {
			ProjectID string `json:"project_id"`
			SiteURL   string `json:"site_url"`
			MaxPages  int    `json:"max_pages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if ok, err := s.verifyProjectAccess(userID, req.ProjectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		if req.MaxPages == 0 {
			req.MaxPages = 200
		}

		crawled, err := ai.SiteCrawler.CrawlFromSitemap(r.Context(), req.ProjectID, req.SiteURL, req.MaxPages)
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Crawl failed: %v", err))
			return
		}

		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"crawled": crawled,
			"status":  "complete",
		})
	}
}

// ---- Usage Tracking Handler ----

func (s *Server) handleUsageSummary(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		projectID := r.URL.Query().Get("project_id")
		if ok, err := s.verifyProjectAccess(userID, projectID); !ok || err != nil {
			s.respondError(w, http.StatusForbidden, "Access denied")
			return
		}

		summary, err := ai.UsageTracker.GetMonthlyUsageSummary(r.Context(), projectID)
		if err != nil {
			s.logger.Warn("Failed to fetch AI usage (table may not exist yet)", zap.Error(err))
			summary = map[string]int{}
		}

		// Get plan limits
		planTier := "starter"
		projectData, _, _ := s.serviceRole.From("projects").
			Select("plan_tier", "", false).
			Eq("id", projectID).
			Limit(1, "").
			Execute()
		var projects []struct {
			PlanTier string `json:"plan_tier"`
		}
		json.Unmarshal(projectData, &projects)
		if len(projects) > 0 && projects[0].PlanTier != "" {
			planTier = projects[0].PlanTier
		}

		limits := usage.Plans[planTier]
		if _, ok := usage.Plans[planTier]; !ok {
			limits = usage.Plans["starter"]
		}

		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"plan":  planTier,
			"usage": summary,
			"limits": map[string]int{
				"briefs_per_month":      limits.BriefsPerMonth,
				"articles_per_month":    limits.ArticlesPerMonth,
				"diagnostics_per_month": limits.DiagnosticsPerMonth,
			},
		})
	}
}

// ---- Helpers ----

func parseVoiceComponents(text string) map[string]interface{} {
	components := map[string]interface{}{}
	sections := map[string]string{
		"TONE:":           "tone",
		"STRUCTURE:":      "structure",
		"SENTENCE_STYLE:": "sentence_style",
		"BRAND_CONTEXT:":  "brand_context",
		"AVOID_LIST:":     "avoid_list",
	}

	for label, field := range sections {
		idx := strings.Index(strings.ToUpper(text), strings.ToUpper(label))
		if idx == -1 {
			continue
		}
		start := idx + len(label)
		end := len(text)

		// Find the next section start
		for otherLabel := range sections {
			if otherLabel == label {
				continue
			}
			otherIdx := strings.Index(strings.ToUpper(text[start:]), strings.ToUpper(otherLabel))
			if otherIdx != -1 && start+otherIdx < end {
				end = start + otherIdx
			}
		}

		components[field] = strings.TrimSpace(text[start:end])
	}

	return components
}

// ---- Weekly Digest Job ----

func (s *Server) handleWeeklyDigest(ai *AISuiteServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		secret := r.Header.Get("X-Cron-Secret")
		if secret == "" {
			secret = r.URL.Query().Get("secret")
		}
		if s.cronSecret == "" || secret != s.cronSecret {
			s.respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Get all projects
		data, _, err := s.serviceRole.From("projects").
			Select("id", "", false).
			Execute()
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, "Failed to load projects")
			return
		}

		var projects []struct {
			ID string `json:"id"`
		}
		json.Unmarshal(data, &projects)

		generated := 0
		for _, proj := range projects {
			ctx := r.Context()

			messages, err := ai.ContextBuilder.BuildMessages(ctx, proj.ID, aicontext.ContextDigest, aicontext.Params{TopN: 10})
			if err != nil {
				s.logger.Warn("Digest context build failed", zap.String("project_id", proj.ID), zap.Error(err))
				continue
			}

			messages = append(messages, providers.Message{
				Role:    "user",
				Content: "Generate a concise weekly SEO digest. Include the biggest GSC movers (impressions, clicks, position changes), any new keyword appearances, and one prioritized action for the coming week. Keep it to 3-5 bullet points.",
			})

			response, err := ai.LiteProvider.Completion(ctx, messages)
			if err != nil {
				s.logger.Warn("Digest generation failed", zap.String("project_id", proj.ID), zap.Error(err))
				continue
			}

			digest := map[string]interface{}{
				"id":           uuid.New().String(),
				"project_id":   proj.ID,
				"content":      response,
				"generated_at": time.Now().UTC().Format(time.RFC3339),
			}
			digestData, _ := json.Marshal(digest)
			s.serviceRole.From("digests").Insert(json.RawMessage(digestData), false, "", "", "").Execute()
			generated++
		}

		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"digests_generated": generated,
		})
	}
}

// ---- Helpers ----

func extractJSON(text string) interface{} {
	// Try to extract JSON from the response (may be wrapped in markdown code blocks)
	cleaned := text
	if idx := strings.Index(cleaned, "```json"); idx != -1 {
		cleaned = cleaned[idx+7:]
	} else if idx := strings.Index(cleaned, "```"); idx != -1 {
		cleaned = cleaned[idx+3:]
	}
	if idx := strings.LastIndex(cleaned, "```"); idx != -1 {
		cleaned = cleaned[:idx]
	}
	cleaned = strings.TrimSpace(cleaned)

	var result interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return cleaned
	}
	return result
}
