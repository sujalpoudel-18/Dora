package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dora/dora/pkg/template"
	"gopkg.in/yaml.v3"
)

// Loader walks directories and loads YAML template definitions
type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

// LoadFile reads a single YAML template file
func (l *Loader) LoadFile(path string) (*template.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	var tmpl template.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template YAML %s: %w", path, err)
	}
	tmpl.Normalize()

	return &tmpl, nil
}

// LoadDir recursively loads all .yaml and .yml templates from a directory
func (l *Loader) LoadDir(dirPath string) ([]*template.Template, error) {
	var templates []*template.Template

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".yml" {
			tmpl, loadErr := l.LoadFile(path)
			if loadErr == nil && tmpl.ID != "" {
				templates = append(templates, tmpl)
			}
		}
		return nil
	})

	return templates, err
}
