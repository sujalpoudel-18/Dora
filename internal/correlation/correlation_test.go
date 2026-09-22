package correlation

import (
	"testing"
	"github.com/dora/dora/pkg/models"
)

func TestCorrelateFindings(t *testing.T) {
	httpFindings := []models.Finding{
		{
			ID:        "HTTP-01",
			Endpoint:  "/login",
			Parameter: "username",
		},
	}

	dbFindings := []models.Finding{
		{
			ID:         "DB-01",
			SourceFile: "auth.go",
			LineNumber: 42,
			Query:      "SELECT * FROM users WHERE username = '" + "username" + "'",
		},
	}

	correlated := CorrelateFindings(httpFindings, dbFindings)
	if len(correlated) != 1 {
		t.Fatalf("Expected 1 correlated finding, got %d", len(correlated))
	}

	if correlated[0].Severity != models.SeverityCritical {
		t.Errorf("Expected Critical severity for correlated finding, got %v", correlated[0].Severity)
	}
}
