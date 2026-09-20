package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dora/dora/pkg/models"
)

type CodeFinding struct {
	Finding models.Finding
}

// Multi-language dynamic query concatenation regexes
var (
	// Concatenation patterns: db.Query("SELECT ... " + var)
	sqlConcatRegex = regexp.MustCompile(`(?i)(db|conn|sql|repo|stmt|client)\.(Query|Exec|Raw|Select|Where|Execute|Find|Model)\s*\([^)]*(\+|,|%\s*\(|fmt\.Sprintf|String\.format|f"|\$)\s*`)
	
	// Embedded Connection Strings
	connStrRegex = regexp.MustCompile(`(?i)(postgresql|postgres|mysql|mariadb|mongodb|mongodb\+srv|redis|sqlserver)://[^:]+:[^@]+@[^\s"']+`)
)

// ScanDirectory walks a codebase and performs multi-language database code analysis
func ScanDirectory(dirPath string) ([]models.Finding, error) {
	var findings []models.Finding

	supportedExts := map[string]string{
		".go":      "Go",
		".py":      "Python",
		".java":    "Java",
		".js":      "JavaScript",
		".ts":      "TypeScript",
		".php":     "PHP",
		".cs":      "C#",
		".rb":      "Ruby",
		".kt":      "Kotlin",
		".sql":     "SQL",
		".pgsql":   "PostgreSQL",
		".psql":    "PostgreSQL",
		".mysql":   "MySQL",
		".sqlite":  "SQLite",
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		lang, supported := supportedExts[ext]
		if !supported && !strings.Contains(strings.ToLower(info.Name()), "migration") {
			return nil
		}

		fileFindings, scanErr := ScanFile(path, lang)
		if scanErr == nil {
			findings = append(findings, fileFindings...)
		}

		return nil
	})

	return findings, err
}

// ScanFile reads a single source file and analyzes dynamic SQL concatenation or exposed secrets
func ScanFile(filePath string, language string) ([]models.Finding, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var findings []models.Finding
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "--") {
			continue
		}

		// Check 1: Dynamic SQL Concatenation / Unsafe ORM Query
		if sqlConcatRegex.MatchString(line) {
			findings = append(findings, models.Finding{
				ID:             "DORA-DB-S001",
				TemplateID:     "database/sql-injection/unsafe-dynamic-sql",
				Name:           "Unsafe Dynamic SQL Query Construction",
				Severity:       models.SeverityHigh,
				Confidence:     models.ConfidenceMedium,
				Target:         filePath,
				SourceFile:     filePath,
				LineNumber:     lineNum,
				Query:          trimmed,
				Description:    "Query built using string concatenation or format string instead of parameterization.",
				Remediation:    "Use parameterized queries, prepared statements, or ORM placeholder bindings.",
				ScannerVersion: "v0.1.0",
				Evidence: models.Evidence{
					What:  "Unsafe Dynamic SQL",
					Where: filePath,
					Why:   "Line contains string concatenation or Sprintf inside database query invocation.",
				},
			})
		}

		// Check 2: Embedded Database Connection Strings / Secrets
		if match := connStrRegex.FindString(line); match != "" {
			findings = append(findings, models.Finding{
				ID:             "DORA-DB-S002",
				TemplateID:     "database/secrets/exposed-connection-string",
				Name:           "Exposed Database Connection String",
				Severity:       models.SeverityCritical,
				Confidence:     models.ConfidenceHigh,
				Target:         filePath,
				SourceFile:     filePath,
				LineNumber:     lineNum,
				Description:    "Plain-text database connection string with embedded credentials found in source code.",
				Remediation:    "Move database credentials to environment variables or a secure key management system.",
				ScannerVersion: "v0.1.0",
				Evidence: models.Evidence{
					What:  "Hardcoded Connection String",
					Where: filePath,
					Why:   "Embedded database URI detected.",
				},
			})
		}
	}

	return findings, scanner.Err()
}
