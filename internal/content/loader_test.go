package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPages(t *testing.T) {
	tmpDir := t.TempDir()
	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	// Hello.md with full frontmatter
	helloContent := `---
title: "Hello World"
titleKo: "안녕 세계"
type: Article
tags:
  - test
---
# Hello World

This is the first paragraph with enough content to extract a brief from it.
`
	if err := os.WriteFile(filepath.Join(pagesDir, "Hello.md"), []byte(helloContent), 0644); err != nil {
		t.Fatalf("WriteFile Hello.md error = %v", err)
	}

	// Secret.md with private:true
	secretContent := `---
title: "Secret Page"
private: true
---
This page is private.
`
	if err := os.WriteFile(filepath.Join(pagesDir, "Secret.md"), []byte(secretContent), 0644); err != nil {
		t.Fatalf("WriteFile Secret.md error = %v", err)
	}

	// Simple.md with no frontmatter
	simpleContent := `# Simple Page

This is a simple page without any frontmatter. It has some content here.
`
	if err := os.WriteFile(filepath.Join(pagesDir, "Simple.md"), []byte(simpleContent), 0644); err != nil {
		t.Fatalf("WriteFile Simple.md error = %v", err)
	}

	pages, err := LoadPages(pagesDir)
	if err != nil {
		t.Fatalf("LoadPages() error = %v", err)
	}

	if len(pages) != 3 {
		t.Fatalf("len(pages) = %d, want 3", len(pages))
	}

	// Build a map for easy lookup
	pageMap := make(map[string]*Page)
	for i := range pages {
		pageMap[pages[i].Name] = &pages[i]
	}

	// Check Hello page
	hello, ok := pageMap["Hello"]
	if !ok {
		t.Fatal("Hello page not found")
	}
	if hello.Title != "Hello World" {
		t.Errorf("Hello.Title = %q, want %q", hello.Title, "Hello World")
	}
	if hello.TitleKo != "안녕 세계" {
		t.Errorf("Hello.TitleKo = %q, want %q", hello.TitleKo, "안녕 세계")
	}

	// Check Secret page
	secret, ok := pageMap["Secret"]
	if !ok {
		t.Fatal("Secret page not found")
	}
	if !secret.Private {
		t.Error("Secret.Private = false, want true")
	}

	// Check Simple page
	simple, ok := pageMap["Simple"]
	if !ok {
		t.Fatal("Simple page not found")
	}
	// Title should default to filename
	if simple.Title != "Simple" {
		t.Errorf("Simple.Title = %q, want %q (filename)", simple.Title, "Simple")
	}
	// Brief should be non-empty
	if simple.Brief == "" {
		t.Error("Simple.Brief is empty, want non-empty")
	}
}
