package analyzer

import "strings"

type MigrationIssue struct {
	File  string
	Issue string
}

// AnalyzeMigration inspects migration script contents for plain-text secrets or risky statements
func AnalyzeMigration(filename, content string) []MigrationIssue {
	var issues []MigrationIssue
	upper := strings.ToUpper(content)

	if strings.Contains(upper, "PASSWORD") && (strings.Contains(upper, "IDENTIFIED BY") || strings.Contains(upper, "WITH PASSWORD")) {
		issues = append(issues, MigrationIssue{
			File:  filename,
			Issue: "Hardcoded password pattern found in migration script",
		})
	}
	return issues
}
