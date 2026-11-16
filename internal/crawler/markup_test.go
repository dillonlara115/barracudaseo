package crawler

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func docOf(t *testing.T, html string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	return doc
}

func TestExtractStructuredData_JSONLD(t *testing.T) {
	items, errors := extractStructuredData(docOf(t, `<html><head>
		<script type="application/ld+json">
		{"@context": "https://schema.org", "@type": "Product", "name": "Widget",
		 "offers": {"@type": "Offer", "price": "9.99"}}
		</script>
		<script type="application/ld+json">
		{"@context": "https://schema.org", "@graph": [
			{"@type": "Organization", "name": "Acme"},
			{"@type": ["Article", "NewsArticle"], "headline": "Hello"}
		]}
		</script>
		<script type="application/ld+json">{"@type": "https://schema.org/BreadcrumbList"}</script>
	</head><body></body></html>`))

	if len(errors) != 0 {
		t.Errorf("errors = %v, want none", errors)
	}

	var types []string
	for _, item := range items {
		types = append(types, item.Type)
	}
	want := []string{"Product", "Organization", "Article", "NewsArticle", "BreadcrumbList"}
	if !reflect.DeepEqual(types, want) {
		t.Errorf("types = %v, want %v (nested Offer is not listed separately)", types, want)
	}
	if got := items[0].Properties; !reflect.DeepEqual(got, []string{"name", "offers"}) {
		t.Errorf("Product properties = %v, want [name offers] (@-keys excluded)", got)
	}
}

func TestExtractStructuredData_Errors(t *testing.T) {
	_, errors := extractStructuredData(docOf(t, `<html><head>
		<script type="application/ld+json">{"@type": "Product", "name": "Widget",}</script>
		<script type="application/ld+json">{"name": "No type here"}</script>
		<script type="application/ld+json">   </script>
	</head></html>`))

	if len(errors) != 3 {
		t.Fatalf("errors = %v, want 3", errors)
	}
	for i, fragment := range []string{"not valid JSON", "no @type", "is empty"} {
		if !strings.Contains(errors[i], fragment) {
			t.Errorf("errors[%d] = %q, want it to mention %q", i, errors[i], fragment)
		}
	}
}

func TestExtractStructuredData_Microdata(t *testing.T) {
	items, _ := extractStructuredData(docOf(t, `<html><body>
		<div itemscope itemtype="https://schema.org/Product">
			<span itemprop="name">Widget</span>
			<div itemprop="offers" itemscope itemtype="https://schema.org/Offer">
				<span itemprop="price">9.99</span>
			</div>
		</div>
	</body></html>`))

	if len(items) != 1 || items[0].Format != "microdata" || items[0].Type != "Product" {
		t.Fatalf("items = %+v, want one microdata Product (nested Offer belongs to it)", items)
	}
	if got := items[0].Properties; !reflect.DeepEqual(got, []string{"name", "offers", "price"}) {
		t.Errorf("properties = %v", got)
	}
}

func TestExtractHreflang(t *testing.T) {
	links := extractHreflang(docOf(t, `<html><head>
		<link rel="alternate" hreflang="en" href="/en/page">
		<link rel="alternate" hreflang="de-DE" href="https://example.com/de/page">
		<link rel="alternate" type="application/rss+xml" href="/feed">
		<link rel="stylesheet" hreflang="fr" href="/style.css">
	</head></html>`), "https://example.com/en/page")

	if len(links) != 2 {
		t.Fatalf("links = %+v, want 2 (feed has no hreflang; stylesheet is not an alternate)", links)
	}
	if links[0].Lang != "en" || links[0].URL != "https://example.com/en/page" {
		t.Errorf("links[0] = %+v, want the relative href resolved", links[0])
	}
	if links[1].Lang != "de-DE" {
		t.Errorf("links[1] = %+v", links[1])
	}
}

func TestExtractMixedContent(t *testing.T) {
	html := `<html><head>
		<link rel="stylesheet" href="http://cdn.example.com/style.css">
		<script src="https://cdn.example.com/ok.js"></script>
	</head><body>
		<img src="http://example.com/a.jpg">
		<img src="http://example.com/a.jpg">
		<img src="/relative.jpg">
		<a href="http://other.example.com/">A plain link is not mixed content</a>
		<iframe src="http://embed.example.com/"></iframe>
	</body></html>`

	got := extractMixedContent(docOf(t, html), "https://example.com/page")
	want := []string{"http://example.com/a.jpg", "http://embed.example.com/", "http://cdn.example.com/style.css"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mixed content = %v, want %v", got, want)
	}

	if onHTTP := extractMixedContent(docOf(t, html), "http://example.com/page"); onHTTP != nil {
		t.Errorf("mixed content on an http page = %v, want nil", onHTTP)
	}
}

func TestExtractSecurityHeaders(t *testing.T) {
	header := http.Header{}
	header.Set("Strict-Transport-Security", "max-age=31536000")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Server", "nginx")

	got := extractSecurityHeaders(header)
	want := map[string]string{
		"strict-transport-security": "max-age=31536000",
		"x-content-type-options":    "nosniff",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("headers = %v, want %v", got, want)
	}

	if none := extractSecurityHeaders(http.Header{}); none == nil || len(none) != 0 {
		t.Errorf("no headers = %v, want an empty non-nil map", none)
	}
}
