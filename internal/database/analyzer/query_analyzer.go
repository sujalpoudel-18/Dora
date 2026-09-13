package analyzer

import "strings"

type QueryIssue struct {
	Issue string
	Advice string
}

// AnalyzeQuery checks query construction for unparameterized values
func AnalyzeQuery(query string) []QueryIssue {
	var issues []QueryIssue
	if strings.Contains(query, " = '") || strings.Contains(query, " = \"") {
		issues = append(issues, QueryIssue{
			Issue:  "Possible dynamic string interpolation in query",
			Advice: "Use parameterized queries or prepared statements",
		})
	}
	return issues
}
