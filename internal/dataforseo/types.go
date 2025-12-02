package dataforseo

// OrganicTaskPost represents a request to create an organic SERP task
type OrganicTaskPost struct {
	LanguageName string `json:"language_name"`
	LocationName string `json:"location_name"`
	Keyword      string `json:"keyword"`
	Device       string `json:"device,omitempty"`        // "desktop" or "mobile"
	SearchEngine string `json:"search_engine_name,omitempty"` // "google.com" etc.
}

// OrganicTaskPostRequest is the request format for DataForSEO API
type OrganicTaskPostRequest map[string]OrganicTaskPost

// TaskResponse represents a single task in the response
type TaskResponse struct {
	ID         string `json:"id"`
	StatusCode int    `json:"status_code"`
	StatusMessage string `json:"status_message,omitempty"`
}

// OrganicTaskPostResponse represents the response from task_post endpoint
type OrganicTaskPostResponse struct {
	Version string         `json:"version"`
	StatusCode int         `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Tasks []TaskResponse  `json:"tasks"`
}

// OrganicResultItem represents a single organic result item
type OrganicResultItem struct {
	RankAbsolute int      `json:"rank_absolute"` // Overall position in SERP
	RankGroup   int      `json:"rank_group"`     // Organic-only position
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	SERPFeatures []string `json:"serp_features,omitempty"`
}

// OrganicResult represents the result data for a task
type OrganicResult struct {
	Items []OrganicResultItem `json:"items"`
}

// TaskResult represents a task result
type TaskResult struct {
	ID     string         `json:"id"`
	StatusCode int        `json:"status_code"`
	StatusMessage string  `json:"status_message"`
	Result []OrganicResult `json:"result"`
}

// OrganicTaskGetResponse represents the response from task_get endpoint
type OrganicTaskGetResponse struct {
	Version string       `json:"version"`
	StatusCode int       `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Tasks []TaskResult   `json:"tasks"`
}

// TasksReadyItem represents a task ID that's ready for retrieval
type TasksReadyItem struct {
	ID string `json:"id"`
}

// OrganicTasksReadyResponse represents the response from tasks_ready endpoint
type OrganicTasksReadyResponse struct {
	Version string       `json:"version"`
	StatusCode int       `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Tasks []TasksReadyItem `json:"tasks"`
}

// OrganicLiveRequest represents a request to the Live API (returns results immediately)
type OrganicLiveRequest map[string]OrganicTaskPost

// OrganicLiveResponse represents the response from Live API endpoint
// This is the same structure as OrganicTaskGetResponse since Live API returns results directly
type OrganicLiveResponse = OrganicTaskGetResponse

// RankedKeywordsTask represents a request to get ranked keywords for a domain/URL
// Based on DataForSEO API docs: https://docs.dataforseo.com/v3/dataforseo_labs/google/ranked_keywords/live/
type RankedKeywordsTask struct {
	Target           string `json:"target"`                     // Domain (e.g., "dataforseo.com") or URL
	LocationName     string `json:"location_name"`               // e.g., "United States"
	LanguageName     string `json:"language_name"`               // e.g., "English"
	LoadRankAbsolute bool  `json:"load_rank_absolute,omitempty"` // Load absolute rank
	LoadKeywordInfo  bool  `json:"load_keyword_info,omitempty"`  // Load keyword metrics (search volume, competition, CPC)
	Limit            int    `json:"limit,omitempty"`             // Max results (default 1000)
	// Note: Filters, SortBy, OrderBy, Offset may not be supported in live endpoint
	// Only use fields shown in official API examples
}

// RankedKeywordsRequest is the request format for Ranked Keywords API
// Note: This API uses an array format, not a map like SERP API
type RankedKeywordsRequest []RankedKeywordsTask

// RankedKeywordItem represents a single discovered keyword
// Based on actual API response structure
type RankedKeywordItem struct {
	SEType        string `json:"se_type"`
	LocationCode  int    `json:"location_code"`
	LanguageCode  string `json:"language_code"`
	
	// Keyword data is nested under keyword_data
	KeywordData struct {
		SEType        string `json:"se_type"`
		Keyword       string `json:"keyword"`
		LocationCode  int    `json:"location_code"`
		LanguageCode  string `json:"language_code"`
		KeywordInfo   struct {
			SEType          string   `json:"se_type"`
			LastUpdatedTime string   `json:"last_updated_time"`
			Competition     interface{} `json:"competition"` // Can be null, number, or string - we use competition_level instead
			CompetitionLevel string `json:"competition_level"`
			CPC             *float64 `json:"cpc"` // Can be null
			SearchVolume    int     `json:"search_volume"`
		} `json:"keyword_info"`
	} `json:"keyword_data"`
	
	RankedSERPElement struct {
		SERPItem struct {
			RankGroup    int    `json:"rank_group"`
			RankAbsolute int    `json:"rank_absolute"`
			URL          string `json:"url"`
			Title        string `json:"title"`
			Type         string `json:"type"` // "organic", "featured_snippet", etc.
		} `json:"serp_item"`
		CheckURL         string `json:"check_url"`
		SERPItemTypes    []string `json:"serp_item_types"`
		KeywordDifficulty int    `json:"keyword_difficulty"`
		IsLost           bool   `json:"is_lost"`
		LastUpdatedTime  string `json:"last_updated_time"`
	} `json:"ranked_serp_element"`
}

// RankedKeywordsResult represents the result data for ranked keywords
type RankedKeywordsResult struct {
	Items []RankedKeywordItem `json:"items"`
}

// RankedKeywordsTaskResult represents a task result for ranked keywords
type RankedKeywordsTaskResult struct {
	ID       string                 `json:"id"`
	StatusCode int                  `json:"status_code"`
	StatusMessage string            `json:"status_message"`
	Result  []RankedKeywordsResult  `json:"result"`
}

// RankedKeywordsResponse represents the response from Ranked Keywords API
type RankedKeywordsResponse struct {
	Version      string                    `json:"version"`
	StatusCode   int                       `json:"status_code"`
	StatusMessage string                   `json:"status_message"`
	Tasks        []RankedKeywordsTaskResult `json:"tasks"`
}

