package api

import "testing"

func TestRenderPageLimit(t *testing.T) {
	cases := map[string]int{"": 0, "5": 5, "0": 0, "-3": 0, "many": 0}
	for value, want := range cases {
		t.Setenv("RENDER_MAX_PAGES", value)
		if got := renderPageLimit(); got != want {
			t.Errorf("RENDER_MAX_PAGES=%q: limit = %d, want %d", value, got, want)
		}
	}
}

func TestPageExperienceEnabled(t *testing.T) {
	t.Setenv("PSI_API_KEY", "")
	t.Setenv("RENDER_MAX_PAGES", "")
	if pageExperienceEnabled() {
		t.Fatal("enabled with neither PSI_API_KEY nor RENDER_MAX_PAGES set")
	}

	t.Setenv("RENDER_MAX_PAGES", "5")
	if !pageExperienceEnabled() {
		t.Fatal("not enabled with RENDER_MAX_PAGES set")
	}

	t.Setenv("RENDER_MAX_PAGES", "")
	t.Setenv("PSI_API_KEY", "key")
	if !pageExperienceEnabled() {
		t.Fatal("not enabled with PSI_API_KEY set")
	}
}

func TestMergedDataKeepsExistingKeys(t *testing.T) {
	existing := map[string]interface{}{"h2": []string{"a"}, "rendered": "old"}

	got := mergedData(existing, map[string]interface{}{"rendered": "new", "web_vitals": 1})

	if got["rendered"] != "new" || got["web_vitals"] != 1 {
		t.Fatalf("new values not applied: %v", got)
	}
	if _, ok := got["h2"]; !ok {
		t.Fatalf("existing key lost: %v", got)
	}
	if existing["rendered"] != "old" {
		t.Fatal("mergedData modified its input")
	}
	if got := mergedData(nil, map[string]interface{}{"a": 1}); len(got) != 1 {
		t.Fatalf("nil existing: got %v", got)
	}
}

func TestStoredPageURLMatchesNormalizedForm(t *testing.T) {
	// Pages are stored under their normalized URL, so lookups must use it too
	if a, b := storedPageURL("https://example.com/page/"), storedPageURL("https://example.com/page"); a != b {
		t.Fatalf("trailing slash changes the stored URL: %q vs %q", a, b)
	}
}
