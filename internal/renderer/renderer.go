package renderer

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Renderer loads, caches, and executes Go text/templates from an embed.FS.
type Renderer struct {
	fs    embed.FS
	cache sync.Map // map[string]*template.Template
}

// New creates a Renderer backed by the given embedded filesystem.
func New(fs embed.FS) *Renderer {
	return &Renderer{fs: fs}
}

// Render executes the named template with the given data and returns the result.
func (r *Renderer) Render(templatePath string, data any) (string, error) {
	tmpl, err := r.load(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderBytes is like Render but returns bytes.
func (r *Renderer) RenderBytes(templatePath string, data any) ([]byte, error) {
	s, err := r.Render(templatePath, data)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func (r *Renderer) load(templatePath string) (*template.Template, error) {
	if cached, ok := r.cache.Load(templatePath); ok {
		return cached.(*template.Template), nil
	}

	content, err := r.fs.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("reading template %s: %w", templatePath, err)
	}

	name := filepath.Base(templatePath)
	name = strings.TrimSuffix(name, ".tmpl")

	tmpl, err := template.New(name).Funcs(FuncMap()).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing template %s: %w", templatePath, err)
	}

	r.cache.Store(templatePath, tmpl)
	return tmpl, nil
}
