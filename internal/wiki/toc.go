package wiki

import (
	"crypto/md5"
	"fmt"
	"regexp"
)

// Heading은 마크다운 문서에서 추출한 헤딩 정보를 나타냅니다.
type Heading struct {
	Level int
	Text  string
	ID    string
}

var headingRe = regexp.MustCompile(`(?m)^(#{2,6})\s+(.+)$`)

// ExtractTOC는 마크다운에서 h2 이상의 헤딩을 추출합니다.
// h1(#)은 제외됩니다.
func ExtractTOC(markdown []byte) []Heading {
	matches := headingRe.FindAllSubmatch(markdown, -1)
	headings := make([]Heading, 0, len(matches))

	for _, m := range matches {
		level := len(m[1]) // # 개수 = 레벨
		text := string(m[2])
		id := headingID(text)

		headings = append(headings, Heading{
			Level: level,
			Text:  text,
			ID:    id,
		})
	}

	return headings
}

// headingID는 텍스트의 md5 해시 앞 4바이트를 hex로 변환하여 'h' 접두사를 붙인 ID를 반환합니다.
func headingID(text string) string {
	sum := md5.Sum([]byte(text))
	return fmt.Sprintf("h%x", sum[:4])
}
