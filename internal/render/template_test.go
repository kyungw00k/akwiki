package render

import (
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/content"
	"github.com/kyungw00k/akwiki/internal/wiki"
	"github.com/kyungw00k/akwiki/theme"
)

func TestRenderPage(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{Title: "Test Wiki", Language: "ko"},
		Theme: config.ThemeConfig{
			Layout: config.ThemeLayout{TOC: true, Backlinks: true, Related: true, Search: true},
		},
	}
	now := time.Now()
	ctx := &TemplateContext{
		Site:    &cfg,
		Page:    content.Page{Name: "Hello", Title: "Hello World", CreatedAt: now, ModifiedAt: now},
		Content: template.HTML("<p>Hello</p>"),
		TOC:     []wiki.Heading{{Level: 2, Text: "Section", ID: "h1234"}},
		JSONLD:  template.JS(`{"@type":"Article"}`),
	}
	engine, err := NewTemplateEngine(theme.DefaultTheme, "")
	if err != nil {
		t.Fatalf("NewTemplateEngine failed: %v", err)
	}

	output, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	html := string(output)

	checks := []string{
		"Hello World",
		"위키 홈",
		"본문으로 건너뛰기",
		"akngs",
		"Section",
		"<p>Hello</p>",
		`lang="ko"`,
		"Test Wiki",
		`role="search"`,
	}
	for _, want := range checks {
		if !strings.Contains(html, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestRenderPage_NoTOC(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{Title: "Test Wiki", Language: "ko"},
		Theme: config.ThemeConfig{
			Layout: config.ThemeLayout{TOC: false, Backlinks: false, Related: false, Search: false},
		},
	}
	now := time.Now()
	ctx := &TemplateContext{
		Site:    &cfg,
		Page:    content.Page{Name: "Hello", Title: "Hello World", CreatedAt: now, ModifiedAt: now},
		Content: template.HTML("<p>Hello</p>"),
		JSONLD:  template.JS(`{}`),
	}
	engine, err := NewTemplateEngine(theme.DefaultTheme, "")
	if err != nil {
		t.Fatalf("NewTemplateEngine failed: %v", err)
	}

	output, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	html := string(output)

	// TOC disabled, should not contain TOC section
	if strings.Contains(html, `aria-label="목차"`) {
		t.Error("TOC should not be rendered when disabled")
	}
	// Search disabled, should not contain search form
	if strings.Contains(html, `role="search"`) {
		t.Error("search should not be rendered when disabled")
	}
}

func TestRenderPage_WithBacklinks(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{Title: "Test Wiki", Language: "ko", URL: "https://example.com"},
		Theme: config.ThemeConfig{
			Layout: config.ThemeLayout{Backlinks: true},
		},
	}
	now := time.Now()
	ctx := &TemplateContext{
		Site: &cfg,
		Page: content.Page{Name: "Hello", Title: "Hello World", CreatedAt: now, ModifiedAt: now},
		Content: template.HTML("<p>Hello</p>"),
		Backlinks: []PageRef{
			{Name: "OtherPage", Title: "Other Page"},
		},
		JSONLD: template.JS(`{}`),
	}
	engine, err := NewTemplateEngine(theme.DefaultTheme, "")
	if err != nil {
		t.Fatalf("NewTemplateEngine failed: %v", err)
	}

	output, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	html := string(output)
	if !strings.Contains(html, "역링크") {
		t.Error("output should contain backlinks section")
	}
	if !strings.Contains(html, "Other Page") {
		t.Error("output should contain backlink title")
	}
}

func TestRenderPage_WithRelated(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{Title: "Test Wiki", Language: "ko", URL: "https://example.com"},
		Theme: config.ThemeConfig{
			Layout: config.ThemeLayout{Related: true},
		},
	}
	now := time.Now()
	ctx := &TemplateContext{
		Site:    &cfg,
		Page:    content.Page{Name: "Hello", Title: "Hello World", CreatedAt: now, ModifiedAt: now},
		Content: template.HTML("<p>Hello</p>"),
		Related: map[string][]PageRef{
			"Book": {{Name: "GoBook", Title: "Go Programming"}},
		},
		JSONLD: template.JS(`{}`),
	}
	engine, err := NewTemplateEngine(theme.DefaultTheme, "")
	if err != nil {
		t.Fatalf("NewTemplateEngine failed: %v", err)
	}

	output, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	html := string(output)
	if !strings.Contains(html, "관련 책") {
		t.Error("output should contain related label for Book type")
	}
	if !strings.Contains(html, "Go Programming") {
		t.Error("output should contain related page title")
	}
}

func TestRenderPage_EditURL(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{Title: "Test Wiki", Language: "ko"},
		Theme: config.ThemeConfig{
			Edit: config.ThemeEdit{URL: "https://github.com/repo/edit/main/{{pagename}}.md"},
		},
	}
	now := time.Now()
	ctx := &TemplateContext{
		Site:    &cfg,
		Page:    content.Page{Name: "Hello", Title: "Hello World", CreatedAt: now, ModifiedAt: now},
		Content: template.HTML("<p>Hello</p>"),
		JSONLD:  template.JS(`{}`),
	}
	engine, err := NewTemplateEngine(theme.DefaultTheme, "")
	if err != nil {
		t.Fatalf("NewTemplateEngine failed: %v", err)
	}

	output, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	html := string(output)
	if !strings.Contains(html, "https://github.com/repo/edit/main/Hello.md") {
		t.Error("output should contain resolved edit URL")
	}
}

func TestRenderPage_DifferentDates(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{Language: "ko"},
	}
	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	modified := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	ctx := &TemplateContext{
		Site:    &cfg,
		Page:    content.Page{Name: "Hello", Title: "Hello World", CreatedAt: created, ModifiedAt: modified},
		Content: template.HTML("<p>Hello</p>"),
		JSONLD:  template.JS(`{}`),
	}
	engine, err := NewTemplateEngine(theme.DefaultTheme, "")
	if err != nil {
		t.Fatalf("NewTemplateEngine failed: %v", err)
	}

	output, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	html := string(output)
	if !strings.Contains(html, "2024-01-01") {
		t.Error("output should contain created date")
	}
	if !strings.Contains(html, "modified: 2024-06-15") {
		t.Error("output should contain modified date when different from created")
	}
}
