# akwiki — 개인 위키 정적 사이트 생성기

> wiki.g15e.com에서 영감을 받은 개인 위키 시스템.
> Go 단일 바이너리로 배포. 마크다운 파일만으로 위키를 구축.

## 1. 개요

### 목표

- wiki.g15e.com의 기능과 디자인을 충실히 재현
- 누구나 `akwiki init` → `akwiki build`로 동일한 형태의 위키를 구축
- Go 단일 바이너리, 런타임 의존성 제로 (supply chain 공격 방지)
- GitHub Pages 배포 기본 지원

### 비목표

- CMS나 관리자 UI
- 서버 사이드 렌더링
- 멀티 유저/권한 시스템

## 2. CLI

```
akwiki init [dir]    새 위키 생성 (디렉토리 구조 + 샘플 + GitHub Actions)
akwiki dev           로컬 개발 서버 (파일 워치 + 라이브 리로드)
akwiki build         정적 사이트 생성 → dist/
akwiki serve         dist/ 미리보기 서버
```

### init이 생성하는 구조

```
my-wiki/
├── pages/
│   └── Home.md              샘플 홈 페이지
├── public/                  정적 자산 (이미지 등)
├── .akwiki/
│   └── config.yml           사이트 설정
└── .github/
    └── workflows/
        └── deploy.yml       GitHub Pages 자동 배포
```

## 3. 콘텐츠 모델

### 마크다운 frontmatter

```yaml
---
title: specdown
titleKo: 스펙다운
type: Article           # Article | Book | ScholarlyArticle | Journal
private: false          # true → 빌드에서 제외
aliases:
  - 스펙다운
tags:
  - markdown
---
```

| 필드 | 필수 | 기본값 | 설명 |
|------|------|--------|------|
| title | N | 파일명 | 영문 제목 |
| titleKo | N | - | 한국어 제목 |
| type | N | Article | 콘텐츠 타입 |
| private | N | false | 공개 여부 |
| aliases | N | [] | 위키링크에서 이 이름으로도 참조 가능 |
| tags | N | [] | 분류 태그 |

- `createdAt` / `modifiedAt`는 **git log에서 자동 추출**. frontmatter에 쓰지 않음.
- git 이력이 없는 파일은 파일 시스템 시간 사용.

### 위키링크 문법

```markdown
[[페이지 이름]]              → <a href="/pages/페이지-이름">페이지 이름</a>
[[페이지 이름|표시 텍스트]]   → <a href="/pages/페이지-이름">표시 텍스트</a>
```

- aliases에 등록된 이름으로도 링크 가능
- `private: true` 페이지로의 링크 → `<a href="#private-link" data-link="private">텍스트</a>`
- 존재하지 않는 페이지 링크 → `<a class="wikilink-missing">텍스트</a>`

## 4. 빌드 파이프라인

```
pages/*.md
    │
    ├─[1] frontmatter 파싱 + private 필터링
    │
    ├─[2] git log → createdAt / modifiedAt 추출
    │
    ├─[3] 위키링크 추출 → 링크맵 → 백링크맵 생성
    │
    ├─[4] TF-IDF 관련 콘텐츠 계산 (타입별 그룹핑)
    │
    ├─[5] goldmark으로 마크다운 → HTML 변환
    │     └─ 커스텀 파서: 위키링크, 목차 추출
    │
    ├─[6] 테마 템플릿 적용
    │     └─ Go html/template 렌더링
    │     └─ JSON-LD Schema.org 메타데이터 삽입
    │
    ├─[7] 검색 인덱스 생성 (search-index.json)
    │
    └─[8] 원본 .txt 복사
    
dist/
├── index.html                 홈 리다이렉트 → /pages/Home
├── pages/
│   ├── {pagename}/index.html  각 페이지 HTML
│   └── {pagename}.txt         마크다운 원본
├── search-index.json          클라이언트 검색용
└── assets/                    CSS, JS, 폰트, 이미지
```

## 5. 테마 시스템

### 3단계 커스터마이징

| 단계 | 대상 사용자 | 방법 |
|------|-----------|------|
| 1. 기본 | "그냥 쓸래" | `akwiki init` → 끝 |
| 2. 설정 | "색상/폰트 바꿀래" | config.yml의 `theme` 섹션 |
| 3. 오버라이드 | "완전 다른 테마" | `.akwiki/theme/`에 템플릿 배치 |

### 기본 내장 테마

Go의 `embed.FS`로 바이너리에 포함. 원본 wiki.g15e.com 디자인 재현.

### config.yml 테마 설정

```yaml
theme:
  colors:
    background: "#ffffff"
    text: "#1a1a1a"
    link: "#0969da"
    link-private: "#999999"
    accent: "#7b2cb5"
  fonts:
    heading: "Noto Serif KR, serif"
    body: "system-ui, sans-serif"
    code: "monospace"
  layout:
    max-width: "50rem"
    toc: true
    backlinks: true
    related: true
    search: true
  footer:
    copyright: "2026 © ak"
    links:
      - label: "Twitter"
        url: "https://twitter.com/@_a6g_"
  edit:
    url: "obsidian://open?vault=wiki&file={{pagename}}"
```

