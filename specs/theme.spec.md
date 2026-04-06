---
type: spec
workdir: .tmp-test
timeout: 30000
---

# Theme

akwiki 기본 테마는 wiki.g15e.com의 디자인을 재현한다.

## setup

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
rm -rf theme-test && mkdir -p theme-test/pages theme-test/.akwiki
cat > theme-test/.akwiki/config.yml << 'EOF'
site:
  title: "Theme Test"
  author: "tester"
theme:
  layout:
    toc: true
    search: true
  footer:
    copyright: "2026 © tester"
EOF
cat > theme-test/pages/Home.md << 'EOF'
---
title: Home
titleKo: 위키 홈
---

# Home

> 환영합니다

## 첫 번째 섹션

내용.

## 두 번째 섹션

내용.
EOF
cd theme-test && git init -q && git add . && git commit -m "init" -q && $AKWIKI build 2>&1
```

## 레이아웃

### 기본 구조

페이지에는 skip link, 네비게이션, 헤더, TOC, 콘텐츠, 푸터가 포함된다.

```run:shell
grep '본문으로 건너뛰기' theme-test/dist/pages/Home/index.html
```

$ ...본문으로 건너뛰기...

```run:shell
grep '위키 홈' theme-test/dist/pages/Home/index.html | head -1
```

$ ...위키 홈...

### 목차 (TOC)

h2 이상의 헤딩에서 자동 생성되는 목차가 포함된다.

```run:shell
grep '목차' theme-test/dist/pages/Home/index.html
```

$ ...목차...

```run:shell
grep '맨 위로' theme-test/dist/pages/Home/index.html
```

$ ...맨 위로...

### 푸터

설정의 copyright이 푸터에 표시된다.

```run:shell
grep '2026 © tester' theme-test/dist/pages/Home/index.html
```

$ ...2026 © tester...

크레딧이 항상 포함된다.

```run:shell
grep 'akngs' theme-test/dist/pages/Home/index.html
```

$ ...akngs...

## CSS

### OKLCh 색상 시스템

기본 CSS에는 OKLCh 기반 색상 변수가 정의되어 있다.

```run:shell
grep 'oklch' theme-test/dist/assets/style.css | head -1
```

$ ...oklch...

### 다크 모드

`prefers-color-scheme: dark` 미디어 쿼리로 다크 모드를 지원한다.

```run:shell
grep 'prefers-color-scheme' theme-test/dist/assets/style.css
```

$ ...prefers-color-scheme...

### 3컬럼 Grid 레이아웃

데스크톱에서는 `2fr 5fr 2fr` 3컬럼 레이아웃을 사용한다.

```run:shell
grep '2fr 5fr 2fr' theme-test/dist/assets/style.css
```

$ ...2fr 5fr 2fr...

### 한국어 줄바꿈

`word-break: keep-all`로 한국어 줄바꿈을 최적화한다.

```run:shell
grep 'keep-all' theme-test/dist/assets/style.css
```

$ ...keep-all...

### 구분선 (HR)

HR은 CSS content 속성으로 section mark 텍스트를 표시한다.

```run:shell
grep 'hr::after' theme-test/dist/assets/style.css
```

$ ...hr::after...

## 검색

### 검색 JS

클라이언트 사이드 검색이 포함된다.

```run:shell
grep 'search-index.json' theme-test/dist/assets/search.js
```

$ ...search-index.json...

### 키보드 단축키

Ctrl+K/Cmd+K로 검색을 토글한다.

```run:shell
grep 'ctrlKey' theme-test/dist/assets/search.js
```

$ ...ctrlKey...

## Schema.org JSON-LD

각 페이지에 Schema.org 메타데이터가 포함된다.

```run:shell
grep 'schema.org' theme-test/dist/pages/Home/index.html
```

$ ...schema.org...

## Noto Serif KR 폰트

Google Fonts에서 Noto Serif KR을 로드한다.

```run:shell
grep 'Noto.Serif.KR' theme-test/dist/pages/Home/index.html
```

$ ...Noto...Serif...KR...

## teardown

```run:shell
rm -rf theme-test
```
