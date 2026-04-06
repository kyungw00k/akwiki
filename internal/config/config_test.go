package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Site.Language != "ko" {
		t.Errorf("Language = %q, want %q", cfg.Site.Language, "ko")
	}
	if cfg.Build.OutDir != "dist" {
		t.Errorf("OutDir = %q, want %q", cfg.Build.OutDir, "dist")
	}
	if cfg.Build.PageRoute != "/pages" {
		t.Errorf("PageRoute = %q, want %q", cfg.Build.PageRoute, "/pages")
	}
	if !cfg.Theme.Layout.TOC {
		t.Error("TOC = false, want true")
	}
	if !cfg.Theme.Layout.Backlinks {
		t.Error("Backlinks = false, want true")
	}
	if !cfg.Theme.Layout.Related {
		t.Error("Related = false, want true")
	}
	if !cfg.Theme.Layout.Search {
		t.Error("Search = false, want true")
	}
}

func TestLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	akwikiDir := filepath.Join(tmpDir, ".akwiki")
	if err := os.MkdirAll(akwikiDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	configContent := `
site:
  title: "My Wiki"
  author: "Test Author"
  language: "en"
build:
  outDir: "output"
theme:
  layout:
    toc: false
`
	configPath := filepath.Join(akwikiDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify overrides
	if cfg.Site.Title != "My Wiki" {
		t.Errorf("Title = %q, want %q", cfg.Site.Title, "My Wiki")
	}
	if cfg.Site.Author != "Test Author" {
		t.Errorf("Author = %q, want %q", cfg.Site.Author, "Test Author")
	}
	if cfg.Site.Language != "en" {
		t.Errorf("Language = %q, want %q", cfg.Site.Language, "en")
	}
	if cfg.Build.OutDir != "output" {
		t.Errorf("OutDir = %q, want %q", cfg.Build.OutDir, "output")
	}
	if cfg.Theme.Layout.TOC != false {
		t.Error("TOC = true, want false")
	}

	// Verify unset fields keep defaults
	if cfg.Build.PageRoute != "/pages" {
		t.Errorf("PageRoute = %q, want %q (default)", cfg.Build.PageRoute, "/pages")
	}
	if !cfg.Theme.Layout.Backlinks {
		t.Error("Backlinks = false, want true (default)")
	}
	if !cfg.Theme.Layout.Related {
		t.Error("Related = false, want true (default)")
	}
	if !cfg.Theme.Layout.Search {
		t.Error("Search = false, want true (default)")
	}
}
