package psi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// fieldResponse mirrors a PageSpeed response for a URL with real-user data.
const fieldResponse = `{
  "loadingExperience": {
    "metrics": {
      "LARGEST_CONTENTFUL_PAINT_MS": {"percentile": 3100, "category": "AVERAGE"},
      "INTERACTION_TO_NEXT_PAINT": {"percentile": 180, "category": "FAST"},
      "CUMULATIVE_LAYOUT_SHIFT_SCORE": {"percentile": 27, "category": "SLOW"}
    }
  },
  "lighthouseResult": {
    "audits": {
      "largest-contentful-paint": {"numericValue": 4212.5},
      "cumulative-layout-shift": {"numericValue": 0.31},
      "total-blocking-time": {"numericValue": 420}
    },
    "categories": {"performance": {"score": 0.615}}
  }
}`

func serve(t *testing.T, status int, body string, check func(r *http.Request)) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if check != nil {
			check(r)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	client := NewClient("test-key")
	client.Endpoint = server.URL
	return client
}

func TestAnalyze_FieldAndLabData(t *testing.T) {
	client := serve(t, http.StatusOK, fieldResponse, func(r *http.Request) {
		query := r.URL.Query()
		if query.Get("url") != "https://example.com/page" || query.Get("strategy") != "mobile" ||
			query.Get("category") != "performance" || query.Get("key") != "test-key" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
	})

	vitals, err := client.Analyze(context.Background(), "https://example.com/page", "mobile")
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}

	if vitals.FieldSource != "page" || vitals.Strategy != "mobile" {
		t.Errorf("FieldSource = %q, Strategy = %q; want page, mobile", vitals.FieldSource, vitals.Strategy)
	}
	if vitals.LCPMs == nil || *vitals.LCPMs != 3100 {
		t.Errorf("LCPMs = %v, want 3100", vitals.LCPMs)
	}
	if vitals.INPMs == nil || *vitals.INPMs != 180 {
		t.Errorf("INPMs = %v, want 180", vitals.INPMs)
	}
	if vitals.CLS == nil || *vitals.CLS != 0.27 {
		t.Errorf("CLS = %v, want 0.27 (the API reports it multiplied by 100)", vitals.CLS)
	}
	if vitals.LabLCPMs == nil || *vitals.LabLCPMs != 4212.5 || vitals.LabTBTMs == nil || *vitals.LabTBTMs != 420 {
		t.Errorf("lab values = %v, %v", vitals.LabLCPMs, vitals.LabTBTMs)
	}
	if vitals.PerformanceScore == nil || *vitals.PerformanceScore != 62 {
		t.Errorf("PerformanceScore = %v, want 62", vitals.PerformanceScore)
	}
}

func TestAnalyze_OriginFallbackAndNoFieldData(t *testing.T) {
	origin := serve(t, http.StatusOK, `{"loadingExperience": {"origin_fallback": true,
		"metrics": {"LARGEST_CONTENTFUL_PAINT_MS": {"percentile": 2000}}}}`, nil)
	vitals, err := origin.Analyze(context.Background(), "https://example.com/", "mobile")
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}
	if vitals.FieldSource != "origin" {
		t.Errorf("FieldSource = %q, want origin", vitals.FieldSource)
	}
	if vitals.INPMs != nil || vitals.CLS != nil {
		t.Error("metrics absent from the response should stay nil")
	}

	labOnly := serve(t, http.StatusOK, `{"loadingExperience": {"metrics": {}},
		"lighthouseResult": {"audits": {"largest-contentful-paint": {"numericValue": 1800}}, "categories": {}}}`, nil)
	vitals, err = labOnly.Analyze(context.Background(), "https://example.com/new", "desktop")
	if err != nil {
		t.Fatalf("Analyze error: %v", err)
	}
	if vitals.FieldSource != "" || vitals.LabLCPMs == nil || *vitals.LabLCPMs != 1800 {
		t.Errorf("vitals = %+v, want lab data only", vitals)
	}
}

func TestAnalyze_APIError(t *testing.T) {
	client := serve(t, http.StatusTooManyRequests,
		`{"error": {"code": 429, "message": "Quota exceeded for quota metric 'Queries'\nmore detail"}}`, nil)

	_, err := client.Analyze(context.Background(), "https://example.com/", "mobile")
	if err == nil || !strings.Contains(err.Error(), "429") || !strings.Contains(err.Error(), "Quota exceeded") {
		t.Errorf("error = %v, want the API's code and message", err)
	}
	if strings.Contains(err.Error(), "more detail") {
		t.Errorf("error = %v, want only the first line of the message", err)
	}
}

func TestEnrich_RecordsFailuresPerPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Query().Get("url"), "bad") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": {"code": 500, "message": "Lighthouse returned error"}}`))
			return
		}
		_, _ = w.Write([]byte(fieldResponse))
	}))
	defer server.Close()
	client := NewClient("")
	client.Endpoint = server.URL

	good := &models.PageResult{URL: "https://example.com/good"}
	bad := &models.PageResult{URL: "https://example.com/bad"}

	measured := client.Enrich(context.Background(), []*models.PageResult{good, bad}, "mobile", 2)

	if measured != 1 {
		t.Errorf("measured = %d, want 1", measured)
	}
	if good.WebVitals == nil || good.WebVitals.Error != "" || good.WebVitals.LCPMs == nil {
		t.Errorf("good page vitals = %+v", good.WebVitals)
	}
	if bad.WebVitals == nil || !strings.Contains(bad.WebVitals.Error, "500") {
		t.Errorf("bad page vitals = %+v, want the error recorded", bad.WebVitals)
	}
}

func TestSelectPages(t *testing.T) {
	ok := func(url string, links ...string) *models.PageResult {
		return &models.PageResult{URL: url, StatusCode: 200, InternalLinks: links,
			IndexabilityStatus: models.IndexabilityIndexable}
	}
	home := ok("https://example.com", "https://example.com/pricing", "https://example.com/blog/")
	pricing := ok("https://example.com/pricing", "https://example.com", "https://example.com/blog")
	blog := ok("https://example.com/blog", "https://example.com", "https://example.com/pricing")
	post := ok("https://example.com/blog/post", "https://example.com", "https://example.com/pricing")
	gone := &models.PageResult{URL: "https://example.com/gone", StatusCode: 404}
	noindex := ok("https://example.com/thanks")
	noindex.IndexabilityStatus = models.IndexabilityNoindex

	selected := SelectPages([]*models.PageResult{post, gone, noindex, blog, pricing, home}, 3)

	var urls []string
	for _, p := range selected {
		urls = append(urls, p.URL)
	}
	// home and pricing have 3 inlinks each, blog has 2; ties go to the shorter URL
	want := []string{"https://example.com", "https://example.com/pricing", "https://example.com/blog"}
	if strings.Join(urls, ",") != strings.Join(want, ",") {
		t.Errorf("selected = %v, want %v", urls, want)
	}
	if got := SelectPages([]*models.PageResult{home}, 0); got != nil {
		t.Errorf("limit 0 selected %v, want nothing", got)
	}
}