설정값은 CSS 변수로 주입. 템플릿 수정 불필요.

### 템플릿 오버라이드

```
.akwiki/theme/
├── templates/
│   ├── page.html           페이지 레이아웃
│   ├── home.html           홈 전용 레이아웃 (선택)
│   └── partials/
│       ├── header.html
│       ├── toc.html
│       ├── backlinks.html
│       ├── related.html
│       ├── search.html
│       └── footer.html
└── static/
    ├── style.css           추가 CSS
    └── search.js           검색 JS 오버라이드
```

파일이 존재하면 해당 부분만 오버라이드. 없으면 내장 기본 사용.

### 템플릿 데이터 컨텍스트

```go
type TemplateContext struct {
    Site      SiteConfig
    Page      Page
    Content   template.HTML     // 렌더링된 HTML 본문
    TOC       []Heading         // 목차 항목
    Links     []PageRef         // 이 페이지의 외부 참조
    Backlinks []PageRef         // 이 페이지를 참조하는 역링크
    Related   map[string][]PageRef  // 타입별 관련 콘텐츠
}

type Page struct {
    Name       string
    Title      string
    TitleKo    string
    Type       string           // Article, Book, ...
    Brief      string           // 첫 문단 요약
    Aliases    []string
    Tags       []string
    CreatedAt  time.Time
    ModifiedAt time.Time
    RawURL     string           // .txt 원본 URL
}

type PageRef struct {
    Name    string
    Title   string
    Brief   string
    Type    string
    Score   float64             // 관련 콘텐츠 유사도 (0-1)
}
```

## 6. 기본 테마 디자인 스펙 (원본 재현)

### 색상 시스템

OKLCh 색공간 기반. CSS 변수로 정의하여 라이트/다크 모드 자동 전환.

```
Primary hue:     90  (녹색 계열)
Complement hue: 270  (보라 — 액센트)
Link hue:       210  (청색 — 내부/외부 링크)
```

**라이트 모드:**

| 용도 | 변수 | 대략적 색상 |
|------|------|-----------|
| 배경 | `--c-bg` | 거의 흰색 (L:97%, 극미한 녹색 틴트) |
| 코드 배경 | `--c-bg-code` | 연한 녹회색 (L:92%) |
| 하이라이트 배경 | `--c-bg-highlight` | 연한 녹색 (L:85%) |
| 본문 텍스트 | `--c-text` | 짙은 회색 (L:30%) |
| 뮤트 텍스트 | `--c-text-muted` | 중간 회색 (L:45%) |
| 내부 링크 | `--c-text-link` | 짙은 청색 (L:25%, C:0.12) |
| 외부 링크 | `--c-text-link-ext` | 선명한 청색 (L:25%, C:0.18) |
| 액센트 | `--c-accent` | 선명한 보라 (L:50%, C:0.28) |

**다크 모드:** `--l-offset: 100`, `--l-sign: -1`로 밝기 반전. Stevens/Hunt 지수 기반 채도 보정.

### 타이포그래피

| 용도 | 폰트 | 크기 |
|------|------|------|
| 제목 (h1) | Noto Serif KR, serif | `max(1.8rem, 2.5dvw)` — 반응형 |
| 제목 (h2-h6) | Noto Serif KR, serif | 1.65rem ~ 1.2em |
| 본문 | 시스템 산세리프 | 기본 |
| 인용문 | Noto Serif KR, serif | 기본 |
| 코드 | monospace | 기본 |

줄간격: 본문 `1.6`, 제목 `1.3`, h1 `1.4`.

`word-break: keep-all` (한국어 줄바꿈 최적화).

### 레이아웃

**모바일 (≤72rem):** 단일 컬럼, max-width `50rem`, 세로 스택.

```
nav → header → toc → content → backlinks → footer
```

**데스크톱 (>72rem):** CSS Grid 3컬럼 `2fr 5fr 2fr`, max-width `100rem`.

```
.     nav       .
.     header    .
toc   content   backlinks
.     footer    .
```

TOC는 `position: sticky; top: 0`으로 스크롤 시 고정.

### UI 디테일

| 요소 | 구현 |
|------|------|
| 링크 밑줄 | 두께 `0.5px`, offset `0.25em` |
| 외부 링크 | 뒤에 ` ↗` (font-size `0.5em`, superscript) |
| Private 링크 | dashed 밑줄, `--c-text-muted` 색상 |
| 강조(bold) | weight 유지, 대신 `--c-text-highlight` + `--c-bg-highlight` |
| 인라인 코드 | `--c-bg-code` 배경, border-radius `0.25em` |
| 코드 블록 | `--c-bg-code` 배경, padding `1em 1.5em` |
| 인용문 | 좌측 border `0.25em solid --c-bg-highlight`, 세리프 폰트 |
| HR (구분선) | 텍스트 `"- - - § - - -"`, 색상 `--c-text-muted` |
| 테이블 | collapse, 셀 padding `0.25em 1em`, border `1px solid --c-text-muted` |
| 프로그레스 바 | `--c-accent` 보라, 높이 `3px`, cubic-bezier 애니메이션 |
| Skip link | "본문으로 건너뛰기" |
| 맨 위로 | 목차 상단 `(맨 위로)` 앵커 |

