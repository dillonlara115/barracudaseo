package analyzer

import (
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

func f(value float64) *float64 { return &value }

func TestAnalyzeWebVitals_FieldData(t *testing.T) {
	slow := page("https://example.com/slow")
	slow.WebVitals = &models.WebVitals{Strategy: "mobile", FieldSource: "page",
		LCPMs: f(4500), INPMs: f(300), CLS: f(0.05)}

	good := page("https://example.com/good")
	good.WebVitals = &models.WebVitals{Strategy: "mobile", FieldSource: "page",
		LCPMs: f(2500), INPMs: f(200), CLS: f(0.1)}

	failed := page("https://example.com/failed")
	failed.WebVitals = &models.WebVitals{Strategy: "mobile", Error: "PageSpeed API error 500"}

	issues := AnalyzeWebVitals([]*models.PageResult{slow, good, failed, page("https://example.com/unmeasured")})

	if len(issues) != 2 {
		t.Fatalf("issues = %+v, want 2 (poor LCP and needs-improvement INP on the slow page)", issues)
	}
	lcp := issuesOfType(issues, IssueSlowLCP)
	if len(lcp) != 1 || lcp[0].Severity != "error" || lcp[0].Value != "4.5 s" ||
		!strings.Contains(lcp[0].Message, "real-user data for this page") {
		t.Errorf("slow_lcp = %+v, want an error of 4.5 s from page field data", lcp)
	}
	inp := issuesOfType(issues, IssueSlowINP)
	if len(inp) != 1 || inp[0].Severity != "warning" || inp[0].Value != "300 ms" {
		t.Errorf("slow_inp = %+v, want a warning of 300 ms", inp)
	}
}

func TestAnalyzeWebVitals_OriginDataReportedOnce(t *testing.T) {
	a := page("https://example.com/a")
	a.WebVitals = &models.WebVitals{Strategy: "mobile", FieldSource: "origin", CLS: f(0.3)}
	b := page("https://example.com/b")
	b.WebVitals = &models.WebVitals{Strategy: "mobile", FieldSource: "origin", CLS: f(0.3)}

	issues := AnalyzeWebVitals([]*models.PageResult{a, b})

	if len(issues) != 1 || issues[0].URL != a.URL || !strings.Contains(issues[0].Message, "whole site") {
		t.Errorf("issues = %+v, want one site-wide high_cls issue on the first page", issues)
	}
}

func TestAnalyzeWebVitals_LabDataIsOnlyAWarning(t *testing.T) {
	p := page("https://example.com/new")
	p.WebVitals = &models.WebVitals{Strategy: "mobile", LabLCPMs: f(6000), LabCLS: f(0.02), LabTBTMs: f(900)}

	issues := AnalyzeWebVitals([]*models.PageResult{p})

	if len(issues) != 1 || issues[0].Type != IssueSlowLCP || issues[0].Severity != "warning" ||
		!strings.Contains(issues[0].Message, "lab test") {
		t.Errorf("issues = %+v, want one lab-based slow_lcp warning and no INP issue", issues)
	}
}

func TestAnalyzeRendering_UnrenderedShell(t *testing.T) {
	shell := contentPage("https://example.com/app", 3, "h", "0000000000000000")
	shell.JS = &models.JSInfo{Framework: "React", ClientRendered: true, Signals: []string{"app mount element holds 0 words"}}

	ssr := contentPage("https://example.com/ssr", 800, "h2", "ffffffffffffffff")
	ssr.JS = &models.JSInfo{Framework: "Next.js"}

	results := []*models.PageResult{shell, ssr}
	issues := AnalyzeRendering(results)

	got := issuesOfType(issues, IssueJSDependentContent)
	if len(got) != 1 || got[0].URL != shell.URL || !strings.Contains(got[0].Message, "React app shell") {
		t.Errorf("js_dependent_content = %+v, want one for the React shell", got)
	}
	// The shell is short because it needs JavaScript, not because it is thin
	if thin := issuesOfType(AnalyzeContent(results), IssueThinContent); len(thin) != 0 {
		t.Errorf("thin_content = %+v, want none for a client-rendered shell", thin)
	}
}

func TestAnalyzeRendering_RenderedComparison(t *testing.T) {
	changed := page("https://example.com/app")
	changed.JS = &models.JSInfo{ClientRendered: true}
	changed.Rendered = &models.RenderComparison{
		RawWordCount: 5, RenderedWordCount: 600,
		RawTitle: "Loading...", RenderedTitle: "Real Title",
		RawCanonical: "https://example.com/", RenderedCanonical: "https://example.com/app",
		RawMetaRobots: "noindex", RenderedMetaRobots: "index, follow",
		LinksOnlyRendered: []string{"https://example.com/a", "https://example.com/b"},
	}

	same := page("https://example.com/static")
	same.Rendered = &models.RenderComparison{
		RawWordCount: 500, RenderedWordCount: 540,
		RawTitle: "Static", RenderedTitle: "Static",
		RawCanonical: "https://example.com/static", RenderedCanonical: "https://example.com/static/",
	}

	failed := page("https://example.com/timeout")
	failed.Rendered = &models.RenderComparison{Error: "failed to load"}

	issues := AnalyzeRendering([]*models.PageResult{changed, same, failed})

	want := map[IssueType]string{
		IssueContentRequiresJS:      "warning",
		IssueLinksRequireJS:         "warning",
		IssueRenderedCanonicalDiff:  "error",
		IssueRenderedRobotsMismatch: "error",
		IssueRenderedTitleMismatch:  "info",
	}
	if len(issues) != len(want) {
		t.Fatalf("issues = %+v, want %d, all on the changed page", issues, len(want))
	}
	for _, issue := range issues {
		if issue.URL != changed.URL {
			t.Errorf("unexpected issue on %s: %+v", issue.URL, issue)
		}
		if severity, ok := want[issue.Type]; !ok || severity != issue.Severity {
			t.Errorf("issue %s has severity %q, want %q", issue.Type, issue.Severity, severity)
		}
	}
	// Rendering confirmed the dependence, so the raw-HTML guess is not repeated
	if guess := issuesOfType(issues, IssueJSDependentContent); len(guess) != 0 {
		t.Errorf("js_dependent_content = %+v, want none once the page was rendered", guess)
	}
}
