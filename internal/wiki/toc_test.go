package wiki

import (
	"testing"
)

func TestExtractTOC(t *testing.T) {
	input := []byte(`# Title

## Section One

### Subsection

## Section Two
`)

	headings := ExtractTOC(input)

	// h1 제외, h2+ 만 반환 → 3개 (Section One, Subsection, Section Two)
	if len(headings) != 3 {
		t.Fatalf("expected 3 headings (h2+), got %d: %v", len(headings), headings)
	}

	// 레벨 검증
	expected := []struct {
		level int
		text  string
	}{
		{2, "Section One"},
		{3, "Subsection"},
		{2, "Section Two"},
	}

	for i, e := range expected {
		h := headings[i]
		if h.Level != e.level {
			t.Errorf("heading[%d]: expected level %d, got %d", i, e.level, h.Level)
		}
		if h.Text != e.text {
			t.Errorf("heading[%d]: expected text %q, got %q", i, e.text, h.Text)
		}
		if h.ID == "" {
			t.Errorf("heading[%d]: expected non-empty ID", i)
		}
	}
}
