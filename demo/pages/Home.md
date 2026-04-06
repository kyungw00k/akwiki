---
title: Home
titleKo: akwiki demo
---

# akwiki demo

> 이 사이트는 [akwiki](https://github.com/kyungw00k/akwiki)로 만든 데모 위키입니다.

[akngs](https://github.com/akngs)의 [위키](https://wiki.g15e.com/pages/Home)가 너무 마음에 들어서, 누구나 같은 형태의 위키를 쉽게 만들 수 있도록 akwiki를 만들었습니다.

## 시작하기

akwiki는 Go 단일 바이너리로, 마크다운 파일만으로 개인 위키를 구축합니다.

```
brew install kyungw00k/cli/akwiki
akwiki init my-wiki
cd my-wiki
akwiki dev
```

자세한 설치 방법은 [[설치 가이드]]를 참고하세요.

## 기능

- **[[위키링크]]** — `[[페이지명]]` 문법으로 페이지 간 연결
- **[[관련 콘텐츠]]** — TF-IDF 유사도 기반 자동 추천
- **[[테마 시스템]]** — 3단계 커스터마이징 (기본 → 설정 → 오버라이드)
- **[[마크다운 문법]]** — frontmatter, 위키링크, 일반 마크다운

## 위키 지도

### 가이드

- [[설치 가이드]]
- [[마크다운 문법]]
- [[테마 시스템]]
- [[GitHub Pages 배포]]

### 기능 설명

- [[위키링크]]
- [[관련 콘텐츠]]
- [[검색 기능]]

### 읽을거리

- [[디지털 정원]]
- [[개인 지식 관리]]
- [[Zettelkasten]]
