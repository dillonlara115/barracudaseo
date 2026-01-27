package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/dillonlara115/barracuda/internal/utils"
	"github.com/dillonlara115/barracuda/pkg/models"
)

type apiClient struct {
	baseURL string
	token   string
}

type apiProject struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type apiProjectsResponse struct {
	Projects []apiProject `json:"projects"`
	Count    int          `json:"count"`
}

type apiCreateProjectRequest struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type apiCreateCrawlRequest struct {
	ProjectID string               `json:"project_id"`
	Pages     []*models.PageResult `json:"pages"`
	Source    string               `json:"source"`
}

type apiCreateCrawlResponse struct {
	CrawlID     string `json:"crawl_id"`
	ProjectID   string `json:"project_id"`
	TotalPages  int    `json:"total_pages"`
	TotalIssues int    `json:"total_issues"`
	Status      string `json:"status"`
}

func newAPIClient(ctx context.Context) (*apiClient, *utils.Credentials, error) {
	loadEnv()

	creds, err := utils.LoadCredentials()
	if err != nil {
		return nil, nil, err
	}
	if creds == nil {
		return nil, nil, errors.New("not authenticated. Run `barracuda auth login`")
	}

	supabaseURL, supabaseAnonKey := resolveSupabaseConfig()
	if creds.SupabaseURL == "" && supabaseURL != "" {
		creds.SupabaseURL = supabaseURL
	}
	if creds.SupabaseAnonKey == "" && supabaseAnonKey != "" {
		creds.SupabaseAnonKey = supabaseAnonKey
	}

	token, updated, err := utils.EnsureValidAccessToken(ctx, creds)
	if err != nil {
		return nil, nil, err
	}

	apiURL := strings.TrimSuffix(updated.APIURL, "/")
	if apiURL == "" {
		apiURL = strings.TrimSuffix(resolveAPIURL(), "/")
		updated.APIURL = apiURL
		_ = utils.SaveCredentials(updated)
	}

	return &apiClient{
		baseURL: apiURL,
		token:   token,
	}, updated, nil
}

func (c *apiClient) doJSON(ctx context.Context, method, path string, payload interface{}, out interface{}) error {
	var body *bytes.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func listProjects(ctx context.Context, c *apiClient) ([]apiProject, error) {
	var resp apiProjectsResponse
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/projects", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Projects, nil
}

func createProject(ctx context.Context, c *apiClient, name, domain string) (*apiProject, error) {
	payload := apiCreateProjectRequest{Name: name, Domain: domain}
	var project apiProject
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/projects", payload, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func uploadCrawl(ctx context.Context, c *apiClient, projectID string, pages []*models.PageResult) (*apiCreateCrawlResponse, error) {
	payload := apiCreateCrawlRequest{
		ProjectID: projectID,
		Pages:     pages,
		Source:    "cli",
	}
	var resp apiCreateCrawlResponse
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/crawls", payload, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func ensureProjectSelection(ctx context.Context, c *apiClient, startURL string, projectID string) (string, error) {
	if projectID != "" {
		return projectID, nil
	}

	projects, err := listProjects(ctx, c)
	if err != nil {
		return "", err
	}

	domain := deriveDomain(startURL)

	if len(projects) == 0 {
		fmt.Println("No projects found. Let's create one.")
		return promptCreateProject(ctx, c, domain)
	}

	sort.Slice(projects, func(i, j int) bool {
		return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
	})

	choices := make([]string, 0, len(projects)+1)
	choiceToID := make(map[string]string, len(projects))
	for _, project := range projects {
		label := fmt.Sprintf("%s (%s)", project.Name, project.Domain)
		choices = append(choices, label)
		choiceToID[label] = project.ID
	}
	createLabel := "Create new project..."
	choices = append(choices, createLabel)

	selected, err := utils.PromptSelect("Select a project for this crawl", choices, choices[0])
	if err != nil {
		return "", err
	}
	if selected == createLabel {
		return promptCreateProject(ctx, c, domain)
	}
	if id, ok := choiceToID[selected]; ok {
		return id, nil
	}
	return "", errors.New("invalid project selection")
}

func promptCreateProject(ctx context.Context, c *apiClient, defaultDomain string) (string, error) {
	name, err := utils.PromptString("Project name", defaultDomain, true)
	if err != nil {
		return "", err
	}
	domain, err := utils.PromptString("Project domain", defaultDomain, true)
	if err != nil {
		return "", err
	}
	project, err := createProject(ctx, c, name, domain)
	if err != nil {
		return "", err
	}
	fmt.Printf("✓ Created project %s (%s)\n", project.Name, project.Domain)
	return project.ID, nil
}

func deriveDomain(startURL string) string {
	if startURL == "" {
		return ""
	}
	parsed, err := url.Parse(startURL)
	if err != nil {
		return ""
	}
	host := parsed.Hostname()
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}
	return host
}
