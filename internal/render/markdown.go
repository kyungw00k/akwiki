package render

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/kyungw00k/akwiki/internal/wiki"
)

// RenderMarkdown converts markdown source to HTML using goldmark with
// GFM, wikilink support, automatic heading IDs, and unsafe HTML enabled.
func RenderMarkdown(source []byte, pageRoute string) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			wiki.NewWikilinkExtension(pageRoute),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
