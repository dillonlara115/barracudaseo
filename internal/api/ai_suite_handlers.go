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

	ctxBuilder := aicontext.NewBuilder(s.serviceRole, s.logger)
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
			Content: fmt.Sprintf(`Generate a content brief for the keyword: "%s"

Output as JSON with these fields:
- target_keyword (string)
- secondary_keywords (array of 3-5 strings)
- recommended_title (string)
- meta_description (string, under 160 chars)
- recommended_word_count (number)
- outline (array of objects: {heading: string, description: string})
- internal_links (array of objects: {url: string, anchor_text: string})
- content_angle (string)
- intent (string: informational/navigational/commercial/transactional)`, req.Keyword),
		})

		fullText, err := ai.StreamHandler.StreamResponse(r.Context(), w, ai.FlashProvider, messages)
		if err != nil {
			s.logger.Error("Brief generation stream failed", zap.Error(err))
			return
		}

		// Store the brief
		brief := map[string]interface{}{
			"id":         uuid.New().String(),
			"project_id": req.ProjectID,
			"keyword":    req.Keyword,
			"brief_data": extractJSON(fullText),
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

	// Load GSC rows with query dimension data
	data, _, err := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", "query").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch GSC performance rows for quick wins", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var rawRows []map[string]interface{}
	json.Unmarshal(data, &rawRows)

	// Filter to keywords ranking 5-15 with decent impressions, compute opportunity score
	var results []map[string]interface{}
	for _, raw := range rawRows {
		flat := gscRowToFlat(raw)
		position, _ := flat["position"].(float64)
		impressions, _ := flat["impressions"].(float64)
		ctr, _ := flat["ctr"].(float64)

		if position >= 5 && position <= 15 && impressions >= 100 {
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

	// Load GSC rows with page dimension data
	data, _, err := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", "page").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch GSC performance rows for declining pages", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var rawRows []map[string]interface{}
	json.Unmarshal(data, &rawRows)

	var results []map[string]interface{}
	for _, raw := range rawRows {
		flat := gscRowToFlat(raw)
		results = append(results, flat)
	}

	// Sort by clicks descending
	sort.Slice(results, func(i, j int) bool {
		ci, _ := results[i]["clicks"].(float64)
		cj, _ := results[j]["clicks"].(float64)
		return ci > cj
	})

	if len(results) > 50 {
		results = results[:50]
	}

	s.respondJSON(w, http.StatusOK, results)
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

// ---- Internal Link Suggester Handlers ----

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

	data, _, err := s.serviceRole.From("crawled_pages").
		Select("url,title,meta_description,word_count", "", false).
		Eq("project_id", projectID).
		Neq("url", pageURL).
		Order("word_count", nil).
		Limit(20, "").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch link suggestions (crawled_pages may lack columns)", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var pages []map[string]interface{}
	json.Unmarshal(data, &pages)
	s.respondJSON(w, http.StatusOK, pages)
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

	data, _, err := s.serviceRole.From("crawled_pages").
		Select("url,title,word_count", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to fetch pages for orphan check", zap.Error(err))
		s.respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	var pages []map[string]interface{}
	json.Unmarshal(data, &pages)

	// Build set of all internally linked URLs
	linkedURLs := make(map[string]bool)
	allLinksData, _, _ := s.serviceRole.From("crawled_pages").
		Select("internal_links", "", false).
		Eq("project_id", projectID).
		Execute()

	var allLinks []struct {
		InternalLinks []string `json:"internal_links"`
	}
	json.Unmarshal(allLinksData, &allLinks)
	for _, page := range allLinks {
		for _, link := range page.InternalLinks {
			linkedURLs[link] = true
		}
	}

	// Filter to orphaned pages
	var orphaned []map[string]interface{}
	for _, page := range pages {
		url, _ := page["url"].(string)
		if !linkedURLs[url] {
			orphaned = append(orphaned, page)
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