### 페이지 구조 (위→아래)

```html
<a href="#content">본문으로 건너뛰기</a>

<nav>
  <a href="/pages/Home">위키 홈</a>
  <a href="obsidian://...">Edit</a>
</nav>

<header>
  <h1>{{ .Page.Title }}</h1>
  <time>{{ .Page.CreatedAt }} (modified: {{ .Page.ModifiedAt }})</time>
</header>

<nav aria-label="목차">
  <h2>목차</h2>
  <ol>
    <li><a href="#top">(맨 위로)</a></li>
    <li><a href="#section-id">섹션 제목</a></li>
    ...
  </ol>
</nav>

<article data-content>
  {{ .Content }}
</article>

<aside>  <!-- 관련 콘텐츠 (타입별) -->
  <section>
    <h2>관련 책</h2>
    <ul>...</ul>
  </section>
  <section>
    <h2>관련 글</h2>
    <ul>...</ul>
  </section>
  <section>
    <h2>관련 논문</h2>
    <ul>...</ul>
  </section>
</aside>

<footer>
  <p>{{ .Site.Footer.Copyright }} | <a href="{{ .Page.RawURL }}">markdown</a></p>
</footer>

<script type="application/ld+json">
  { "@context": "https://schema.org/", "@type": "Article", ... }
</script>
```

## 7. 검색

- 빌드 시 `search-index.json` 생성 (각 페이지의 title, titleKo, brief, tags, aliases)
- 클라이언트에서 경량 JS로 퍼지 검색 (인라인 스크립트, 외부 의존성 없음)
- 검색 UI는 테마의 `search.html` partial

## 8. 관련 콘텐츠

- 빌드 타임에 TF-IDF 텍스트 유사도 계산 (Go 자체 구현)
- 타입별 그룹핑: 관련 책 (Book), 관련 글 (Article), 관련 논문 (ScholarlyArticle)
- 유사도 상위 N개 (기본 10) 노출
- 점수 기반 정렬

## 9. GitHub Pages 배포

`akwiki init`이 생성하는 `.github/workflows/deploy.yml`:

```yaml
name: Deploy wiki
on:
  push:
    branches: [main]
permissions:
  contents: read
  pages: write
  id-token: write
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0          # git log 전체 (날짜 추출)
      - name: Install akwiki
        run: |
          curl -fsSL https://github.com/kyungw00k/akwiki/releases/latest/download/akwiki-linux-amd64 -o akwiki
          chmod +x akwiki
      - name: Build
        run: ./akwiki build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: dist
  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
```

## 10. 기술 스택

| 영역 | 선택 | 이유 |
|------|------|------|
| 언어 | Go | 단일 바이너리, 의존성 제로 |
| CLI | cobra | Go CLI 표준 |
| 마크다운 | goldmark | Hugo도 사용, 확장성 |
| 위키링크 | goldmark 커스텀 파서 | 직접 구현 |
| 템플릿 | html/template | 표준 라이브러리 |
| 파일 워치 | fsnotify | dev 서버 핫 리로드 |
| HTTP 서버 | net/http | 표준 라이브러리 |
| 유사도 | TF-IDF 자체 구현 | 알고리즘 단순, 외부 의존 불필요 |
| 검색 (클라이언트) | 인라인 JS | 외부 라이브러리 없음 |
| 내장 테마 | embed.FS | 바이너리에 포함 |
| git 날짜 | os/exec → git log | 외부 의존 없음 |

## 11. 설정 파일 전체 스키마

`.akwiki/config.yml`:

```yaml
site:
  title: "나의 위키"          # 필수
  author: "홍길동"            # 필수
  url: ""                     # GitHub Pages URL (빌드 시 base path)
  language: "ko"              # 기본 언어

build:
  outDir: "dist"
  pageRoute: "/pages"

analytics:
  ga: ""                      # Google Analytics ID (선택)

theme:
  colors:
    background: ""            # 빈 값 = 기본 테마 색상
    text: ""
    link: ""
    link-private: ""
    accent: ""
  fonts:
    heading: ""
    body: ""
    code: ""
  layout:
    max-width: ""
    toc: true
    backlinks: true
    related: true
    search: true
  footer:
    copyright: ""
    links: []
  edit:
    url: ""                   # {{pagename}} 플레이스홀더 지원
```

모든 값은 선택. 빈 값이면 원본 wiki.g15e.com 기본 테마 적용.

## 12. 크레딧

akwiki가 생성하는 모든 위키의 기본 테마 푸터에 다음 문구를 포함한다:

> Inspired by [akngs](https://github.com/akngs)'s [wiki](https://wiki.g15e.com/pages/Home).
> Powered by [akwiki](https://github.com/kyungw00k/akwiki).

이 프로젝트는 akngs의 위키 시스템이 너무 마음에 들어서 만들었다. 누구나 같은 형태의 개인 위키를 쉽게 구축할 수 있도록 하는 것이 목표다.
