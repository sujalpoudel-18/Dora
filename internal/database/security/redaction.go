package security

import (
	"regexp"
	"strings"
)

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|pass|pwd|secret|token|key|api_key|auth|connectionstring|connstr)\s*[:=]\s*['"]?([^'"\s;,]+)['"]?`),
	regexp.MustCompile(`(?i)(postgresql|postgres|mysql|mariadb|mongodb|mongodb\+srv|redis|sqlserver)://([^:]+):([^@]+)@`),
}

// RedactSecrets replaces sensitive passwords, tokens, and credentials with REDACTED
func RedactSecrets(input string) string {
	if input == "" {
		return ""
	}

	result := input
	// Connection string URI password redaction: scheme://user:pass@host
	uriRegex := regexp.MustCompile(`(?i)(postgresql|postgres|mysql|mariadb|mongodb|mongodb\+srv|redis|sqlserver)://([^:]+):([^@]+)@`)
	result = uriRegex.ReplaceAllString(result, "${1}://${2}:REDACTED@")

	// Key-value pair redaction: password=secret
	kvRegex := regexp.MustCompile(`(?i)\b(password|pass|pwd|secret|token|key|api_key|auth)=['"]?([^'"\s;,]+)['"]?`)
	result = kvRegex.ReplaceAllString(result, "${1}=REDACTED")

	return result
}

// RedactMap performs secret redaction on key-value string maps
func RedactMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	redacted := make(map[string]string, len(m))
	for k, v := range m {
		lowerK := strings.ToLower(k)
		if strings.Contains(lowerK, "password") || strings.Contains(lowerK, "secret") || strings.Contains(lowerK, "token") || strings.Contains(lowerK, "auth") {
			redacted[k] = "REDACTED"
		} else {
			redacted[k] = RedactSecrets(v)
		}
	}
	return redacted
}
