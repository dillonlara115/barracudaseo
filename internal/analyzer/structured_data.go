package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Structured data issue types.
const (
	IssueInvalidStructuredData       IssueType = "invalid_structured_data"
	IssueStructuredDataMissingFields IssueType = "structured_data_missing_fields"
)

// schemaRequirement lists what a Schema.org type needs to be usable.
type schemaRequirement struct {
	// All of these properties must be present
	required []string
	// At least one of these must be present (ignored when empty)
	oneOf []string
}

// requiredSchemaProperties maps a Schema.org type to the properties Google
// requires before it will consider the markup for a rich result.
//
// TODO(dillon): this is a deliberately small starting set. Decide which types
// matter to your users and extend it, and check each entry against Google's
// current structured data documentation, which changes: types are added,
// retired or restricted over time. A wrong entry here produces confident,
// wrong advice, so a short accurate list beats a long approximate one.
var requiredSchemaProperties = map[string]schemaRequirement{
	"Product":        {required: []string{"name"}, oneOf: []string{"offers", "review", "aggregateRating"}},
	"BreadcrumbList": {required: []string{"itemListElement"}},
	"LocalBusiness":  {required: []string{"name", "address"}},
	"Event":          {required: []string{"name", "startDate", "location"}},
	"Recipe":         {required: []string{"name", "image"}},
	"VideoObject":    {required: []string{"name", "thumbnailUrl", "uploadDate"}},
	"JobPosting":     {required: []string{"title", "description", "datePosted", "hiringOrganization"}},
}

// missingSchemaProperties describes what a structured data entity lacks, or
// returns nil if it is complete or its type has no known requirements.
func missingSchemaProperties(item models.SchemaItem) []string {
	requirement, known := requiredSchemaProperties[item.Type]
	if !known {
		return nil
	}

	present := make(map[string]bool, len(item.Properties))
	for _, property := range item.Properties {
		present[property] = true
	}

	var missing []string
	for _, property := range requirement.required {
		if !present[property] {
			missing = append(missing, property)
		}
	}

	if len(requirement.oneOf) > 0 {
		hasOne := false
		for _, property := range requirement.oneOf {
			if present[property] {
				hasOne = true
				break
			}
		}
		if !hasOne {
			missing = append(missing, "one of "+strings.Join(requirement.oneOf, ", "))
		}
	}

	return missing
}

// AnalyzeStructuredData reports unreadable JSON-LD and entities missing
// required properties. It also returns how many pages use each Schema.org type.
func AnalyzeStructuredData(results []*models.PageResult) ([]Issue, map[string]int) {
	var issues []Issue
	typeCounts := make(map[string]int)

	for _, p := range results {
		if !isHTMLPage(p) {
			continue
		}

		if len(p.SchemaErrors) > 0 {
			message := "Structured data could not be read: " + p.SchemaErrors[0]
			if len(p.SchemaErrors) > 1 {
				message += fmt.Sprintf(" (and %d more problem(s))", len(p.SchemaErrors)-1)
			}
			issues = append(issues, Issue{
				Type:           IssueInvalidStructuredData,
				Severity:       "error",
				URL:            p.URL,
				Message:        message,
				Value:          strings.Join(p.SchemaErrors, " | "),
				Recommendation: "Fix the JSON-LD so it is valid JSON with an @type; search engines ignore a block they cannot parse",
			})
		}

		// A page may declare the same type twice (JSON-LD and microdata, or two
		// blocks). Report a type once per page, and only if every declaration
		// of it is incomplete.
		typesOnPage := make(map[string]bool)
		missingByType := make(map[string][]string)
		complete := make(map[string]bool)
		for _, item := range p.SchemaItems {
			typesOnPage[item.Type] = true
			if missing := missingSchemaProperties(item); len(missing) > 0 {
				if _, seen := missingByType[item.Type]; !seen {
					missingByType[item.Type] = missing
				}
			} else {
				complete[item.Type] = true
			}
		}
		for typeName := range typesOnPage {
			typeCounts[typeName]++
		}

		incomplete := make([]string, 0, len(missingByType))
		for typeName := range missingByType {
			if !complete[typeName] {
				incomplete = append(incomplete, typeName)
			}
		}
		sort.Strings(incomplete)
		for _, typeName := range incomplete {
			issues = append(issues, Issue{
				Type:           IssueStructuredDataMissingFields,
				Severity:       "warning",
				URL:            p.URL,
				Message:        fmt.Sprintf("%s structured data is missing: %s", typeName, strings.Join(missingByType[typeName], "; ")),
				Value:          typeName,
				Recommendation: "Add the missing properties; without them the markup is not eligible for rich results",
			})
		}
	}

	return issues, typeCounts
}
