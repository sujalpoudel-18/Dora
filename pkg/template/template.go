package template

import "time"

// InfoBlock holds metadata about a template
type InfoBlock struct {
	Name             string    `yaml:"name" json:"name"`
	Author           string    `yaml:"author" json:"author"`
	Description      string    `yaml:"description" json:"description"`
	Severity         string    `yaml:"severity" json:"severity"`
	Impact           string    `yaml:"impact,omitempty" json:"impact,omitempty"`
	Remediation      string    `yaml:"remediation,omitempty" json:"remediation,omitempty"`
	Category         string    `yaml:"category,omitempty" json:"category,omitempty"`
	Tags             []string  `yaml:"tags,omitempty" json:"tags,omitempty"`
	Technology       string    `yaml:"technology,omitempty" json:"technology,omitempty"`
	Product          string    `yaml:"product,omitempty" json:"product,omitempty"`
	Version          string    `yaml:"version,omitempty" json:"version,omitempty"`
	AffectedVersions string    `yaml:"affected_versions,omitempty" json:"affected_versions,omitempty"`
	CWE              string    `yaml:"cwe,omitempty" json:"cwe,omitempty"`
	CVE              string    `yaml:"cve,omitempty" json:"cve,omitempty"`
	OWASP            string    `yaml:"owasp,omitempty" json:"owasp,omitempty"`
	Created          time.Time `yaml:"created,omitempty" json:"created,omitempty"`
	Updated          time.Time `yaml:"updated,omitempty" json:"updated,omitempty"`
}

// Template represents a modular security assessment template
type Template struct {
	ID               string                 `yaml:"id" json:"id"`
	Name             string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Author           string                 `yaml:"author,omitempty" json:"author,omitempty"`
	Description      string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Category         string                 `yaml:"category,omitempty" json:"category,omitempty"`
	Tags             []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
	Severity         string                 `yaml:"severity,omitempty" json:"severity,omitempty"`
	Confidence       string                 `yaml:"confidence,omitempty" json:"confidence,omitempty"`
	Technology       string                 `yaml:"technology,omitempty" json:"technology,omitempty"`
	Product          string                 `yaml:"product,omitempty" json:"product,omitempty"`
	Version          string                 `yaml:"version,omitempty" json:"version,omitempty"`
	AffectedVersions string                 `yaml:"affected_versions,omitempty" json:"affected_versions,omitempty"`
	CWE              string                 `yaml:"cwe,omitempty" json:"cwe,omitempty"`
	CVE              string                 `yaml:"cve,omitempty" json:"cve,omitempty"`
	OWASP            string                 `yaml:"owasp,omitempty" json:"owasp,omitempty"`
	Created          time.Time              `yaml:"created,omitempty" json:"created,omitempty"`
	Updated          time.Time              `yaml:"updated,omitempty" json:"updated,omitempty"`

	Info InfoBlock `yaml:"info,omitempty" json:"info,omitempty"`

	// Target Protocol Specifications
	HTTP     []HTTPRequestBlock   `yaml:"http,omitempty" json:"http,omitempty"`
	TCP      []TCPRequestBlock    `yaml:"tcp,omitempty" json:"tcp,omitempty"`
	DNS      []DNSRequestBlock    `yaml:"dns,omitempty" json:"dns,omitempty"`
	TLS      []TLSCheckBlock      `yaml:"tls,omitempty" json:"tls,omitempty"`
	Code     []CodeAnalysisBlock  `yaml:"code,omitempty" json:"code,omitempty"`
	Database []DatabaseCheckBlock `yaml:"database,omitempty" json:"database,omitempty"`
	Cloud    []CloudCheckBlock    `yaml:"cloud,omitempty" json:"cloud,omitempty"`
	K8s      []K8sCheckBlock      `yaml:"k8s,omitempty" json:"k8s,omitempty"`

	// Logic & Rules
	Matchers          []Matcher          `yaml:"matchers,omitempty" json:"matchers,omitempty"`
	MatchersCondition string             `yaml:"matchers-condition,omitempty" json:"matchers_condition,omitempty"`
	Extractors        []Extractor        `yaml:"extractors,omitempty" json:"extractors,omitempty"`
	Verification      VerificationConfig `yaml:"verification,omitempty" json:"verification,omitempty"`
	Remediation       string             `yaml:"remediation,omitempty" json:"remediation,omitempty"`
	References        []string           `yaml:"references,omitempty" json:"references,omitempty"`
}

