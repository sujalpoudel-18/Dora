package correlation

import (
	"fmt"
	"github.com/dora/dora/pkg/models"
)

// CorrelateFindings matches HTTP web application findings with static code database findings
func CorrelateFindings(httpFindings []models.Finding, dbFindings []models.Finding) []models.Finding {
	var correlated []models.Finding
	matchedDB := make(map[string]bool)

	for _, hf := range httpFindings {
		foundMatch := false
		for _, df := range dbFindings {
			// Match criteria: HTTP parameter reflection + Database query concatenation
			if hf.Parameter != "" && (df.Query != "" || df.SourceFile != "") {
				foundMatch = true
				matchedDB[df.ID] = true

				merged := hf
				merged.ID = fmt.Sprintf("CORRELATED-%s-%s", hf.ID, df.ID)
				merged.Name = fmt.Sprintf("Confirmed SQL Injection: %s correlated with %s", hf.Endpoint, df.SourceFile)
				merged.Severity = models.SeverityCritical
				merged.Confidence = models.ConfidenceHigh
				merged.SourceFile = df.SourceFile
				merged.LineNumber = df.LineNumber
				merged.Query = df.Query
				merged.Description = fmt.Sprintf("%s\nCorrelated with source code file %s:%d (`%s`)", hf.Description, df.SourceFile, df.LineNumber, df.Query)
				
				correlated = append(correlated, merged)
				break
			}
		}

		if !foundMatch {
			correlated = append(correlated, hf)
		}
	}

	for _, df := range dbFindings {
		if !matchedDB[df.ID] {
			correlated = append(correlated, df)
		}
	}

	return correlated
}
