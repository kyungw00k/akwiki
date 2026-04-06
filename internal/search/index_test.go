package search

import (
	"encoding/json"
	"testing"

	"github.com/kyungw00k/akwiki/internal/content"
)

func TestBuildSearchIndex(t *testing.T) {
	pages := []content.Page{
		{
			Name:    "go-intro",
			Title:   "Introduction to Go",
			TitleKo: "Go 소개",
			Brief:   "A brief intro to the Go programming language.",
			Tags:    []string{"go", "programming"},
			Aliases: []string{"golang-intro"},
		},
		{
			Name:    "rust-intro",
			Title:   "Introduction to Rust",
			TitleKo: "Rust 소개",
			Brief:   "A brief intro to the Rust programming language.",
			Tags:    []string{"rust", "systems"},
			Aliases: []string{"rust-lang-intro"},
		},
	}

	data := BuildSearchIndex(pages)
	if data == nil {
		t.Fatal("BuildSearchIndex returned nil")
	}

	// Verify valid JSON
	var entries []SearchEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Verify 2 entries
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Verify correct fields for first entry
	e := entries[0]
	if e.Name != "go-intro" {
		t.Errorf("expected name=go-intro, got %s", e.Name)
	}
	if e.Title != "Introduction to Go" {
		t.Errorf("expected title=Introduction to Go, got %s", e.Title)
	}
	if e.TitleKo != "Go 소개" {
		t.Errorf("expected titleKo=Go 소개, got %s", e.TitleKo)
	}
	if e.Brief != "A brief intro to the Go programming language." {
		t.Errorf("unexpected brief: %s", e.Brief)
	}
	if len(e.Tags) != 2 || e.Tags[0] != "go" {
		t.Errorf("unexpected tags: %v", e.Tags)
	}
	if len(e.Aliases) != 1 || e.Aliases[0] != "golang-intro" {
		t.Errorf("unexpected aliases: %v", e.Aliases)
	}
}
