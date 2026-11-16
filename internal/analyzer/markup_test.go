package analyzer

import (
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

func TestAnalyzeStructuredData(t *testing.T) {
	complete := page("https://example.com/complete")
	complete.SchemaItems = []models.SchemaItem{
		{Format: "json-ld", Type: "Product", Properties: []string{"name", "offers"}},
		{Format: "json-ld", Type: "Organization", Properties: []string{"name"}},
	}

	incomplete := page("https://example.com/incomplete")
	incomplete.SchemaItems = []models.SchemaItem{
		{Format: "json-ld", Type: "Product", Properties: []string{"description"}},
		{Format: "json-ld", Type: "LocalBusiness", Properties: []string{"name"}},
	}

	// Declared twice; one declaration is complete, so nothing to report
	twice := page("https://example.com/twice")
	twice.SchemaItems = []models.SchemaItem{
		{Format: "microdata", Type: "Product", Properties: []string{"name"}},
		{Format: "json-ld", Type: "Product", Properties: []string{"name", "aggregateRating"}},
	}

	broken := page("https://example.com/broken")
	broken.SchemaErrors = []string{"JSON-LD block 1 is not valid JSON: unexpected end", "JSON-LD block 2 has no @type"}

	issues, types := AnalyzeStructuredData([]*models.PageResult{complete, incomplete, twice, broken})

	if types["Product"] != 3 || types["Organization"] != 1 || types["LocalBusiness"] != 1 {
		t.Errorf("type counts = %v, want Product 3, Organization 1, LocalBusiness 1", types)
	}

	missing := issuesOfType(issues, IssueStructuredDataMissingFields)
	if len(missing) != 2 {
		t.Fatalf("structured_data_missing_fields = %+v, want 2 (both on the incomplete page)", missing)
	}
	if missing[0].Value != "LocalBusiness" || !strings.Contains(missing[0].Message, "address") {
		t.Errorf("first issue = %+v, want LocalBusiness missing address", missing[0])
	}
	if missing[1].Value != "Product" || !strings.Contains(missing[1].Message, "name") ||
		!strings.Contains(missing[1].Message, "one of offers, review, aggregateRating") {
		t.Errorf("second issue = %+v, want Product missing name and an offer/review/rating", missing[1])
	}
	if missing[0].Fingerprint() == missing[1].Fingerprint() {
		t.Error("two schema types on one page should have different fingerprints")
	}

	invalid := issuesOfType(issues, IssueInvalidStructuredData)
	if len(invalid) != 1 || invalid[0].URL != broken.URL || !strings.Contains(invalid[0].Message, "1 more") {
		t.Errorf("invalid_structured_data = %+v, want one on the broken page noting a second problem", invalid)
	}
}

func TestValidHreflangCode(t *testing.T) {
	valid := []string{"en", "en-GB", "EN-gb", "zh-Hant", "zh-Hant-TW", "x-default"}
	for _, code := range valid {
		if !validHreflangCode(code) {
			t.Errorf("validHreflangCode(%q) = false, want true", code)
		}
	}
	invalid := []string{"en_US", "english", "en-UK", "e", "en-", "en-GBR", ""}
	for _, code := range invalid {
		if validHreflangCode(code) {
			t.Errorf("validHreflangCode(%q) = true, want false", code)
		}
	}
}

func TestAnalyzeHreflang(t *testing.T) {
	const (
		en = "https://example.com/en"
		de = "https://example.com/de"
		fr = "https://example.com/fr"
		es = "https://example.com/es"
	)
	set := func(links ...string) []models.HreflangLink {
		var out []models.HreflangLink
		for i := 0; i < len(links); i += 2 {
			out = append(out, models.HreflangLink{Lang: links[i], URL: links[i+1]})
		}
		return out
	}

	// en and de are a correct reciprocal pair with x-default
	enPage := page(en)
	enPage.Hreflang = set("en", en, "de", de, "fr", fr, "es", es, "x-default", en, "en_US", en)
	dePage := page(de)
	dePage.Hreflang = set("en", en, "de", de+"/", "x-default", en)

	// fr does not link back to en, and leaves itself and x-default out
	frPage := page(fr)
	frPage.Hreflang = set("de", de)

	// es is noindexed
	esPage := page(es)
	esPage.IndexabilityStatus = models.IndexabilityNoindex

	issues := AnalyzeHreflang([]*models.PageResult{enPage, dePage, frPage, esPage})

	invalid := issuesOfType(issues, IssueHreflangInvalidLanguage)
	if len(invalid) != 1 || invalid[0].URL != en || invalid[0].Value != "en_US" {
		t.Errorf("hreflang_invalid_language = %+v, want one for en_US on the en page", invalid)
	}

	returnLinks := issuesOfType(issues, IssueHreflangMissingReturnLink)
	if len(returnLinks) != 2 {
		t.Fatalf("hreflang_missing_return_link = %+v, want 2 (en->fr, fr->de)", returnLinks)
	}
	if returnLinks[0].URL != en || returnLinks[0].Value != fr {
		t.Errorf("first missing return link = %+v, want en pointing at fr", returnLinks[0])
	}
	if returnLinks[1].URL != fr || returnLinks[1].Value != de {
		t.Errorf("second missing return link = %+v, want fr pointing at de", returnLinks[1])
	}

	nonIndexable := issuesOfType(issues, IssueHreflangToNonIndexable)
	if len(nonIndexable) != 1 || nonIndexable[0].Value != es || !strings.Contains(nonIndexable[0].Message, "noindex") {
		t.Errorf("hreflang_to_non_indexable = %+v, want one for the noindexed es page", nonIndexable)
	}

	missingSelf := issuesOfType(issues, IssueHreflangMissingSelf)
	if len(missingSelf) != 1 || missingSelf[0].URL != fr {
		t.Errorf("hreflang_missing_self_reference = %+v, want only the fr page", missingSelf)
	}
	missingDefault := issuesOfType(issues, IssueHreflangMissingXDefault)
	if len(missingDefault) != 1 || missingDefault[0].URL != fr {
		t.Errorf("hreflang_missing_x_default = %+v, want only the fr page", missingDefault)
	}
}

func TestAnalyzeSecurity(t *testing.T) {
	allHeaders := map[string]string{
		"strict-transport-security": "max-age=31536000",
		"content-security-policy":   "default-src 'self'",
		"x-content-type-options":    "nosniff",
		"x-frame-options":           "DENY",
		"referrer-policy":           "strict-origin-when-cross-origin",
	}

	home := page("https://example.com")
	home.SecurityHeaders = map[string]string{"x-content-type-options": "nosniff"}

	deep := page("https://example.com/a/deep/page")
	deep.SecurityHeaders = allHeaders
	deep.MixedContent = []string{"http://a/1.jpg", "http://a/2.jpg", "http://a/3.jpg", "http://a/4.jpg"}

	insecure := page("http://example.com/legacy")
	insecure.SecurityHeaders = map[string]string{}

	// Headers were not recorded for this page: it must not count as missing
	imported := page("https://example.com/imported")

	issues := AnalyzeSecurity([]*models.PageResult{deep, home, insecure, imported})

	mixed := issuesOfType(issues, IssueMixedContent)
	if len(mixed) != 1 || mixed[0].URL != deep.URL || !strings.Contains(mixed[0].Message, "4 resource(s)") ||
		!strings.Contains(mixed[0].Message, "+1 more") {
		t.Errorf("mixed_content = %+v, want one on the deep page listing 3 of 4", mixed)
	}

	http := issuesOfType(issues, IssueInsecurePage)
	if len(http) != 1 || http[0].URL != insecure.URL {
		t.Errorf("insecure_page = %+v, want only the http page", http)
	}

	headers := issuesOfType(issues, IssueMissingSecurityHeader)
	if len(headers) != 5 {
		t.Fatalf("missing_security_header count = %d, want 5: %+v", len(headers), headers)
	}
	byHeader := make(map[string]Issue)
	fingerprints := make(map[string]bool)
	for _, issue := range headers {
		byHeader[issue.Value] = issue
		fingerprints[issue.Fingerprint()] = true
	}
	if len(fingerprints) != 5 {
		t.Error("each missing header should have its own fingerprint")
	}

	// HSTS only applies to https pages: home is missing it, deep has it
	hsts := byHeader["Strict-Transport-Security"]
	if hsts.URL != home.URL || !strings.Contains(hsts.Message, "1 of 2 pages") {
		t.Errorf("HSTS issue = %+v, want it on the homepage, 1 of 2 https pages", hsts)
	}
	// nosniff: home and deep have it, the http page does not
	nosniff := byHeader["X-Content-Type-Options"]
	if nosniff.URL != insecure.URL || !strings.Contains(nosniff.Message, "1 of 3 pages") {
		t.Errorf("X-Content-Type-Options issue = %+v, want 1 of 3 pages with recorded headers", nosniff)
	}
}

func TestAnalyzeSecurity_NoHeaderDataReportsNothing(t *testing.T) {
	issues := AnalyzeSecurity([]*models.PageResult{page("https://example.com"), page("https://example.com/a")})
	if got := issuesOfType(issues, IssueMissingSecurityHeader); len(got) != 0 {
		t.Errorf("missing_security_header = %+v, want none when headers were never recorded", got)
	}
}
