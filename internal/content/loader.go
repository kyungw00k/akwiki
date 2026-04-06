package content

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// LoadPages reads all .md files from pagesDir and returns a slice of Pages.
func LoadPages(pagesDir string) ([]Page, error) {
	entries, err := os.ReadDir(pagesDir)
	if err != nil {
		return nil, err
	}

	repoDir := filepath.Dir(pagesDir)

	var pages []Page
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(pagesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		name := strings.TrimSuffix(entry.Name(), ".md")

		fm, body, err := ParseFrontmatter(data)
		if err != nil {
			return nil, err
		}

		title := fm.Title
		if title == "" {
			title = name
		}

		relPath := filepath.Join(filepath.Base(pagesDir), entry.Name())
		created, modified, _ := GitDates(repoDir, relPath)

		brief := extractBrief(body)

		page := Page{
			Name:       name,
			Title:      title,
			TitleKo:    fm.TitleKo,
			Type:       fm.Type,
			Brief:      brief,
			Private:    fm.Private,
			Aliases:    fm.Aliases,
			Tags:       fm.Tags,
			CreatedAt:  created,
			ModifiedAt: modified,
			RawBody:    body,
			RawSource:  data,
		}

		pages = append(pages, page)
	}

	return pages, nil
}

// extractBrief returns the first non-heading, non-empty paragraph line (max 200 chars).
func extractBrief(body []byte) string {
	lines := bytes.Split(body, []byte("\n"))
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		// Skip headings
		if bytes.HasPrefix(trimmed, []byte("#")) {
			continue
		}
		// Skip horizontal rules
		if bytes.Equal(trimmed, []byte("---")) || bytes.Equal(trimmed, []byte("***")) || bytes.Equal(trimmed, []byte("___")) {
			continue
		}
		brief := string(trimmed)
		if len(brief) > 200 {
			brief = brief[:200]
		}
		return brief
	}
	return ""
}
