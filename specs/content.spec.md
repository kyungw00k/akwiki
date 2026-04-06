---
type: spec
workdir: .tmp-test
timeout: 30000
---

# Content

akwiki의 콘텐츠 모델은 YAML frontmatter가 있는 마크다운 파일로 구성된다.

## setup

```run:shell
rm -rf content-test && mkdir -p content-test/pages content-test/.akwiki
cat > content-test/.akwiki/config.yml << 'EOF'
site:
  title: "Content Test"
  author: "tester"
EOF
```

## Frontmatter

### 전체 필드

```run:shell
cat > content-test/pages/FullPage.md << 'EOF'
---
title: Full Page
titleKo: 전체 페이지
type: Book
aliases:
  - 풀페이지
tags:
  - test
  - example
---

# Full Page

본문 내용입니다.
EOF
```

### 최소 필드

frontmatter 없이도 페이지가 생성된다. 파일명이 제목이 된다.

```run:shell
cat > content-test/pages/MinimalPage.md << 'EOF'
# Minimal Page

frontmatter 없는 페이지.
EOF
```

### 빌드 검증

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
cd content-test && git init -q && git add . && git commit -m "init" -q && $AKWIKI build 2>&1
```

두 페이지 모두 빌드된다.

```run:shell
test -f content-test/dist/pages/FullPage/index.html
test -f content-test/dist/pages/MinimalPage/index.html
```

## 날짜 추출

`createdAt`과 `modifiedAt`는 git 이력에서 자동 추출된다.

```run:shell
grep 'dateCreated' content-test/dist/pages/FullPage/index.html
```

$ ...dateCreated...

## teardown

```run:shell
rm -rf content-test
```
