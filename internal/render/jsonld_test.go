package render

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kyungw00k/akwiki/internal/content"
)

func TestGenerateJSONLD(t *testing.T) {
	page := content.Page{
		Title:      "Test Page",
		Brief:      "A test page",
		CreatedAt:  time.Date(2024, 9, 8, 0, 0, 0, 0, time.UTC),
		ModifiedAt: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		Type:       "Article",
	}

	result := GenerateJSONLD(page)
	if result == "" {
		t.Fatal("GenerateJSONLD returned empty string")
	}

	// Verify valid JSON
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, result)
	}

	// Verify @context
	if m["@context"] != "https://schema.org/" {
		t.Errorf("expected @context=https://schema.org/, got %v", m["@context"])
	}

	// Verify @type
	if m["@type"] != "Article" {
		t.Errorf("expected @type=Article, got %v", m["@type"])
	}

	// Verify name
	if m["name"] != "Test Page" {
		t.Errorf("expected name=Test Page, got %v", m["name"])
	}

	// Verify dateCreated contains 2024-09-08
	if dc, ok := m["dateCreated"].(string); !ok || len(dc) < 10 || dc[:10] != "2024-09-08" {
		t.Errorf("expected dateCreated starting with 2024-09-08, got %v", m["dateCreated"])
	}

	// Verify dateModified contains 2026-03-15
	if dm, ok := m["dateModified"].(string); !ok || len(dm) < 10 || dm[:10] != "2026-03-15" {
		t.Errorf("expected dateModified starting with 2026-03-15, got %v", m["dateModified"])
	}
}
