package content

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := []byte(`---
title: "Hello World"
titleKo: "안녕 세계"
type: Book
private: true
aliases:
  - hello
  - world
tags:
  - go
  - wiki
---
# Body

This is the body content.
`)

	fm, body, err := ParseFrontmatter(input)
	if err != nil {
		t.Fatalf("ParseFrontmatter() error = %v", err)
	}

	if fm.Title != "Hello World" {
		t.Errorf("Title = %q, want %q", fm.Title, "Hello World")
	}
	if fm.TitleKo != "안녕 세계" {
		t.Errorf("TitleKo = %q, want %q", fm.TitleKo, "안녕 세계")
	}
	if fm.Type != "Book" {
		t.Errorf("Type = %q, want %q", fm.Type, "Book")
	}
	if !fm.Private {
		t.Error("Private = false, want true")
	}
	if len(fm.Aliases) != 2 || fm.Aliases[0] != "hello" || fm.Aliases[1] != "world" {
		t.Errorf("Aliases = %v, want [hello world]", fm.Aliases)
	}
	if len(fm.Tags) != 2 || fm.Tags[0] != "go" || fm.Tags[1] != "wiki" {
		t.Errorf("Tags = %v, want [go wiki]", fm.Tags)
	}

	bodyStr := string(body)
	if bodyStr == "" {
		t.Error("body is empty, want non-empty")
	}
	// Body should not contain the frontmatter delimiters
	if len(body) >= 3 && string(body[:3]) == "---" {
		t.Error("body starts with ---, frontmatter not stripped")
	}
}

func TestParseFrontmatterDefaults(t *testing.T) {
	input := []byte(`---
---
`)

	fm, _, err := ParseFrontmatter(input)
	if err != nil {
		t.Fatalf("ParseFrontmatter() error = %v", err)
	}

	if fm.Type != "Article" {
		t.Errorf("Type = %q, want %q", fm.Type, "Article")
	}
	if fm.Private != false {
		t.Error("Private = true, want false")
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	input := []byte(`# Just a heading

Some content here.
`)

	fm, body, err := ParseFrontmatter(input)
	if err != nil {
		t.Fatalf("ParseFrontmatter() error = %v", err)
	}

	if fm.Type != "Article" {
		t.Errorf("Type = %q, want %q", fm.Type, "Article")
	}

	// body should be the entire input
	if string(body) != string(input) {
		t.Errorf("body = %q, want entire input %q", string(body), string(input))
	}
}
