package search

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// SimilarDoc represents a document with its similarity score.
type SimilarDoc struct {
	Name  string
	Score float64
}

// TFIDFIndex holds TF-IDF vectors for a collection of documents.
type TFIDFIndex struct {
	docs    map[string]string
	vectors map[string]map[string]float64
}

// tokenize splits text into lowercase tokens, filtering short tokens.
func tokenize(text string) []string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	tokens := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.ToLower(w)
		if len(w) > 1 {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

// NewTFIDFIndex builds a TF-IDF index from a map of document name → text.
func NewTFIDFIndex(docs map[string]string) *TFIDFIndex {
	N := float64(len(docs))

	// Compute term frequencies per document
	tfMap := make(map[string]map[string]float64, len(docs))
	for name, text := range docs {
		tokens := tokenize(text)
		counts := make(map[string]int)
		for _, t := range tokens {
			counts[t]++
		}
		total := len(tokens)
		tf := make(map[string]float64, len(counts))
		for term, cnt := range counts {
			tf[term] = float64(cnt) / float64(total)
		}
		tfMap[name] = tf
	}

	// Compute document frequency per term
	df := make(map[string]int)
	for _, tf := range tfMap {
		for term := range tf {
			df[term]++
		}
	}

	// Compute TF-IDF vectors
	vectors := make(map[string]map[string]float64, len(docs))
	for name, tf := range tfMap {
		vec := make(map[string]float64, len(tf))
		for term, tfVal := range tf {
			idf := math.Log(1 + N/float64(df[term]))
			vec[term] = tfVal * idf
		}
		vectors[name] = vec
	}

	return &TFIDFIndex{
		docs:    docs,
		vectors: vectors,
	}
}

// cosineSimilarity computes cosine similarity between two TF-IDF vectors.
func cosineSimilarity(a, b map[string]float64) float64 {
	var dot, normA, normB float64
	for term, va := range a {
		dot += va * b[term]
		normA += va * va
	}
	for _, vb := range b {
		normB += vb * vb
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// MostSimilar returns up to limit documents most similar to name, sorted by score descending.
// Only documents with score > 0.01 are included.
func (idx *TFIDFIndex) MostSimilar(name string, limit int) []SimilarDoc {
	vec, ok := idx.vectors[name]
	if !ok {
		return nil
	}

	var results []SimilarDoc
	for other, otherVec := range idx.vectors {
		if other == name {
			continue
		}
		score := cosineSimilarity(vec, otherVec)
		if score > 0.01 {
			results = append(results, SimilarDoc{Name: other, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > limit {
		results = results[:limit]
	}
	return results
}
