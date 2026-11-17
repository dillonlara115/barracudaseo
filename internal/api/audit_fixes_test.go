package api

import (
	"reflect"
	"testing"
)

func TestPickFields_DropsEverythingNotAllowed(t *testing.T) {
	body := map[string]interface{}{
		"id":         "brief-1",
		"project_id": "someone-elses-project",
		"created_by": "attacker",
		"keyword":    "seo crawler",
		"status":     "approved",
	}

	got := pickFields(body, "keyword", "brief_data", "status")

	want := map[string]interface{}{"keyword": "seo crawler", "status": "approved"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pickFields = %v, want %v", got, want)
	}
}

func TestPublicReportIDFromPath(t *testing.T) {
	cases := map[string]string{
		// What the handler receives: the v1 router strips /api/v1
		"/reports/public/abc-123":        "abc-123",
		"/reports/public/abc-123/":       "abc-123",
		"/api/v1/reports/public/abc-123": "abc-123",
		"/reports/public/":               "",
	}
	for path, want := range cases {
		if got := publicReportIDFromPath(path); got != want {
			t.Errorf("publicReportIDFromPath(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestProjectsOrFilter(t *testing.T) {
	const (
		me       = "11111111-1111-1111-1111-111111111111"
		teammate = "22222222-2222-2222-2222-222222222222"
		shared   = "33333333-3333-3333-3333-333333333333"
	)

	got := projectsOrFilter([]string{me, teammate, me}, []string{shared, "x),owner_id.neq.(y"})
	want := "owner_id.in.(" + me + "," + teammate + "),id.in.(" + shared + ")"
	if got != want {
		t.Errorf("projectsOrFilter = %q, want %q", got, want)
	}

	if got := projectsOrFilter([]string{me}, nil); got != "owner_id.in.("+me+")" {
		t.Errorf("projectsOrFilter without memberships = %q", got)
	}
}

func TestAttachIssueURLs(t *testing.T) {
	pageURLs := map[int64]string{7: "https://example.com/pricing"}
	issues := []map[string]interface{}{
		{"type": "missing_h1", "page_id": float64(7)},
		{"type": "missing_security_header", "page_id": nil},
		{"type": "orphan_page", "page_id": float64(99)},
	}

	got := attachIssueURLs(issues, pageURLs)

	if len(got) != 1 || got[0]["url"] != "https://example.com/pricing" || got[0]["type"] != "missing_h1" {
		t.Errorf("attachIssueURLs = %v, want only the issue on page 7 with its URL", got)
	}
}

func TestAIPlanForTier(t *testing.T) {
	for tier, want := range map[string]string{"free": "starter", "": "starter", "pro": "pro", "team": "agency"} {
		if got := aiPlanForTier(tier); got != want {
			t.Errorf("aiPlanForTier(%q) = %q, want %q", tier, got, want)
		}
	}
}
