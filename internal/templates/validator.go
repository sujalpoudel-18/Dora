package templates

import (
	"fmt"
	"strings"

	"github.com/dora/dora/pkg/template"
)

type ValidationError struct {
	TemplateID string
	Field      string
	Message    string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Template [%s] validation error on field '%s': %s", e.TemplateID, e.Field, e.Message)
}

// ValidateTemplate checks required template schema fields
func ValidateTemplate(t *template.Template) error {
	if strings.TrimSpace(t.ID) == "" {
		return ValidationError{TemplateID: "UNKNOWN", Field: "id", Message: "template ID cannot be empty"}
	}
	if strings.TrimSpace(t.Name) == "" {
		return ValidationError{TemplateID: t.ID, Field: "name", Message: "template name cannot be empty"}
	}
	if strings.TrimSpace(t.Author) == "" {
		return ValidationError{TemplateID: t.ID, Field: "author", Message: "template author cannot be empty"}
	}
	if strings.TrimSpace(t.Severity) == "" {
		return ValidationError{TemplateID: t.ID, Field: "severity", Message: "template severity cannot be empty"}
	}

	validSeverities := map[string]bool{
		"info": true, "low": true, "medium": true, "high": true, "critical": true,
	}
	if !validSeverities[strings.ToLower(t.Severity)] {
		return ValidationError{TemplateID: t.ID, Field: "severity", Message: "invalid severity level (must be info, low, medium, high, or critical)"}
	}

	return nil
}
