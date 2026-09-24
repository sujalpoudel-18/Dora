package templates

import (
	"path/filepath"
	"testing"
)

func TestLoaderAndValidator(t *testing.T) {
	loader := NewLoader()
	manager := NewManager()

	// Load templates from templates directory
	templates, err := loader.LoadDir(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	if len(templates) == 0 {
		t.Fatalf("Expected templates to be loaded, got 0")
	}

	for _, tmpl := range templates {
		if err := manager.Add(tmpl); err != nil {
			t.Errorf("Validation failed for template %s: %v", tmpl.ID, err)
		}
	}

	if manager.Count() < 4 {
		t.Errorf("Expected at least 4 loaded templates, got %d", manager.Count())
	}

	// Test Search
	results := manager.Search("exposure")
	if len(results) == 0 {
		t.Errorf("Expected search results for 'exposure', got 0")
	}

	// Test FilterByTags
	filtered := manager.FilterByTags([]string{"xss"})
	if len(filtered) == 0 {
		t.Errorf("Expected filtered templates for tag 'xss', got 0")
	}
}
