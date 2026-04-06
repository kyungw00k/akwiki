---
type: spec
workdir: .tmp-test
timeout: 30000
---

# Build

akwiki 빌드 파이프라인은 마크다운 파일을 정적 HTML 사이트로 변환한다.

## setup

테스트용 위키를 준비한다.

```run:shell
rm -rf build-test && mkdir -p build-test/pages build-test/public build-test/.akwiki
cat > build-test/.akwiki/config.yml << 'EOF'
site:
  title: "테스트 위키"
  author: "tester"
  language: "ko"
EOF
```

## 페이지 렌더링

### 기본 페이지

마크다운 파일이 HTML 페이지로 변환된다.

```run:shell
cat > build-test/pages/Home.md << 'EOF'
---
title: Home
titleKo: 위키 홈
---

# Home

> 환영합니다 :)

[[About]] 페이지를 참고하세요.
EOF

cat > build-test/pages/About.md << 'EOF'
---
title: About
titleKo: 소개
type: Article
---

# About

이 위키는 [[Home|홈]]에서 시작합니다.
EOF
```

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
cd build-test && git init -q && git add . && git commit -m "init" -q && $AKWIKI build 2>&1
```

각 페이지는 `/pages/{name}/index.html` 경로에 생성된다.

```run:shell
test -f build-test/dist/pages/Home/index.html
test -f build-test/dist/pages/About/index.html
```

### Private 페이지 제외

`private: true`인 페이지는 빌드에서 완전히 제외된다.

```run:shell
cat > build-test/pages/Secret.md << 'EOF'
---
title: Secret
private: true
---

비밀 내용.
EOF
```

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
cd build-test && git add . && git commit -m "add secret" -q && $AKWIKI build 2>&1
```

```run:shell
test ! -f build-test/dist/pages/Secret/index.html
```

## 위키링크

### 내부 링크

`[[페이지명]]` 문법이 HTML 링크로 변환된다.

```run:shell
grep 'class="wikilink"' build-test/dist/pages/Home/index.html
```

$ ...wikilink...

### 표시 텍스트

`[[페이지명|표시 텍스트]]` 문법으로 표시 텍스트를 지정할 수 있다.

```run:shell
grep '홈' build-test/dist/pages/About/index.html
```

$ ...홈...

## 원본 마크다운

각 페이지의 마크다운 원본이 `.txt` 확장자로 제공된다.

```run:shell
test -f build-test/dist/pages/Home.txt
test -f build-test/dist/pages/About.txt
```

## 검색 인덱스

빌드 시 클라이언트 검색용 JSON 인덱스가 생성된다.

```run:shell
test -f build-test/dist/search-index.json
```

인덱스에는 각 페이지의 제목이 포함된다.

```run:shell
cat build-test/dist/search-index.json | grep '"title"'
```

$ ...Home...
$ ...About...

## 홈 리다이렉트

루트 index.html은 /pages/Home으로 리다이렉트한다.

```run:shell
grep 'Home' build-test/dist/index.html
```

$ ...Home...

## teardown

```run:shell
rm -rf build-test
```
