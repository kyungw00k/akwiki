package search

import (
	"encoding/json"

	"github.com/kyungw00k/akwiki/internal/content"
)

// SearchEntry represents a single entry in the search index JSON.
type SearchEntry struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	TitleKo string   `json:"titleKo,omitempty"`
	Brief   string   `json:"brief,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
}

// BuildSearchIndex converts a slice of pages into a JSON search index.
func BuildSearchIndex(pages []content.Page) []byte {
	entries := make([]SearchEntry, 0, len(pages))
	for _, p := range pages {
		entries = append(entries, SearchEntry{
			Name:    p.Name,
			Title:   p.Title,
			TitleKo: p.TitleKo,
			Brief:   p.Brief,
			Tags:    p.Tags,
			Aliases: p.Aliases,
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return nil
	}
	return data
}
