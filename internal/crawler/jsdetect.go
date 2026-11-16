package crawler

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// clientRenderedMaxWords is the main-text word count at or below which a page
// with an app mount point is judged to be an empty shell.
const clientRenderedMaxWords = 50

// frameworkSignatures maps a framework to selectors that identify it in the
// HTML the server sends. Order matters: meta-frameworks come before the
// libraries they are built on.
var frameworkSignatures = []struct {
	name     string
	selector string
}{
	{"Next.js", "#__next, script#__NEXT_DATA__, script[src*='/_next/']"},
	{"Nuxt", "#__nuxt, script[src*='/_nuxt/']"},
	{"Gatsby", "#___gatsby"},
	{"SvelteKit", "[data-sveltekit-preload-data], script[src*='/_app/immutable/']"},
	{"Astro", "astro-island"},
	{"Angular", "[ng-version], app-root"},
	{"Vue", "[data-v-app], [data-server-rendered]"},
	{"React", "[data-reactroot], #root"},
}

// mountSelector matches the element a client-side app renders into.
const mountSelector = "#root, #app, #__next, #__nuxt, #___gatsby, app-root, [data-reactroot]"

// detectJS judges, from the HTML the server sent, whether a page is built with
// a JavaScript framework and whether its content exists without running
// JavaScript. It does not execute any script, so it can say a page looks like
// an empty app shell but not what the page contains once rendered.
func detectJS(doc *goquery.Document, content PageContent) *models.JSInfo {
	info := &models.JSInfo{}

	for _, signature := range frameworkSignatures {
		if doc.Find(signature.selector).Length() > 0 {
			info.Framework = signature.name
			info.Signals = append(info.Signals, signature.name+" markers in the HTML")
			break
		}
	}

	// Count executable scripts, inline or external. JSON blocks (structured
	// data, framework state) carry data, not code.
	scripts := doc.Find("script").FilterFunction(func(_ int, script *goquery.Selection) bool {
		return !strings.Contains(strings.ToLower(script.AttrOr("type", "")), "json")
	}).Length()
	mount := doc.Find(mountSelector).First()

	switch {
	case mount.Length() > 0 && content.WordCount <= clientRenderedMaxWords && scripts > 0:
		// A framework mount point with next to no text: the server sent a
		// shell and the content is filled in by JavaScript.
		mountText := strings.Fields(mount.Text())
		if len(mountText) <= clientRenderedMaxWords {
			info.ClientRendered = true
			info.Signals = append(info.Signals,
				fmt.Sprintf("app mount element holds %d words", len(mountText)),
				fmt.Sprintf("%d words of main text in the HTML", content.WordCount),
				fmt.Sprintf("%d scripts", scripts))
		}
	case content.WordCount == 0 && scripts >= 3:
		// No recognisable mount point, but an empty body and several scripts
		info.ClientRendered = true
		info.Signals = append(info.Signals,
			"no text in the HTML body",
			fmt.Sprintf("%d scripts", scripts))
	}

	if info.Framework == "" && !info.ClientRendered {
		return nil
	}
	return info
}
