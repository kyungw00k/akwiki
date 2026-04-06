package render

import (
	"bytes"
	"embed"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/content"
	"github.com/kyungw00k/akwiki/internal/wiki"
)

// TemplateContext holds all data passed to the page template.
type TemplateContext struct {
	Site      *config.Config
	Page      content.Page
	Content   template.HTML
	TOC       []wiki.Heading
	Links     []PageRef
	Backlinks []PageRef
	Related   map[string][]PageRef
	JSONLD    template.JS
	ThemeCSS  template.CSS
}

// PageRef represents a reference to another wiki page.
type PageRef struct {
	Name  string
	Title string
	Brief string
	Type  string
	Score float64
}

// funcMap provides template helper functions.
var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"sameDay": func(a, b time.Time) bool {
		y1, m1, d1 := a.Date()
		y2, m2, d2 := b.Date()
		return y1 == y2 && m1 == m2 && d1 == d2
	},
	"editURL": func(pattern, pagename string) string {
		return strings.ReplaceAll(pattern, "{{pagename}}", pagename)
	},
	"urlPathEscape": func(s string) string {
		return url.PathEscape(s)
	},
	"relatedLabel": func(typ string) string {
		switch typ {
		case "Book":
			return "관련 책"
		case "Article":
			return "관련 글"
		case "ScholarlyArticle":
			return "관련 논문"
		default:
			return "관련 문서"
		}
	},
}

// TemplateEngine renders wiki pages using Go html/template.
type TemplateEngine struct {
	tmpl *template.Template
}

// templateFiles defines the loading order: partials first, then main page template.
var templateFiles = []struct {
	name    string
	embPath string // path inside embed.FS
}{
	{"header", "default/templates/partials/header.html"},
	{"toc", "default/templates/partials/toc.html"},
	{"backlinks", "default/templates/partials/backlinks.html"},
	{"related", "default/templates/partials/related.html"},
	{"search", "default/templates/partials/search.html"},
	{"footer", "default/templates/partials/footer.html"},
	{"page", "default/templates/page.html"},
}

// NewTemplateEngine creates a TemplateEngine by loading templates.
// For each template file, it first checks themeOverrideDir for a user override.
// If not found, it falls back to the embedded defaultFS.
func NewTemplateEngine(defaultFS embed.FS, themeOverrideDir string) (*TemplateEngine, error) {
	tmpl := template.New("").Funcs(funcMap)

	for _, tf := range templateFiles {
		data, err := loadTemplateFile(defaultFS, tf.embPath, themeOverrideDir)
		if err != nil {
			return nil, err
		}
		if _, err := tmpl.Parse(string(data)); err != nil {
			return nil, err
		}
	}

	return &TemplateEngine{tmpl: tmpl}, nil
}

// loadTemplateFile tries to read from themeOverrideDir first, then falls back to embed.FS.
func loadTemplateFile(defaultFS embed.FS, embPath string, themeOverrideDir string) ([]byte, error) {
	if themeOverrideDir != "" {
		// Map embed path to override path: "default/templates/partials/header.html" → "partials/header.html"
		// or "default/templates/page.html" → "page.html"
		relPath := strings.TrimPrefix(embPath, "default/templates/")
		overridePath := filepath.Join(themeOverrideDir, relPath)
		if data, err := os.ReadFile(overridePath); err == nil {
			return data, nil
		}
	}
	return defaultFS.ReadFile(embPath)
}

// RenderPage renders a wiki page to HTML bytes using the loaded templates.
func (e *TemplateEngine) RenderPage(ctx *TemplateContext) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.tmpl.ExecuteTemplate(&buf, "page", ctx); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
