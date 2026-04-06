package render

import (
	"bytes"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	input := `# My Wiki Page

This is a paragraph with a [[Wiki Link]] and [[Another|display text]].

## Section One

- item one
- item two

> This is a blockquote

` + "```go\nfmt.Println(\"hello\")\n```"

	out, err := RenderMarkdown([]byte(input), "/pages")
	if err != nil {
		t.Fatalf("RenderMarkdown error: %v", err)
	}

	// Verify wikilinks rendered with class="wikilink" and href containing /pages/
	if !bytes.Contains(out, []byte(`class="wikilink"`)) {
		t.Errorf("expected wikilink class in output, got:\n%s", out)
	}
	if !bytes.Contains(out, []byte(`href="/pages/`)) {
		t.Errorf("expected href=/pages/ in output, got:\n%s", out)
	}

	// Verify display text for pipe-style wikilink
	if !bytes.Contains(out, []byte(`display text`)) {
		t.Errorf("expected display text in output, got:\n%s", out)
	}

	// Verify standard elements
	if !bytes.Contains(out, []byte(`<h1`)) {
		t.Errorf("expected h1 in output, got:\n%s", out)
	}
	if !bytes.Contains(out, []byte(`<blockquote`)) {
		t.Errorf("expected blockquote in output, got:\n%s", out)
	}
}

func TestRenderMarkdownHeadingIDs(t *testing.T) {
	input := `## First Heading

Some text.

## Second Heading

More text.
`
	out, err := RenderMarkdown([]byte(input), "/pages")
	if err != nil {
		t.Fatalf("RenderMarkdown error: %v", err)
	}

	if !bytes.Contains(out, []byte(`id="`)) {
		t.Errorf("expected id= attributes on headings, got:\n%s", out)
	}
}
