package reporting

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dora/dora/internal/database/security"
	"github.com/dora/dora/pkg/models"
)

// PrintTerminal prints findings formatted for terminal output
func PrintTerminal(f *models.Finding, silent bool) {
	if silent {
		fmt.Printf("[%s] %s - %s\n", f.Severity, f.Name, f.Target)
		return
	}

	color := "\033[33m" // Yellow default
	switch f.Severity {
	case models.SeverityCritical:
		color = "\033[35m" // Magenta
	case models.SeverityHigh:
		color = "\033[31m" // Red
	case models.SeverityMedium:
		color = "\033[33m" // Yellow
	case models.SeverityLow:
		color = "\033[36m" // Cyan
	case models.SeverityInfo:
		color = "\033[32m" // Green
	}
	reset := "\033[0m"

	fmt.Printf("%s[%s]%s %s\n", color, f.Severity, reset, f.Name)
	fmt.Printf("  Target:     %s\n", f.Target)
	if f.Endpoint != "" {
		fmt.Printf("  Endpoint:   %s\n", f.Endpoint)
	}
	if f.SourceFile != "" {
		fmt.Printf("  Source:     %s:%d\n", f.SourceFile, f.LineNumber)
	}
	if f.Query != "" {
		fmt.Printf("  Query:      %s\n", security.RedactSecrets(f.Query))
	}
	fmt.Printf("  Confidence: %s\n", f.Confidence)
	if f.Remediation != "" {
		fmt.Printf("  Fix:        %s\n", f.Remediation)
	}
	fmt.Println()
}

// WriteJSON writes findings as formatted JSON array to file
func WriteJSON(filePath string, findings []models.Finding) error {
	redacted := make([]models.Finding, len(findings))
	for i, f := range findings {
		redacted[i] = f
		redacted[i].Query = security.RedactSecrets(f.Query)
	}

	data, err := json.MarshalIndent(redacted, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// WriteJSONL writes findings in JSON Lines format
func WriteJSONL(w io.Writer, findings []models.Finding) error {
	for _, f := range findings {
		f.Query = security.RedactSecrets(f.Query)
		data, err := json.Marshal(f)
		if err != nil {
			return err
		}
		if _, err := w.Write(append(data, '\n')); err != nil {
			return err
		}
	}
	return nil
}

// WriteCSV writes findings in CSV format
func WriteCSV(filePath string, findings []models.Finding) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = w.Write([]string{"ID", "TemplateID", "Name", "Severity", "Confidence", "Target", "Endpoint", "SourceFile", "LineNumber", "Description", "Remediation"})

	for _, fn := range findings {
		_ = w.Write([]string{
			fn.ID, fn.TemplateID, fn.Name, string(fn.Severity), string(fn.Confidence),
			fn.Target, fn.Endpoint, fn.SourceFile, fmt.Sprintf("%d", fn.LineNumber),
			fn.Description, fn.Remediation,
		})
	}

	return nil
}
