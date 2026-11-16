package render

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/internal/crawler"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// spaPage is a page whose content, title, canonical and links are all written
// by JavaScript, the way a client-rendered app behaves.
const spaPage = `<!DOCTYPE html>
<html><head>
<title>Loading...</title>
<link rel="canonical" href="/raw-canonical">
<meta name="robots" content="noindex">
</head><body>
<div id="root"></div>
<script>
	document.title = 'Rendered Title';
	document.querySelector('link[rel=canonical]').href = '/rendered-canonical';
	document.querySelector('meta[name=robots]').content = 'index, follow';
	var words = [];
	for (var i = 0; i < 300; i++) { words.push('word' + i); }
	document.getElementById('root').innerHTML =
		'<main><h1>Rendered heading</h1><p>' + words.join(' ') + '</p>' +
		'<a href="/js-only-link">A link added by JavaScript</a></main>';
</script>
</body></html>`

// TestCompare_ClientRenderedPage drives a real headless Chrome. It is skipped
// when Chrome is not installed.
func TestCompare_ClientRenderedPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping headless Chrome test in short mode")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, spaPage)
	}))
	defer server.Close()

	renderer, err := NewRenderer("barracuda-test")
	if err != nil {
		t.Skipf("headless Chrome is not available: %v", err)
	}
	defer renderer.Close()

	// Parse the raw HTML the way the crawler does
	parser, err := crawler.NewParser(server.URL + "/")
	if err != nil {
		t.Fatalf("NewParser: %v", err)
	}
	raw, err := parser.Parse([]byte(spaPage))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	raw.URL = server.URL + "/"
	raw.StatusCode = 200

	if raw.JS == nil || !raw.JS.ClientRendered {
		t.Errorf("raw JS info = %+v, want the shell detected as client rendered", raw.JS)
	}

	comparison := renderer.Compare(raw)
	if comparison.Error != "" {
		t.Fatalf("Compare error: %s", comparison.Error)
	}

	if comparison.RawWordCount != 0 {
		t.Errorf("RawWordCount = %d, want 0", comparison.RawWordCount)
	}
	if comparison.RenderedWordCount < 300 {
		t.Errorf("RenderedWordCount = %d, want at least 300", comparison.RenderedWordCount)
	}
	if comparison.RawTitle != "Loading..." || comparison.RenderedTitle != "Rendered Title" {
		t.Errorf("titles = %q -> %q", comparison.RawTitle, comparison.RenderedTitle)
	}
	if !strings.HasSuffix(comparison.RawCanonical, "/raw-canonical") ||
		!strings.HasSuffix(comparison.RenderedCanonical, "/rendered-canonical") {
		t.Errorf("canonicals = %q -> %q", comparison.RawCanonical, comparison.RenderedCanonical)
	}
	if comparison.RawMetaRobots != "noindex" || comparison.RenderedMetaRobots != "index, follow" {
		t.Errorf("meta robots = %q -> %q", comparison.RawMetaRobots, comparison.RenderedMetaRobots)
	}
	if len(comparison.LinksOnlyRendered) != 1 || !strings.HasSuffix(comparison.LinksOnlyRendered[0], "/js-only-link") {
		t.Errorf("LinksOnlyRendered = %v, want the JavaScript-added link", comparison.LinksOnlyRendered)
	}
}

func TestSelectPages(t *testing.T) {
	page := func(url string, clientRendered bool) *models.PageResult {
		p := &models.PageResult{URL: url, StatusCode: 200}
		if clientRendered {
			p.JS = &models.JSInfo{ClientRendered: true}
		}
		return p
	}
	home := page("https://example.com", false)
	app := page("https://example.com/app", true)
	about := page("https://example.com/about", false)
	broken := page("https://example.com/broken", true)
	broken.StatusCode = 500

	selected := SelectPages([]*models.PageResult{home, app, about, broken}, []*models.PageResult{home, app, about}, 2)

	if len(selected) != 2 || selected[0] != app || selected[1] != home {
		var urls []string
		for _, p := range selected {
			urls = append(urls, p.URL)
		}
		t.Errorf("selected = %v, want the client-rendered page first, then the top priority page", urls)
	}
}
