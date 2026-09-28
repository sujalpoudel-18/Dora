package fingerprint

import (
	"strings"

	"github.com/dora/dora/internal/protocols/http"
	"github.com/dora/dora/pkg/models"
)

type Technology struct {
	Name       string            `json:"name"`
	Category   string            `json:"category"`
	Confidence models.Confidence `json:"confidence"`
	Version    string            `json:"version,omitempty"`
}

// Fingerprinter runs multi-signal detection across HTTP headers, cookies, and body content
type Fingerprinter struct{}

func NewFingerprinter() *Fingerprinter {
	return &Fingerprinter{}
}

// Identify identifies technologies from an HTTP response
func (f *Fingerprinter) Identify(resp *http.Response) []Technology {
	var techs []Technology
	serverHeader := strings.ToLower(resp.Headers["Server"])
	poweredBy := strings.ToLower(resp.Headers["X-Powered-By"])
	bodyLower := strings.ToLower(resp.Body)

	// Nginx
	if strings.Contains(serverHeader, "nginx") {
		techs = append(techs, Technology{
			Name:       "Nginx",
			Category:   "Web Server",
			Confidence: models.ConfidenceHigh,
		})
	}

	// Apache
	if strings.Contains(serverHeader, "apache") {
		techs = append(techs, Technology{
			Name:       "Apache",
			Category:   "Web Server",
			Confidence: models.ConfidenceHigh,
		})
	}

	// Express / Node.js
	if strings.Contains(poweredBy, "express") {
		techs = append(techs, Technology{
			Name:       "Express",
			Category:   "Web Framework",
			Confidence: models.ConfidenceHigh,
		})
		techs = append(techs, Technology{
			Name:       "Node.js",
			Category:   "Runtime",
			Confidence: models.ConfidenceMedium,
		})
	}

	// Django / Python
	if strings.Contains(resp.Headers["Set-Cookie"], "csrftoken") || strings.Contains(bodyLower, "csrfmiddlewaretoken") {
		techs = append(techs, Technology{
			Name:       "Django",
			Category:   "Web Framework",
			Confidence: models.ConfidenceHigh,
		})
	}

	// WordPress
	if strings.Contains(bodyLower, "wp-content") || strings.Contains(bodyLower, "wp-includes") {
		techs = append(techs, Technology{
			Name:       "WordPress",
			Category:   "CMS",
			Confidence: models.ConfidenceHigh,
		})
	}

	// Next.js / React
	if strings.Contains(bodyLower, "__next") || strings.Contains(bodyLower, "_next/static") {
		techs = append(techs, Technology{
			Name:       "Next.js",
			Category:   "Frontend Framework",
			Confidence: models.ConfidenceHigh,
		})
		techs = append(techs, Technology{
			Name:       "React",
			Category:   "JavaScript Library",
			Confidence: models.ConfidenceHigh,
		})
	}

	return techs
}
