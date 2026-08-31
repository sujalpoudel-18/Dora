package models

import "time"

// Severity represents vulnerability severity levels
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityLow      Severity = "LOW"
	SeverityMedium   Severity = "MEDIUM"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

// Confidence represents detection confidence levels
type Confidence string

const (
	ConfidenceLow    Confidence = "LOW"
	ConfidenceMedium Confidence = "MEDIUM"
	ConfidenceHigh   Confidence = "HIGH"
)

// Evidence holds structured empirical evidence for a finding
type Evidence struct {
	What             string            `json:"what"`
	Where            string            `json:"where"`
	Why              string            `json:"why"`
	MatchedCondition string            `json:"matched_condition,omitempty"`
	Request          string            `json:"request,omitempty"`
	Response         string            `json:"response,omitempty"`
	StatusCode       int               `json:"status_code,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	ExtractedValues  map[string]string `json:"extracted_values,omitempty"`
}

// Finding represents a single security observation or vulnerability finding
type Finding struct {
	ID                 string     `json:"id"`
	TemplateID         string     `json:"template_id"`
	Name               string     `json:"name"`
	Severity           Severity   `json:"severity"`
	Confidence         Confidence `json:"confidence"`
	Target             string     `json:"target"`
	Host               string     `json:"host"`
	Port               int        `json:"port"`
	Protocol           string     `json:"protocol"`
	Endpoint           string     `json:"endpoint"`
	Method             string     `json:"method,omitempty"`
	Parameter          string     `json:"parameter,omitempty"`
	Technology         string     `json:"technology,omitempty"`
	Description        string     `json:"description"`
	Impact             string     `json:"impact,omitempty"`
	Remediation        string     `json:"remediation,omitempty"`
	References         []string   `json:"references,omitempty"`
	Evidence           Evidence   `json:"evidence"`
	Timestamp          time.Time  `json:"timestamp"`
	ScannerVersion     string     `json:"scanner_version"`
	CVSSScore          float64    `json:"cvss_score,omitempty"`
	CVSSVector         string     `json:"cvss_vector,omitempty"`

	// Database Subsystem Fields
	DatabaseType       string     `json:"database_type,omitempty"`
	DatabaseVersion    string     `json:"database_version,omitempty"`
	DatabaseDialect    string     `json:"database_dialect,omitempty"`
	Query              string     `json:"query,omitempty"`
	SourceFile         string     `json:"source_file,omitempty"`
	LineNumber         int        `json:"line_number,omitempty"`
	Schema             string     `json:"schema,omitempty"`
	Table              string     `json:"table,omitempty"`
	Function           string     `json:"function,omitempty"`
	Permission         string     `json:"permission,omitempty"`
	ConnectionSecurity string     `json:"connection_security,omitempty"`
}
