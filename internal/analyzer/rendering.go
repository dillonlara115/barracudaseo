package analyzer

import (
	"fmt"
	"strings"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Core Web Vitals and JavaScript rendering issue types.
const (
	IssueSlowLCP IssueType = "slow_lcp"
	IssueSlowINP IssueType = "slow_inp"
	IssueHighCLS IssueType = "high_cls"

	IssueJSDependentContent     IssueType = "js_dependent_content"
	IssueContentRequiresJS      IssueType = "content_requires_javascript"
	IssueLinksRequireJS         IssueType = "links_require_javascript"
	IssueRenderedCanonicalDiff  IssueType = "rendered_canonical_mismatch"
	IssueRenderedRobotsMismatch IssueType = "rendered_robots_mismatch"
	IssueRenderedTitleMismatch  IssueType = "rendered_title_mismatch"
)

// Core Web Vitals thresholds, at the 75th percentile of page loads. Up to the
// "good" value passes; above the "poor" value fails; between is "needs
// improvement".
const (
	lcpGoodMs = 2500
	lcpPoorMs = 4000
	inpGoodMs = 200
	inpPoorMs = 500
	clsGood   = 0.1
	clsPoor   = 0.25
)

const (
	// renderedContentMinGain and renderedContentRatio decide when rendering
	// added enough text to say the content depends on JavaScript: at least
	// this many more words, and at least this multiple of the raw count.
	renderedContentMinGain = 100
	renderedContentRatio   = 2
)

// vitalSeverity grades a measurement: "" when good, "warning" when it needs
// improvement, "error" when poor.
func vitalSeverity(value, good, poor float64) string {
	switch {
	case value > poor:
		return "error"
	case value > good:
		return "warning"
	default:
		return ""
	}
}

// vitalMeasurement is one Core Web Vital reading ready to report.
type vitalMeasurement struct {
	issueType  IssueType
	name       string
	value      float64
	good, poor float64
	format     func(float64) string
	advice     string
}

func formatMs(value float64) string    { return fmt.Sprintf("%.1f s", value/1000) }
func formatMsRaw(value float64) string { return fmt.Sprintf("%.0f ms", value) }
func formatCLS(value float64) string   { return fmt.Sprintf("%.2f", value) }

// measurementsFor lists the readings available for a page. Field data (real
// users) is preferred. Lab data is used only when Google has no field data,
// and never for INP, which a lab test cannot measure.
func measurementsFor(v *models.WebVitals) (readings []vitalMeasurement, source string) {
	lcp := vitalMeasurement{issueType: IssueSlowLCP, name: "Largest Contentful Paint", good: lcpGoodMs, poor: lcpPoorMs, format: formatMs,
		advice: "Speed up the largest above-the-fold element: optimise and preload the hero image, cut server response time, and remove render-blocking resources"}
	inp := vitalMeasurement{issueType: IssueSlowINP, name: "Interaction to Next Paint", good: inpGoodMs, poor: inpPoorMs, format: formatMsRaw,
		advice: "Reduce JavaScript work on interaction: break up long tasks, defer non-critical scripts, and trim third-party code"}
	cls := vitalMeasurement{issueType: IssueHighCLS, name: "Cumulative Layout Shift", good: clsGood, poor: clsPoor, format: formatCLS,
		advice: "Reserve space for images, ads and embeds with width and height, and avoid inserting content above what is already on screen"}

	if v.FieldSource != "" {
		if v.LCPMs != nil {
			lcp.value = *v.LCPMs
			readings = append(readings, lcp)
		}
		if v.INPMs != nil {
			inp.value = *v.INPMs
			readings = append(readings, inp)
		}
		if v.CLS != nil {
			cls.value = *v.CLS
			readings = append(readings, cls)
		}
		return readings, v.FieldSource
	}

	if v.LabLCPMs != nil {
		lcp.value = *v.LabLCPMs
		readings = append(readings, lcp)
	}
	if v.LabCLS != nil {
		cls.value = *v.LabCLS
		readings = append(readings, cls)
	}
	return readings, "lab"
}

// AnalyzeWebVitals reports Core Web Vitals that miss Google's thresholds, for
// the pages that were measured.
func AnalyzeWebVitals(results []*models.PageResult) []Issue {
	var issues []Issue

	// Site-wide (origin) field data is the same for every page that falls back
	// to it, so it is reported once, on the first such page.
	originReported := false

	for _, p := range results {
		v := p.WebVitals
		if v == nil || v.Error != "" {
			continue
		}

		readings, source := measurementsFor(v)
		if source == "origin" {
			if originReported {
				continue
			}
			originReported = true
		}

		var basis string
		switch source {
		case "page":
			basis = "real-user data for this page"
		case "origin":
			basis = "real-user data for the whole site; Google has too little for this page alone"
		default:
			basis = "a lab test; Google has no real-user data for this page"
		}

		for _, reading := range readings {
			severity := vitalSeverity(reading.value, reading.good, reading.poor)
			if severity == "" {
				continue
			}
			// A single simulated load is weaker evidence than real users
			if source == "lab" {
				severity = "warning"
			}
			issues = append(issues, Issue{
				Type:     reading.issueType,
				Severity: severity,
				URL:      p.URL,
				Message: fmt.Sprintf("%s is %s on %s (good is %s or less), from %s",
					reading.name, reading.format(reading.value), v.Strategy, reading.format(reading.good), basis),
				Value:          reading.format(reading.value),
				Recommendation: reading.advice,
			})
		}
	}

	return issues
}

// hasNoindex reports whether a robots directive contains noindex.
func hasNoindex(directive string) bool {
	return strings.Contains(strings.ToLower(directive), "noindex")
}

// AnalyzeRendering reports pages whose content depends on JavaScript. Pages
// that were rendered in a browser are compared with their raw HTML; for the
// rest, the judgement rests on what the raw HTML looks like.
func AnalyzeRendering(results []*models.PageResult) []Issue {
	var issues []Issue

	for _, p := range results {
		if !isHTMLPage(p) {
			continue
		}

		r := p.Rendered
		if r == nil || r.Error != "" {
			// Not rendered: fall back to the raw-HTML judgement
			if p.JS != nil && p.JS.ClientRendered {
				framework := "a JavaScript"
				if p.JS.Framework != "" {
					framework = "a " + p.JS.Framework
				}
				issues = append(issues, Issue{
					Type:           IssueJSDependentContent,
					Severity:       "warning",
					URL:            p.URL,
					Message:        fmt.Sprintf("Page looks like %s app shell: the HTML has %d words of content before JavaScript runs", framework, p.WordCount),
					Value:          strings.Join(p.JS.Signals, "; "),
					Recommendation: "Render the content on the server (SSR or static generation) so search engines and AI crawlers that do not run JavaScript can read it",
				})
			}
			continue
		}

		gain := r.RenderedWordCount - r.RawWordCount
		if gain >= renderedContentMinGain && r.RenderedWordCount >= r.RawWordCount*renderedContentRatio {
			issues = append(issues, Issue{
				Type:           IssueContentRequiresJS,
				Severity:       "warning",
				URL:            p.URL,
				Message:        fmt.Sprintf("Most content appears only after JavaScript runs: %d words in the HTML, %d once rendered", r.RawWordCount, r.RenderedWordCount),
				Value:          fmt.Sprintf("%d -> %d words", r.RawWordCount, r.RenderedWordCount),
				Recommendation: "Render the content on the server (SSR or static generation); Google renders JavaScript late and many AI crawlers do not render it at all",
			})
		}

		if len(r.LinksOnlyRendered) > 0 {
			issues = append(issues, Issue{
				Type:           IssueLinksRequireJS,
				Severity:       "warning",
				URL:            p.URL,
				Message:        fmt.Sprintf("%d internal link(s) exist only after JavaScript runs, e.g. %s", len(r.LinksOnlyRendered), r.LinksOnlyRendered[0]),
				Value:          fmt.Sprintf("%d links", len(r.LinksOnlyRendered)),
				Recommendation: "Put navigation and internal links in the server-rendered HTML as real <a href> elements so crawlers can follow them",
			})
		}

		rawCanonical, renderedCanonical := urlKey(strings.TrimSpace(r.RawCanonical)), urlKey(strings.TrimSpace(r.RenderedCanonical))
		if rawCanonical != renderedCanonical {
			issues = append(issues, Issue{
				Type:           IssueRenderedCanonicalDiff,
				Severity:       "error",
				URL:            p.URL,
				Message:        fmt.Sprintf("Canonical changes when JavaScript runs: %q in the HTML, %q once rendered", r.RawCanonical, r.RenderedCanonical),
				Value:          r.RenderedCanonical,
				Recommendation: "Serve the same canonical in the initial HTML as after rendering; when they differ, Google may use either one",
			})
		}

		if hasNoindex(r.RawMetaRobots) != hasNoindex(r.RenderedMetaRobots) {
			issues = append(issues, Issue{
				Type:           IssueRenderedRobotsMismatch,
				Severity:       "error",
				URL:            p.URL,
				Message:        fmt.Sprintf("noindex changes when JavaScript runs: meta robots is %q in the HTML, %q once rendered", r.RawMetaRobots, r.RenderedMetaRobots),
				Value:          r.RenderedMetaRobots,
				Recommendation: "Serve the intended robots directive in the initial HTML; a noindex in the raw HTML can be honoured even if JavaScript removes it",
			})
		}

		if strings.TrimSpace(r.RawTitle) != strings.TrimSpace(r.RenderedTitle) {
			issues = append(issues, Issue{
				Type:           IssueRenderedTitleMismatch,
				Severity:       "info",
				URL:            p.URL,
				Message:        fmt.Sprintf("Title changes when JavaScript runs: %q in the HTML, %q once rendered", r.RawTitle, r.RenderedTitle),
				Value:          r.RenderedTitle,
				Recommendation: "Set the final title in the server-rendered HTML so every crawler sees the same one",
			})
		}
	}

	return issues
}
