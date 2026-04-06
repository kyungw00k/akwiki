package render

import (
	"encoding/json"

	"github.com/kyungw00k/akwiki/internal/content"
)

// GenerateJSONLD produces a Schema.org JSON-LD string for the given page.
func GenerateJSONLD(page content.Page) string {
	schemaType := page.Type
	if schemaType == "" {
		schemaType = "Article"
	}

	data := map[string]interface{}{
		"@context":     "https://schema.org/",
		"@type":        schemaType,
		"name":         page.Title,
		"abstract":     page.Brief,
		"dateCreated":  page.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		"dateModified": page.ModifiedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}
