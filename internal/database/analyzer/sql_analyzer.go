package analyzer

import "strings"

type SQLRisk struct {
	RuleID      string
	Description string
	Severity    string
}

// AnalyzeSQL inspects SQL query strings for dangerous patterns
func AnalyzeSQL(query string) []SQLRisk {
	var risks []SQLRisk
	upper := strings.ToUpper(query)

	if strings.Contains(upper, "DROP TABLE") || strings.Contains(upper, "DROP DATABASE") {
		risks = append(risks, SQLRisk{
			RuleID:      "DORA-DB-001",
			Description: "Destructive DROP DDL statement detected",
			Severity:    "HIGH",
		})
	}

	if strings.Contains(upper, "GRANT ALL") || strings.Contains(upper, "SUPERUSER") {
		risks = append(risks, SQLRisk{
			RuleID:      "DORA-DB-002",
			Description: "Excessive grant or privilege escalation statement detected",
			Severity:    "MEDIUM",
		})
	}

	if strings.Contains(upper, "UNION SELECT") {
		risks = append(risks, SQLRisk{
			RuleID:      "DORA-DB-003",
			Description: "Potential UNION-based SQL query pattern detected",
			Severity:    "HIGH",
		})
	}

	return risks
}
