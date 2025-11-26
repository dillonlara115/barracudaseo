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

