package crawler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// boilerplateSelector matches elements whose text is not page content: code,
// and the navigation / header / footer chrome repeated on every page. Leaving
// it in would make unrelated pages look alike.
const boilerplateSelector = "script, style, noscript, template, svg, iframe, nav, header, footer, aside, form"

// shingleSize is the number of consecutive words hashed together for the
// similarity hash. Word order matters at this size, so two pages sharing
// vocabulary but not sentences do not look alike.
const shingleSize = 3

// PageContent summarises a page's main text for content analysis.
type PageContent struct {
	WordCount int
	// Hash identifies the exact (case- and whitespace-insensitive) main text.
	Hash string
	// SimHash is a 64-bit similarity hash, hex encoded: near-identical texts
	// differ in only a few bits.
	SimHash string
}

// ExtractContent returns the word count and content hashes of a page's main
// text. It prefers the <main> / <article> region and falls back to <body>,
// dropping boilerplate either way.
func ExtractContent(doc *goquery.Document) PageContent {
	root := doc.Find("main, article, [role='main']").First()
	if root.Length() == 0 {
		root = doc.Find("body").First()
	}
	if root.Length() == 0 {
		return PageContent{}
	}

	root = root.Clone()
	root.Find(boilerplateSelector).Remove()

	var text strings.Builder
	collectText(root, &text)

	words := tokenize(text.String())
	if len(words) == 0 {
		return PageContent{}
	}

	sum := sha256.Sum256([]byte(strings.Join(words, " ")))
	return PageContent{
		WordCount: len(words),
		Hash:      hex.EncodeToString(sum[:16]),
		SimHash:   fmt.Sprintf("%016x", simHash(words)),
	}
}

// collectText appends the text of a selection, separating nodes with spaces.
// goquery's Text() joins adjacent elements with nothing between them, which
// would glue the last word of one paragraph to the first word of the next.
func collectText(sel *goquery.Selection, out *strings.Builder) {
	sel.Contents().Each(func(_ int, node *goquery.Selection) {
		if goquery.NodeName(node) == "#text" {
			out.WriteString(node.Text())
			out.WriteByte(' ')
			return
		}
		collectText(node, out)
	})
}

// tokenize splits text into lowercase words of letters and digits.
func tokenize(text string) []string {
	return strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// simHash computes a 64-bit SimHash over word shingles. Each shingle votes on
// every bit, so texts that share most shingles end up with nearly equal hashes.
func simHash(words []string) uint64 {
	var votes [64]int

	vote := func(shingle string) {
		h := fnv.New64a()
		_, _ = h.Write([]byte(shingle))
		sum := h.Sum64()
		for bit := 0; bit < 64; bit++ {
			if sum&(1<<uint(bit)) != 0 {
				votes[bit]++
			} else {
				votes[bit]--
			}
		}
	}

	if len(words) < shingleSize {
		vote(strings.Join(words, " "))
	} else {
		for i := 0; i+shingleSize <= len(words); i++ {
			vote(strings.Join(words[i:i+shingleSize], " "))
		}
	}

	var hash uint64
	for bit := 0; bit < 64; bit++ {
		if votes[bit] > 0 {
			hash |= 1 << uint(bit)
		}
	}
	return hash
}
