# akwiki

개인 위키 정적 사이트 생성기.

마크다운 파일만으로 [akngs](https://github.com/akngs)의 [위키](https://wiki.g15e.com/pages/Home)와 동일한 형태의 개인 위키를 구축합니다.

## 특징

- **Go 단일 바이너리** — 런타임 의존성 제로, supply chain 걱정 없음
- **위키링크** — `[[페이지명]]` 문법으로 페이지 간 연결, 백링크 자동 추적
- **관련 콘텐츠** — TF-IDF 유사도 기반 관련 글/책/논문 자동 추천
- **검색** — 클라이언트 사이드 퍼지 검색 (Ctrl+K)
- **OKLCh 색상 시스템** — 라이트/다크 모드 자동 전환
- **3컬럼 레이아웃** — 데스크톱에서 TOC | 콘텐츠 | 관련 글
- **테마 커스터마이징** — config.yml 한 줄로 색상/폰트 변경, 템플릿 오버라이드도 가능
- **GitHub Pages** — push 한 번으로 자동 배포
- **private 페이지** — `private: true`로 비공개 설정, 빌드에서 완전 제외
- **Schema.org JSON-LD** — SEO 메타데이터 자동 생성

## 설치

### Homebrew

```bash
brew install kyungw00k/cli/akwiki
```

### Go

```bash
go install github.com/kyungw00k/akwiki/cmd/akwiki@latest
```

### Shell script

```bash
curl -sSfL https://kyungw00k.github.io/akwiki/install.sh | sh
```

### 바이너리 직접 다운로드

[Releases](https://github.com/kyungw00k/akwiki/releases) 페이지에서 OS/Arch에 맞는 바이너리를 다운로드하세요.

## Quick Start

```bash
# 위키 초기화
akwiki init my-wiki
cd my-wiki

# 마크다운 작성
cat > pages/About.md << 'EOF'
---
title: About
titleKo: 소개
---

# About

이 위키는 [[Home|홈]]에서 시작합니다.
EOF

# 로컬 미리보기
akwiki dev

# 정적 사이트 빌드
akwiki build
```

## 명령어

| 명령 | 설명 |
|------|------|
| `akwiki init [dir]` | 새 위키 생성 (디렉토리 구조 + 샘플 + GitHub Actions) |
| `akwiki build` | 정적 사이트 생성 → `dist/` |
| `akwiki dev` | 로컬 개발 서버 (파일 변경 시 자동 리빌드) |
| `akwiki serve` | `dist/` 미리보기 서버 |

## 콘텐츠 모델

`pages/` 디렉토리에 마크다운 파일을 추가하면 위키 페이지가 됩니다.

### Frontmatter

```yaml
---
title: specdown
titleKo: 스펙다운
type: Article        # Article | Book | ScholarlyArticle | Journal
private: false       # true면 빌드에서 제외
aliases:
  - 스펙다운
tags:
  - markdown
---
```

모든 필드는 선택입니다. `title`이 없으면 파일명이 제목이 됩니다.

`createdAt`과 `modifiedAt`는 git 이력에서 자동 추출됩니다.

### 위키링크

```markdown
[[페이지 이름]]              → 해당 페이지로 링크
[[페이지 이름|표시 텍스트]]   → 표시 텍스트로 링크
```

- `private: true` 페이지로의 링크 → 비활성 링크로 표시
- 존재하지 않는 페이지 링크 → 점선 밑줄로 표시

## 테마 커스터마이징

### 설정으로 변경 (`.akwiki/config.yml`)

```yaml
theme:
  colors:
    background: "#ffffff"
    text: "#1a1a1a"
    link: "#0969da"
  fonts:
    heading: "Noto Serif KR, serif"
    body: "system-ui, sans-serif"
  layout:
    toc: true
    backlinks: true
    related: true
    search: true
  footer:
    copyright: "2026 © 홍길동"
  edit:
    url: "obsidian://open?vault=wiki&file={{pagename}}"
```

### 템플릿 오버라이드

`.akwiki/theme/templates/partials/` 에 파일을 두면 해당 부분만 교체됩니다:

```
.akwiki/theme/
├── templates/
│   └── partials/
│       └── footer.html    ← 푸터만 커스텀
└── static/
    └── style.css          ← CSS 추가
```

## GitHub Pages 배포

`akwiki init`이 `.github/workflows/deploy.yml`을 자동 생성합니다.

```bash
git push origin main
# → GitHub Actions가 자동으로 빌드 + 배포
```

Repository Settings → Pages → Source를 "GitHub Actions"로 설정하세요.

## 설정 파일 전체 스키마

`.akwiki/config.yml`:

```yaml
site:
  title: "나의 위키"       # 사이트 제목
  author: "홍길동"         # 저자
  url: ""                  # 배포 URL (GitHub Pages 등)
  language: "ko"           # 기본 언어

build:
  outDir: "dist"           # 출력 디렉토리
  pageRoute: "/pages"      # URL 경로 접두사

analytics:
  ga: ""                   # Google Analytics ID (선택)

theme:
  colors: {}               # 빈 값 = 기본 테마 색상
  fonts: {}
  layout:
    toc: true
    backlinks: true
    related: true
    search: true
  footer:
    copyright: ""
    links: []
  edit:
    url: ""
```

## 스펙

[specdown](https://github.com/corca-ai/specdown)으로 작성된 실행 가능한 사양이 포함되어 있습니다.

```bash
# specdown 설치
brew install corca-ai/tap/specdown

# 스펙 실행 (50 cases)
make build && specdown run
```

| 스펙 | 내용 |
|------|------|
| [CLI](specs/cli.spec.md) | init, build, version 명령어 |
| [Build](specs/build.spec.md) | 빌드 파이프라인, 위키링크, private 페이지 |
| [Content](specs/content.spec.md) | frontmatter, 날짜 추출 |
| [Theme](specs/theme.spec.md) | OKLCh CSS, 다크 모드, Grid 레이아웃, 검색 |

## 크레딧

이 프로젝트는 [akngs](https://github.com/akngs)의 [위키](https://wiki.g15e.com/pages/Home)가 너무 마음에 들어서 만들었습니다. 누구나 같은 형태의 개인 위키를 쉽게 구축할 수 있도록 하는 것이 목표입니다.

## 라이선스

MIT
