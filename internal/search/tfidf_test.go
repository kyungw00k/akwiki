package search

import "testing"

func TestTFIDF(t *testing.T) {
	docs := map[string]string{
		"go-intro":     "Go is a programming language designed for simplicity and efficiency. Go programs are compiled and statically typed.",
		"rust-intro":   "Rust is a systems programming language focused on safety and performance. Rust prevents memory errors.",
		"go-advanced":  "Advanced Go programming covers concurrency goroutines channels and Go runtime internals for building scalable programs.",
		"cooking":      "Pasta is made from durum wheat flour and water. Italian cooking uses many types of pasta with various sauces.",
	}

	idx := NewTFIDFIndex(docs)

	// go-intro should be most similar to go-advanced
	similar := idx.MostSimilar("go-intro", 3)
	if len(similar) == 0 {
		t.Fatal("expected at least one similar document")
	}
	if similar[0].Name != "go-advanced" {
		t.Errorf("expected go-advanced as most similar to go-intro, got %s", similar[0].Name)
	}

	// cooking should not be highly similar to go-intro (score < 0.5)
	cookingScore := 0.0
	for _, s := range similar {
		if s.Name == "cooking" {
			cookingScore = s.Score
			break
		}
	}
	if cookingScore >= 0.5 {
		t.Errorf("expected cooking score < 0.5 relative to go-intro, got %f", cookingScore)
	}
}

func TestTFIDFEmpty(t *testing.T) {
	docs := map[string]string{
		"only-doc": "This is the only document in the index.",
	}

	idx := NewTFIDFIndex(docs)
	similar := idx.MostSimilar("only-doc", 5)
	if len(similar) != 0 {
		t.Errorf("expected empty result for single-doc index, got %d results", len(similar))
	}
}
