package crawler

import (
	"math/bits"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func contentOf(t *testing.T, html string) PageContent {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	return ExtractContent(doc)
}

func distance(t *testing.T, a, b PageContent) int {
	t.Helper()
	x, errA := strconv.ParseUint(a.SimHash, 16, 64)
	y, errB := strconv.ParseUint(b.SimHash, 16, 64)
	if errA != nil || errB != nil {
		t.Fatalf("invalid simhash: %q, %q", a.SimHash, b.SimHash)
	}
	return bits.OnesCount64(x ^ y)
}

func TestExtractContent_WordCountExcludesBoilerplate(t *testing.T) {
	content := contentOf(t, `<html><body>
		<header>Site name and tagline</header>
		<nav><a href="/">Home</a><a href="/about">About</a></nav>
		<p>One two three.</p><p>Four five</p>
		<script>var ignored = "not content";</script>
		<style>.ignored { color: red }</style>
		<footer>Copyright notice here</footer>
	</body></html>`)

	// "three." and "Four" sit in adjacent elements and must not be glued together
	if content.WordCount != 5 {
		t.Errorf("WordCount = %d, want 5 (boilerplate excluded, paragraphs kept apart)", content.WordCount)
	}
	if content.Hash == "" || len(content.SimHash) != 16 {
		t.Errorf("hashes not set: %+v", content)
	}
}

func TestExtractContent_PrefersMainRegion(t *testing.T) {
	content := contentOf(t, `<html><body>
		<div>Sidebar promo text that is not the article</div>
		<main><h1>Title here</h1><p>Body text</p></main>
	</body></html>`)

	if content.WordCount != 4 {
		t.Errorf("WordCount = %d, want 4 (only the <main> region)", content.WordCount)
	}
}

func TestExtractContent_HashIgnoresChromeCaseAndWhitespace(t *testing.T) {
	a := contentOf(t, `<html><body><nav>Menu A</nav><main><p>Same   Article text</p></main></body></html>`)
	b := contentOf(t, `<html><body><nav>Different menu</nav><main><p>same article
		TEXT</p></main><footer>Other footer</footer></body></html>`)
	c := contentOf(t, `<html><body><main><p>Another article entirely</p></main></body></html>`)

	if a.Hash != b.Hash {
		t.Error("pages with the same main text should share a content hash")
	}
	if a.Hash == c.Hash {
		t.Error("pages with different main text should have different content hashes")
	}
}

func TestExtractContent_EmptyPage(t *testing.T) {
	content := contentOf(t, `<html><body><nav>Only navigation</nav></body></html>`)
	if content.WordCount != 0 || content.Hash != "" || content.SimHash != "" {
		t.Errorf("content = %+v, want zero value for a page with no main text", content)
	}
}

func TestSimHash_NearDuplicatesAreClose(t *testing.T) {
	// 300 distinct words per text, so no shingle repeats (repeated sentences
	// would give a few shingles outsized weight and make the hash unstable).
	text := func(prefix string) string {
		words := make([]string, 300)
		for i := range words {
			words[i] = prefix + strconv.Itoa(i)
		}
		return strings.Join(words, " ")
	}

	original := contentOf(t, "<main><p>"+text("plumbing")+" Call us in Austin today.</p></main>")
	variant := contentOf(t, "<main><p>"+text("plumbing")+" Call us in Dallas today.</p></main>")
	unrelated := contentOf(t, "<main><p>"+text("museum")+" Tickets are on sale now.</p></main>")

	if original.Hash == variant.Hash {
		t.Fatal("variant differs by one word and should not share the exact hash")
	}
	if d := distance(t, original, variant); d > 3 {
		t.Errorf("near-duplicate distance = %d bits, want <= 3", d)
	}
	if d := distance(t, original, unrelated); d <= 3 {
		t.Errorf("unrelated distance = %d bits, want > 3", d)
	}
}
