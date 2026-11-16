package analyzer

import "testing"

func TestIssueFingerprint(t *testing.T) {
	const page = "https://example.com/page"

	same := []struct {
		name string
		a, b Issue
	}{
		{
			"trailing slash on page URL",
			Issue{Type: IssueMissingH1, URL: page},
			Issue{Type: IssueMissingH1, URL: page + "/"},
		},
		{
			"per-page issue ignores a changed value",
			Issue{Type: IssueLongTitle, URL: page, Value: "An old long title"},
			Issue{Type: IssueLongTitle, URL: page, Value: "A different long title"},
		},
		{
			"large image ignores measured size",
			Issue{Type: IssueLargeImage, URL: page, Value: "https://example.com/a.jpg (250 KB)"},
			Issue{Type: IssueLargeImage, URL: page, Value: "https://example.com/a.jpg (180 KB)"},
		},
		{
			"link to redirect ignores final destination",
			Issue{Type: IssueLinkToRedirect, URL: page, Value: "https://example.com/old -> https://example.com/new"},
			Issue{Type: IssueLinkToRedirect, URL: page, Value: "https://example.com/old -> https://example.com/newer"},
		},
	}
	for _, tt := range same {
		if tt.a.Fingerprint() != tt.b.Fingerprint() {
			t.Errorf("%s: fingerprints differ, want equal", tt.name)
		}
	}

	different := []struct {
		name string
		a, b Issue
	}{
		{
			"different type",
			Issue{Type: IssueMissingH1, URL: page},
			Issue{Type: IssueMissingTitle, URL: page},
		},
		{
			"different page",
			Issue{Type: IssueMissingH1, URL: page},
			Issue{Type: IssueMissingH1, URL: page + "-2"},
		},
		{
			"different image on the same page",
			Issue{Type: IssueMissingImageAlt, URL: page, Value: "https://example.com/a.jpg"},
			Issue{Type: IssueMissingImageAlt, URL: page, Value: "https://example.com/b.jpg"},
		},
		{
			"different broken link target on the same page",
			Issue{Type: IssueBrokenInternalLink, URL: page, Value: "https://example.com/gone"},
			Issue{Type: IssueBrokenInternalLink, URL: page, Value: "https://example.com/also-gone"},
		},
	}
	for _, tt := range different {
		if tt.a.Fingerprint() == tt.b.Fingerprint() {
			t.Errorf("%s: fingerprints equal, want different", tt.name)
		}
	}
}
