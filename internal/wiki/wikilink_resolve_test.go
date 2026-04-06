package wiki

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
)

func TestWikilinkPrivate(t *testing.T) {
	resolver := func(target string) LinkStatus {
		if target == "Secret" {
			return LinkPrivate
		}
		return LinkExists
	}
	md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtensionWithResolver("/pages", resolver)))

	var buf bytes.Buffer
	if err := md.Convert([]byte("See [[Secret]] page."), &buf); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `href="#private-link"`) {
		t.Errorf("expected href=#private-link, got: %s", output)
	}
	if !strings.Contains(output, `data-link="private"`) {
		t.Errorf("expected data-link=private, got: %s", output)
	}
	if !strings.Contains(output, `class="wikilink"`) {
		t.Errorf("expected class=wikilink, got: %s", output)
	}
}

func TestWikilinkMissing(t *testing.T) {
	resolver := func(target string) LinkStatus {
		if target == "NonExistent" {
			return LinkMissing
		}
		return LinkExists
	}
	md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtensionWithResolver("/pages", resolver)))

	var buf bytes.Buffer
	if err := md.Convert([]byte("See [[NonExistent]]."), &buf); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `class="wikilink-missing"`) {
		t.Errorf("expected class=wikilink-missing, got: %s", output)
	}
	if strings.Contains(output, `href=`) {
		t.Errorf("expected no href for missing link, got: %s", output)
	}
}

func TestWikilinkResolverExists(t *testing.T) {
	resolver := func(target string) LinkStatus {
		return LinkExists
	}
	md := goldmark.New(goldmark.WithExtensions(NewWikilinkExtensionWithResolver("/pages", resolver)))

	var buf bytes.Buffer
	if err := md.Convert([]byte("See [[Home]]."), &buf); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `href="/pages/Home"`) {
		t.Errorf("expected normal href, got: %s", output)
	}
	if !strings.Contains(output, `class="wikilink"`) {
		t.Errorf("expected class=wikilink, got: %s", output)
	}
}
