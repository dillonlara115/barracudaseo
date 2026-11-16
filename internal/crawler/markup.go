package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// securityHeaderNames are the response headers recorded for each page.
var securityHeaderNames = []string{
	"Strict-Transport-Security",
	"Content-Security-Policy",
	"X-Frame-Options",
	"X-Content-Type-Options",
	"Referrer-Policy",
}

// extractSecurityHeaders returns the security headers present on a response,
// keyed by lowercase name.
func extractSecurityHeaders(header http.Header) map[string]string {
	found := make(map[string]string)
	for _, name := range securityHeaderNames {
		if value := strings.TrimSpace(header.Get(name)); value != "" {
			found[strings.ToLower(name)] = value
		}
	}
	return found
}

// extractStructuredData reads a page's JSON-LD blocks and microdata items. It
// returns the entities found and a description of each JSON-LD block that
// could not be used.
func extractStructuredData(doc *goquery.Document) ([]models.SchemaItem, []string) {
	var items []models.SchemaItem
	var errors []string

	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		raw := strings.TrimSpace(s.Text())
		if raw == "" {
			errors = append(errors, fmt.Sprintf("JSON-LD block %d is empty", i+1))
			return
		}

		var parsed interface{}
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			errors = append(errors, fmt.Sprintf("JSON-LD block %d is not valid JSON: %v", i+1, err))
			return
		}

		found := collectJSONLDItems(parsed)
		if len(found) == 0 {
			errors = append(errors, fmt.Sprintf("JSON-LD block %d has no @type", i+1))
			return
		}
		items = append(items, found...)
	})

	// Microdata: top-level items only (nested items belong to their parent)
	doc.Find("[itemscope][itemtype]").Each(func(_ int, s *goquery.Selection) {
		if s.ParentsFiltered("[itemscope]").Length() > 0 {
			return
		}
		itemType := schemaTypeName(s.AttrOr("itemtype", ""))
		if itemType == "" {
			return
		}
		seen := make(map[string]bool)
		var properties []string
		s.Find("[itemprop]").Each(func(_ int, prop *goquery.Selection) {
			for _, name := range strings.Fields(prop.AttrOr("itemprop", "")) {
				if !seen[name] {
					seen[name] = true
					properties = append(properties, name)
				}
			}
		})
		sort.Strings(properties)
		items = append(items, models.SchemaItem{Format: "microdata", Type: itemType, Properties: properties})
	})

	return items, errors
}

// collectJSONLDItems returns the entities in a parsed JSON-LD value: the
// object itself, each element of a top-level array, and each node of @graph.
// Entities nested inside a property (a Product's offers, say) are part of
// their parent and are not listed separately.
func collectJSONLDItems(value interface{}) []models.SchemaItem {
	var items []models.SchemaItem

	switch v := value.(type) {
	case []interface{}:
		for _, element := range v {
			items = append(items, collectJSONLDItems(element)...)
		}
	case map[string]interface{}:
		if graph, ok := v["@graph"]; ok {
			items = append(items, collectJSONLDItems(graph)...)
		}

		var properties []string
		for key := range v {
			if !strings.HasPrefix(key, "@") {
				properties = append(properties, key)
			}
		}
		sort.Strings(properties)

		for _, typeName := range jsonLDTypes(v["@type"]) {
			items = append(items, models.SchemaItem{Format: "json-ld", Type: typeName, Properties: properties})
		}
	}

	return items
}

// jsonLDTypes returns the type names of an @type value, which may be a single
// string or a list.
func jsonLDTypes(value interface{}) []string {
	var types []string
	switch v := value.(type) {
	case string:
		if name := schemaTypeName(v); name != "" {
			types = append(types, name)
		}
	case []interface{}:
		for _, element := range v {
			if s, ok := element.(string); ok {
				if name := schemaTypeName(s); name != "" {
					types = append(types, name)
				}
			}
		}
	}
	return types
}

// schemaTypeName strips a schema.org URL prefix: "https://schema.org/Product"
// and "schema:Product" both become "Product".
func schemaTypeName(value string) string {
	value = strings.TrimSpace(value)
	if i := strings.LastIndexAny(value, "/#:"); i >= 0 {
		value = value[i+1:]
	}
	return value
}

// extractHreflang returns the page's hreflang annotations from <link> tags,
// with URLs resolved against the page URL.
func extractHreflang(doc *goquery.Document, pageURL string) []models.HreflangLink {
	var links []models.HreflangLink
	doc.Find("link[hreflang][href]").Each(func(_ int, s *goquery.Selection) {
		if !strings.Contains(strings.ToLower(s.AttrOr("rel", "")), "alternate") {
			return
		}
		lang := strings.TrimSpace(s.AttrOr("hreflang", ""))
		href := strings.TrimSpace(s.AttrOr("href", ""))
		if lang == "" || href == "" {
			return
		}
		resolved, err := utils.ResolveURL(pageURL, href)
		if err != nil {
			resolved = href
		}
		links = append(links, models.HreflangLink{Lang: lang, URL: resolved})
	})
	return links
}

// mixedContentSelectors maps elements that load a resource to the attribute
// holding its URL. Plain links (<a href>) are not included: following a link
// to an http:// page is not mixed content.
var mixedContentSelectors = []struct {
	selector  string
	attribute string
}{
	{"img[src]", "src"},
	{"script[src]", "src"},
	{"iframe[src]", "src"},
	{"video[src], audio[src], source[src], embed[src]", "src"},
	{"link[rel~='stylesheet'][href]", "href"},
	{"object[data]", "data"},
	{"form[action]", "action"},
}

// extractMixedContent returns the http:// resources an https:// page loads.
// Browsers block or warn about these, and they undermine the page's HTTPS.
func extractMixedContent(doc *goquery.Document, pageURL string) []string {
	if !strings.HasPrefix(strings.ToLower(pageURL), "https://") {
		return nil
	}

	seen := make(map[string]bool)
	var insecure []string
	for _, target := range mixedContentSelectors {
		doc.Find(target.selector).Each(func(_ int, s *goquery.Selection) {
			value := strings.TrimSpace(s.AttrOr(target.attribute, ""))
			if !strings.HasPrefix(strings.ToLower(value), "http://") || seen[value] {
				return
			}
			seen[value] = true
			insecure = append(insecure, value)
		})
	}
	return insecure
}
