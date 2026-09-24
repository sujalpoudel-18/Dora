package templates

import (
	"strings"

	"github.com/dora/dora/pkg/template"
)

type Manager struct {
	templates map[string]*template.Template
}

func NewManager() *Manager {
	return &Manager{
		templates: make(map[string]*template.Template),
	}
}

func (m *Manager) Add(t *template.Template) error {
	if err := ValidateTemplate(t); err != nil {
		return err
	}
	m.templates[t.ID] = t
	return nil
}

func (m *Manager) List() []*template.Template {
	res := make([]*template.Template, 0, len(m.templates))
	for _, t := range m.templates {
		res = append(res, t)
	}
	return res
}

func (m *Manager) Search(keyword string) []*template.Template {
	kw := strings.ToLower(keyword)
	var matches []*template.Template

	for _, t := range m.templates {
		if strings.Contains(strings.ToLower(t.ID), kw) ||
			strings.Contains(strings.ToLower(t.Name), kw) ||
			strings.Contains(strings.ToLower(t.Description), kw) {
			matches = append(matches, t)
		}
	}
	return matches
}

func (m *Manager) FilterByTags(tags []string) []*template.Template {
	var filtered []*template.Template
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[strings.ToLower(strings.TrimSpace(tag))] = true
	}

	for _, t := range m.templates {
		for _, tag := range t.Tags {
			if tagMap[strings.ToLower(tag)] {
				filtered = append(filtered, t)
				break
			}
		}
	}
	return filtered
}

func (m *Manager) Get(id string) (*template.Template, bool) {
	t, ok := m.templates[id]
	return t, ok
}

func (m *Manager) Count() int {
	return len(m.templates)
}
