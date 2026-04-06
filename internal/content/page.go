package content

import "time"

// Page represents a single wiki page.
type Page struct {
	Name       string
	Title      string
	TitleKo    string
	Type       string
	Brief      string
	Private    bool
	Aliases    []string
	Tags       []string
	CreatedAt  time.Time
	ModifiedAt time.Time
	RawBody    []byte // markdown body without frontmatter
	RawSource  []byte // original full file content
	RawURL     string // .txt URL (set during build)
}
