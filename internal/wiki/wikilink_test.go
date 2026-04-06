package wiki

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

func TestWikilinkParser(t *testing.T) {
	md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtension("/pages")))

	tests := []struct {
		name     string
		input    string
		contains []string
		absent   []string
	}{
		{
			name:  "basic wikilink",
			input: "See [[Hello World]] for details.",
			contains: []string{
				`<a href="/pages/Hello%20World" class="wikilink">Hello World</a>`,
			},
		},
		{
			name:  "wikilink with display text",
			input: "Check [[Hello World|the greeting]] page.",
			contains: []string{
				`<a href="/pages/Hello%20World" class="wikilink">the greeting</a>`,
			},
		},
		{
			name:  "multiple wikilinks",
			input: "See [[Foo]] and [[Bar]].",
			contains: []string{
				`<a href="/pages/Foo" class="wikilink">Foo</a>`,
				`<a href="/pages/Bar" class="wikilink">Bar</a>`,
			},
		},
		{
			name:   "no wikilinks",
			input:  "Normal text without links.",
			absent: []string{`class="wikilink"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.input), &buf); err != nil {
				t.Fatalf("Convert failed: %v", err)
			}
			output := buf.String()

			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("expected output to contain %q\ngot: %s", want, output)
				}
			}
			for _, notWant := range tt.absent {
				if strings.Contains(output, notWant) {
					t.Errorf("expected output NOT to contain %q\ngot: %s", notWant, output)
				}
			}
		})
	}
}
