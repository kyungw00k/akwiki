package wiki

import "regexp"

var wikilinkRe = regexp.MustCompile(`\[\[([^\]|]+)(?:\|[^\]]+)?\]\]`)

// ExtractWikilinks는 마크다운 본문에서 위키링크 타겟을 추출하여 중복 제거 후 반환합니다.
func ExtractWikilinks(body []byte) []string {
	matches := wikilinkRe.FindAllSubmatch(body, -1)

	seen := make(map[string]bool)
	var result []string

	for _, m := range matches {
		target := string(m[1])
		if !seen[target] {
			seen[target] = true
			result = append(result, target)
		}
	}

	return result
}

// BuildLinkMaps는 페이지 맵으로부터 양방향 링크 맵을 생성합니다.
// links[pageName] = 해당 페이지에서 링크하는 페이지 목록
// backlinks[pageName] = 해당 페이지를 링크하는 페이지 목록
func BuildLinkMaps(pages map[string][]byte) (links map[string][]string, backlinks map[string][]string) {
	links = make(map[string][]string)
	backlinks = make(map[string][]string)

	for name, body := range pages {
		targets := ExtractWikilinks(body)
		links[name] = targets

		for _, target := range targets {
			backlinks[target] = append(backlinks[target], name)
		}
	}

	return links, backlinks
}
