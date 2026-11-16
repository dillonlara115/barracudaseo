package analyzer

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// issueSubject returns what, besides the page, an issue is about. Most issue
// types occur at most once per page, so the page URL alone identifies them and
// the subject is empty. Types that can occur several times on one page (one per
// image or per link) are told apart by the asset or link target.
//
// The subject must stay the same from crawl to crawl while the problem
// persists, so volatile detail such as an image's measured size is stripped.
func issueSubject(issueType IssueType, value string) string {
	value = strings.TrimSpace(value)

	switch issueType {
	case IssueMissingImageAlt, IssueBrokenImage, IssueBrokenInternalLink:
		// Value is the image URL / link target
		return urlKey(value)
	case IssueLargeImage:
		// Value is "<image url> (123 KB)"
		if i := strings.LastIndex(value, " ("); i >= 0 {
			value = value[:i]
		}
		return urlKey(value)
	case IssueLinkToRedirect:
		// Value is "<link target> -> <final destination>"
		if i := strings.Index(value, " -> "); i >= 0 {
			value = value[:i]
		}
		return urlKey(value)
	default:
		return ""
	}
}

// IssueFingerprint returns a stable identity for an issue: the same problem on
// the same page yields the same fingerprint in every crawl. It is what makes
// crawl-to-crawl comparison and carrying an issue's status forward possible.
func IssueFingerprint(issueType IssueType, pageURL, value string) string {
	key := string(issueType) + "\n" + urlKey(strings.TrimSpace(pageURL)) + "\n" + issueSubject(issueType, value)
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:16])
}

// Fingerprint returns the issue's stable identity. See IssueFingerprint.
func (i Issue) Fingerprint() string {
	return IssueFingerprint(i.Type, i.URL, i.Value)
}
