package content

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// Frontmatter holds the parsed YAML front matter of a page.
type Frontmatter struct {
	Title   string   `yaml:"title"`
	TitleKo string   `yaml:"titleKo"`
	Type    string   `yaml:"type"`
	Private bool     `yaml:"private"`
	Aliases []string `yaml:"aliases"`
	Tags    []string `yaml:"tags"`
}

var frontmatterDelimiter = []byte("---")

// ParseFrontmatter parses YAML front matter from the given data.
// It returns the parsed Frontmatter, the body (everything after the front matter),
// and any error encountered.
// If no front matter is found, returns defaults and the full data as body.
func ParseFrontmatter(data []byte) (*Frontmatter, []byte, error) {
	fm := &Frontmatter{
		Type: "Article",
	}

	// Check if the file starts with ---
	if !bytes.HasPrefix(data, frontmatterDelimiter) {
		return fm, data, nil
	}

	// Find the closing ---
	// Skip the first line (opening ---)
	rest := data[3:]
	// Skip optional newline after opening ---
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	// Find the closing delimiter
	closingIdx := bytes.Index(rest, frontmatterDelimiter)
	if closingIdx == -1 {
		// No closing delimiter found, treat as no frontmatter
		return fm, data, nil
	}

	yamlContent := rest[:closingIdx]
	body := rest[closingIdx+3:]
	// Skip optional newline after closing ---
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	} else if len(body) > 1 && body[0] == '\r' && body[1] == '\n' {
		body = body[2:]
	}

	if err := yaml.Unmarshal(yamlContent, fm); err != nil {
		return fm, data, err
	}

	// Apply defaults for zero values
	if fm.Type == "" {
		fm.Type = "Article"
	}

	return fm, body, nil
}
