---
type: spec
workdir: .tmp-test
timeout: 30000
---

# CLI

akwiki는 4개의 커맨드를 제공한다: `init`, `build`, `dev`, `serve`.

## setup

테스트용 바이너리 경로를 확인한다.

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
test -x "$AKWIKI"
```

## version

버전 정보를 출력한다.

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
$AKWIKI --version
```

$ akwiki version ...

## init

새 위키 디렉토리를 초기화한다.

### 기본 구조 생성

빈 디렉토리에 위키를 초기화하면 표준 디렉토리 구조가 생성된다.

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
rm -rf wiki-test && mkdir wiki-test
$AKWIKI init wiki-test
```

### 필수 파일 확인

초기화된 위키에는 다음 파일이 존재해야 한다.

```run:shell
test -f wiki-test/pages/Home.md
test -f wiki-test/.akwiki/config.yml
test -f wiki-test/.github/workflows/deploy.yml
```

### 필수 디렉토리 확인

```run:shell
test -d wiki-test/pages
test -d wiki-test/public
test -d wiki-test/.akwiki
```

### 홈 페이지 내용

Home.md에는 위키 안내가 포함되어야 한다.

```run:shell
grep '위키' wiki-test/pages/Home.md
```

$ ...위키...

### 이미 존재하는 파일은 건너뛰기

같은 디렉토리에 다시 init을 실행하면 기존 파일을 건너뛴다.

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
$AKWIKI init wiki-test 2>&1 | grep -c -E 'skip|건너뜀'
```

$ 3

## build

정적 사이트를 생성한다.

### 빌드 준비

빌드에는 git 이력이 필요하다 (날짜 추출용).

```run:shell
cd wiki-test && git init -q && git add . && git commit -m "init" -q
```

### 빌드 실행

```run:shell
AKWIKI="$(cd .. && pwd)/build/akwiki"
cd wiki-test && $AKWIKI build
```

### 출력 파일 확인

빌드 결과물에는 HTML 페이지, 원본 마크다운, 검색 인덱스, CSS, JS가 포함된다.

```run:shell
test -f wiki-test/dist/index.html
test -f wiki-test/dist/pages/Home/index.html
test -f wiki-test/dist/pages/Home.txt
test -f wiki-test/dist/search-index.json
test -f wiki-test/dist/assets/style.css
test -f wiki-test/dist/assets/search.js
```

### 크레딧 포함

빌드된 페이지에는 akngs 위키에서 영감을 받았다는 크레딧이 포함된다.

```run:shell
grep 'akngs' wiki-test/dist/pages/Home/index.html
```

$ ...akngs...

## teardown

```run:shell
rm -rf wiki-test
```