// Normalize ensures metadata fields in Info block are populated at top level if nested
func (t *Template) Normalize() {
	if t.Info.Name != "" && t.Name == "" {
		t.Name = t.Info.Name
	}
	if t.Info.Author != "" && t.Author == "" {
		t.Author = t.Info.Author
	}
	if t.Info.Description != "" && t.Description == "" {
		t.Description = t.Info.Description
	}
	if t.Info.Severity != "" && t.Severity == "" {
		t.Severity = t.Info.Severity
	}
	if len(t.Info.Tags) > 0 && len(t.Tags) == 0 {
		t.Tags = t.Info.Tags
	}
}

// HTTPRequestBlock defines HTTP request attributes for scanning
type HTTPRequestBlock struct {
	Method            string            `yaml:"method" json:"method"`
	Path              []string          `yaml:"path" json:"path"`
	Headers           map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body              string            `yaml:"body,omitempty" json:"body,omitempty"`
	Matchers          []Matcher         `yaml:"matchers,omitempty" json:"matchers,omitempty"`
	MatchersCondition string            `yaml:"matchers-condition,omitempty" json:"matchers_condition,omitempty"`
	Extractors        []Extractor       `yaml:"extractors,omitempty" json:"extractors,omitempty"`
}

// TCPRequestBlock defines raw TCP probing
type TCPRequestBlock struct {
	Port     int        `yaml:"port" json:"port"`
	Inputs   []TCPInput `yaml:"inputs,omitempty" json:"inputs,omitempty"`
	Matchers []Matcher  `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

type TCPInput struct {
	Data string `yaml:"data" json:"data"`
}

// DNSRequestBlock defines DNS record queries
type DNSRequestBlock struct {
	Types    []string  `yaml:"types" json:"types"`
	Matchers []Matcher `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

// TLSCheckBlock defines TLS configuration checks
type TLSCheckBlock struct {
	MinVersion string    `yaml:"min_version,omitempty" json:"min_version,omitempty"`
	Matchers   []Matcher `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

// CodeAnalysisBlock defines static code rule patterns
type CodeAnalysisBlock struct {
	Languages []string  `yaml:"languages" json:"languages"`
	Pattern   string    `yaml:"pattern" json:"pattern"`
	Matchers  []Matcher `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

// DatabaseCheckBlock defines database service or SQL analysis checks
type DatabaseCheckBlock struct {
	Type     string    `yaml:"type" json:"type"`
	Query    string    `yaml:"query,omitempty" json:"query,omitempty"`
	Matchers []Matcher `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

// CloudCheckBlock defines cloud security check rules
type CloudCheckBlock struct {
	Provider string    `yaml:"provider" json:"provider"`
	Resource string    `yaml:"resource" json:"resource"`
	Matchers []Matcher `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

// K8sCheckBlock defines Kubernetes security check rules
type K8sCheckBlock struct {
	Resource string    `yaml:"resource" json:"resource"`
	Matchers []Matcher `yaml:"matchers,omitempty" json:"matchers,omitempty"`
}

// Matcher defines rule match conditions
type Matcher struct {
	Type        string   `yaml:"type" json:"type"`
	Part        string   `yaml:"part,omitempty" json:"part,omitempty"`
	Words       []string `yaml:"words,omitempty" json:"words,omitempty"`
	Regex       []string `yaml:"regex,omitempty" json:"regex,omitempty"`
	Status      []int    `yaml:"status,omitempty" json:"status,omitempty"`
	Size        []int    `yaml:"size,omitempty" json:"size,omitempty"`
	Condition   string   `yaml:"condition,omitempty" json:"condition,omitempty"`
	Negative    bool     `yaml:"negative,omitempty" json:"negative,omitempty"`
	ContentType []string `yaml:"content_type,omitempty" json:"content_type,omitempty"`
}

// Extractor defines data extraction rules
type Extractor struct {
	Type      string   `yaml:"type" json:"type"`
	Name      string   `yaml:"name" json:"name"`
	Part      string   `yaml:"part,omitempty" json:"part,omitempty"`
	Regex     []string `yaml:"regex,omitempty" json:"regex,omitempty"`
	Header    string   `yaml:"header,omitempty" json:"header,omitempty"`
	Attribute string   `yaml:"attribute,omitempty" json:"attribute,omitempty"`
}

// VerificationConfig defines verification modes
type VerificationConfig struct {
	Mode           string `yaml:"mode" json:"mode"` // safe, baseline, callback
	InertMarker    string `yaml:"inert_marker,omitempty" json:"inert_marker,omitempty"`
	CallbackServer string `yaml:"callback_server,omitempty" json:"callback_server,omitempty"`
}
