package wiki

import (
	"sort"
	"testing"
)

func TestExtractWikilinks(t *testing.T) {
	input := []byte("See [[Hello]] and [[World|the world]] here. Also [[Hello]] again.")
	got := ExtractWikilinks(input)

	// 중복 제거된 타겟만 반환: ["Hello", "World"]
	if len(got) != 2 {
		t.Fatalf("expected 2 unique targets, got %d: %v", len(got), got)
	}

	sort.Strings(got)
	if got[0] != "Hello" || got[1] != "World" {
		t.Errorf("expected [Hello World], got %v", got)
	}
}

func TestBuildLinkMap(t *testing.T) {
	pages := map[string][]byte{
		"Home":   []byte("See [[About]] and [[Blog]]."),
		"About":  []byte("Back to [[Home]]."),
		"Blog":   []byte("See [[About]]."),
		"Orphan": []byte("No links here."),
	}

	links, backlinks := BuildLinkMaps(pages)

	// Home → About, Blog (2개)
	if len(links["Home"]) != 2 {
		t.Errorf("expected Home to have 2 links, got %d: %v", len(links["Home"]), links["Home"])
	}

	// About의 백링크: Home, Blog (2개)
	if len(backlinks["About"]) != 2 {
		t.Errorf("expected About to have 2 backlinks, got %d: %v", len(backlinks["About"]), backlinks["About"])
	}

	// Home의 백링크: About (1개)
	if len(backlinks["Home"]) != 1 {
		t.Errorf("expected Home to have 1 backlink, got %d: %v", len(backlinks["Home"]), backlinks["Home"])
	}
	if backlinks["Home"][0] != "About" {
		t.Errorf("expected Home backlink to be About, got %v", backlinks["Home"])
	}

	// Orphan의 백링크: 없음
	if len(backlinks["Orphan"]) != 0 {
		t.Errorf("expected Orphan to have 0 backlinks, got %d: %v", len(backlinks["Orphan"]), backlinks["Orphan"])
	}
}
