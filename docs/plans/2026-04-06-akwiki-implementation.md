# akwiki Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI tool that generates static wiki sites from markdown files, faithfully reproducing wiki.g15e.com's design and features.

**Architecture:** Single Go binary with embedded default theme. Core pipeline: parse markdown → extract links/metadata → compute relationships → render HTML via templates. CLI commands: init, build, dev, serve.

**Tech Stack:** Go 1.22+, cobra (CLI), goldmark (markdown), html/template (templates), fsnotify (file watch), embed.FS (theme assets)

---

## File Structure

```
akwiki/
├── main.go                          CLI entrypoint
├── go.mod
├── go.sum
├── cmd/
│   ├── root.go                      cobra root command
│   ├── init.go                      akwiki init
│   ├── build.go                     akwiki build
│   ├── dev.go                       akwiki dev
│   └── serve.go                     akwiki serve
├── internal/
│   ├── config/
│   │   └── config.go                config.yml 파싱 + 기본값
│   ├── content/
│   │   ├── frontmatter.go           YAML frontmatter 파싱
│   │   ├── page.go                  Page 타입 정의 + 로딩
│   │   └── gitdate.go              git log → 날짜 추출
│   ├── wiki/
│   │   ├── wikilink.go              goldmark 위키링크 파서 확장
│   │   ├── linkmap.go               링크맵 + 백링크맵
│   │   └── toc.go                   목차(TOC) 추출
│   ├── search/
│   │   ├── tfidf.go                 TF-IDF 유사도 계산
│   │   └── index.go                 검색 인덱스 JSON 생성
│   ├── render/
│   │   ├── markdown.go              goldmark 설정 + HTML 변환
│   │   ├── template.go              html/template 로딩 + 렌더링
│   │   └── jsonld.go                Schema.org JSON-LD 생성
│   └── builder/
│       └── builder.go               전체 빌드 파이프라인 오케스트레이션
├── theme/
│   └── default/
│       ├── templates/
│       │   ├── page.html            메인 페이지 레이아웃
│       │   └── partials/
│       │       ├── header.html
│       │       ├── toc.html
│       │       ├── backlinks.html
│       │       ├── related.html
│       │       ├── search.html
│       │       └── footer.html
│       └── static/
│           ├── style.css            원본 재현 CSS
│           └── search.js            클라이언트 검색 JS
└── scaffold/                        init 명령이 복사할 파일들
    ├── Home.md
    ├── config.yml
    └── deploy.yml
```

---

## Task 1: Go 프로젝트 초기화 + CLI 뼈대

**Files:**
- Create: `main.go`
- Create: `go.mod`
- Create: `cmd/root.go`
- Create: `cmd/build.go`
- Create: `cmd/init_cmd.go`
- Create: `cmd/dev.go`
- Create: `cmd/serve.go`

- [ ] **Step 1: Go 모듈 초기화**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
go mod init github.com/kyungw00k/akwiki
```

- [ ] **Step 2: cobra 의존성 추가**

```bash
go get github.com/spf13/cobra@latest
```

- [ ] **Step 3: main.go 작성**

```go
// main.go
package main

import "github.com/kyungw00k/akwiki/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 4: root 커맨드 작성**

```go
// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "akwiki",
	Short: "Personal wiki static site generator",
	Long:  "akwiki generates static wiki sites from markdown files.\nInspired by akngs's wiki (https://wiki.g15e.com).",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

- [ ] **Step 5: 4개 서브커맨드 스텁 작성**

```go
// cmd/init_cmd.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create a new wiki",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		fmt.Printf("Initializing wiki in %s\n", dir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
```

```go
// cmd/build.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static site",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Building wiki...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
```

```go
// cmd/dev.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with live reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting dev server...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
```

```go
// cmd/serve.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the built site",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Serving dist/...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
```

- [ ] **Step 6: 빌드 및 실행 확인**

```bash
go build -o akwiki .
./akwiki --help
./akwiki init --help
./akwiki build
```

Expected: 각 명령이 스텁 메시지 출력.

- [ ] **Step 7: 커밋**

```bash
git init
git add main.go go.mod go.sum cmd/
git commit -m "feat: initialize Go project with cobra CLI skeleton"
```

---

## Task 2: 설정 시스템 (config.yml)

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: 설정 파싱 테스트 작성**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Site.Language != "ko" {
		t.Errorf("default language = %q, want %q", cfg.Site.Language, "ko")
	}
	if cfg.Build.OutDir != "dist" {
		t.Errorf("default outDir = %q, want %q", cfg.Build.OutDir, "dist")
	}
	if cfg.Build.PageRoute != "/pages" {
		t.Errorf("default pageRoute = %q, want %q", cfg.Build.PageRoute, "/pages")
	}
	if !cfg.Theme.Layout.TOC {
		t.Error("default toc should be true")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	akwikiDir := filepath.Join(dir, ".akwiki")
	os.MkdirAll(akwikiDir, 0o755)

	yml := []byte(`site:
  title: "테스트 위키"
  author: "tester"
build:
  outDir: "output"
theme:
  layout:
    toc: false
`)
	os.WriteFile(filepath.Join(akwikiDir, "config.yml"), yml, 0o644)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Site.Title != "테스트 위키" {
		t.Errorf("title = %q, want %q", cfg.Site.Title, "테스트 위키")
	}
	if cfg.Site.Author != "tester" {
		t.Errorf("author = %q, want %q", cfg.Site.Author, "tester")
	}
	if cfg.Build.OutDir != "output" {
		t.Errorf("outDir = %q, want %q", cfg.Build.OutDir, "output")
	}
	if cfg.Theme.Layout.TOC {
		t.Error("toc should be false")
	}
	// Defaults still applied for unset fields
	if cfg.Site.Language != "ko" {
		t.Errorf("language = %q, want default %q", cfg.Site.Language, "ko")
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
go test ./internal/config/ -v
```

Expected: FAIL — package/types not defined.

- [ ] **Step 3: config.go 구현**

```go
// internal/config/config.go
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Site      SiteConfig      `yaml:"site"`
	Build     BuildConfig     `yaml:"build"`
	Analytics AnalyticsConfig `yaml:"analytics"`
	Theme     ThemeConfig     `yaml:"theme"`
}

type SiteConfig struct {
	Title    string `yaml:"title"`
	Author   string `yaml:"author"`
	URL      string `yaml:"url"`
	Language string `yaml:"language"`
}

type BuildConfig struct {
	OutDir    string `yaml:"outDir"`
	PageRoute string `yaml:"pageRoute"`
}

type AnalyticsConfig struct {
	GA string `yaml:"ga"`
}

type ThemeConfig struct {
	Colors ThemeColors `yaml:"colors"`
	Fonts  ThemeFonts  `yaml:"fonts"`
	Layout ThemeLayout `yaml:"layout"`
	Footer ThemeFooter `yaml:"footer"`
	Edit   ThemeEdit   `yaml:"edit"`
}

type ThemeColors struct {
	Background  string `yaml:"background"`
	Text        string `yaml:"text"`
	Link        string `yaml:"link"`
	LinkPrivate string `yaml:"link-private"`
	Accent      string `yaml:"accent"`
}

type ThemeFonts struct {
	Heading string `yaml:"heading"`
	Body    string `yaml:"body"`
	Code    string `yaml:"code"`
}

type ThemeLayout struct {
	MaxWidth  string `yaml:"max-width"`
	TOC       bool   `yaml:"toc"`
	Backlinks bool   `yaml:"backlinks"`
	Related   bool   `yaml:"related"`
	Search    bool   `yaml:"search"`
}

type ThemeFooter struct {
	Copyright string       `yaml:"copyright"`
	Links     []FooterLink `yaml:"links"`
}

type FooterLink struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}

type ThemeEdit struct {
	URL string `yaml:"url"`
}

func defaults() Config {
	return Config{
		Site: SiteConfig{
			Language: "ko",
		},
		Build: BuildConfig{
			OutDir:    "dist",
			PageRoute: "/pages",
		},
		Theme: ThemeConfig{
			Layout: ThemeLayout{
				TOC:       true,
				Backlinks: true,
				Related:   true,
				Search:    true,
			},
		},
	}
}

func Load(rootDir string) (*Config, error) {
	cfg := defaults()

	path := filepath.Join(rootDir, ".akwiki", "config.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

- [ ] **Step 4: yaml 의존성 추가 + 테스트 통과 확인**

```bash
go get gopkg.in/yaml.v3
go test ./internal/config/ -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/config/ go.mod go.sum
git commit -m "feat: add config system with YAML parsing and defaults"
```

---

## Task 3: 콘텐츠 모델 — frontmatter 파싱

**Files:**
- Create: `internal/content/frontmatter.go`
- Create: `internal/content/frontmatter_test.go`
- Create: `internal/content/page.go`

- [ ] **Step 1: frontmatter 파싱 테스트 작성**

```go
// internal/content/frontmatter_test.go
package content

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := `---
title: specdown
titleKo: 스펙다운
type: Book
private: true
aliases:
  - 스펙다운
  - spec-down
tags:
  - markdown
---

# specdown

This is the body content.
`
	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter() error: %v", err)
	}
	if fm.Title != "specdown" {
		t.Errorf("title = %q, want %q", fm.Title, "specdown")
	}
	if fm.TitleKo != "스펙다운" {
		t.Errorf("titleKo = %q, want %q", fm.TitleKo, "스펙다운")
	}
	if fm.Type != "Book" {
		t.Errorf("type = %q, want %q", fm.Type, "Book")
	}
	if !fm.Private {
		t.Error("private should be true")
	}
	if len(fm.Aliases) != 2 {
		t.Errorf("aliases len = %d, want 2", len(fm.Aliases))
	}
	if len(fm.Tags) != 1 || fm.Tags[0] != "markdown" {
		t.Errorf("tags = %v, want [markdown]", fm.Tags)
	}
	if string(body) != "\n# specdown\n\nThis is the body content.\n" {
		t.Errorf("body = %q", string(body))
	}
}

func TestParseFrontmatterDefaults(t *testing.T) {
	input := `---
---

Just body.
`
	fm, _, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter() error: %v", err)
	}
	if fm.Type != "Article" {
		t.Errorf("default type = %q, want %q", fm.Type, "Article")
	}
	if fm.Private {
		t.Error("default private should be false")
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	input := `# Just markdown

No frontmatter here.
`
	fm, body, err := ParseFrontmatter([]byte(input))
	if err != nil {
		t.Fatalf("ParseFrontmatter() error: %v", err)
	}
	if fm.Type != "Article" {
		t.Errorf("type = %q, want %q", fm.Type, "Article")
	}
	if string(body) != input {
		t.Errorf("body should be entire input")
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/content/ -v
```

Expected: FAIL

- [ ] **Step 3: frontmatter.go 구현**

```go
// internal/content/frontmatter.go
package content

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

type Frontmatter struct {
	Title   string   `yaml:"title"`
	TitleKo string   `yaml:"titleKo"`
	Type    string   `yaml:"type"`
	Private bool     `yaml:"private"`
	Aliases []string `yaml:"aliases"`
	Tags    []string `yaml:"tags"`
}

var fmDelimiter = []byte("---")

func ParseFrontmatter(data []byte) (*Frontmatter, []byte, error) {
	fm := &Frontmatter{
		Type: "Article",
	}

	trimmed := bytes.TrimLeft(data, " \t\n\r")
	if !bytes.HasPrefix(trimmed, fmDelimiter) {
		return fm, data, nil
	}

	// Find end delimiter
	rest := trimmed[len(fmDelimiter):]
	endIdx := bytes.Index(rest, fmDelimiter)
	if endIdx == -1 {
		return fm, data, nil
	}

	yamlBlock := rest[:endIdx]
	body := rest[endIdx+len(fmDelimiter):]

	parsed := &Frontmatter{}
	if err := yaml.Unmarshal(yamlBlock, parsed); err != nil {
		return nil, nil, err
	}

	// Apply defaults
	if parsed.Type == "" {
		parsed.Type = "Article"
	}

	return parsed, body, nil
}
```

- [ ] **Step 4: page.go 타입 정의**

```go
// internal/content/page.go
package content

import "time"

type Page struct {
	Name       string
	Title      string
	TitleKo    string
	Type       string
	Brief      string
	Private    bool
	Aliases    []string
	Tags       []string
	CreatedAt  time.Time
	ModifiedAt time.Time
	RawBody    []byte // markdown body without frontmatter
	RawSource  []byte // original full file content
}
```

- [ ] **Step 5: 테스트 통과 확인**

```bash
go test ./internal/content/ -v
```

Expected: PASS

- [ ] **Step 6: 커밋**

```bash
git add internal/content/
git commit -m "feat: add frontmatter parser and page type"
```

---

## Task 4: git 날짜 추출

**Files:**
- Create: `internal/content/gitdate.go`
- Create: `internal/content/gitdate_test.go`

- [ ] **Step 1: git 날짜 추출 테스트 작성**

```go
// internal/content/gitdate_test.go
package content

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestGitDates(t *testing.T) {
	dir := t.TempDir()

	// Init git repo
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s %v", args, out, err)
		}
	}

	run("init")
	run("config", "user.email", "test@test.com")
	run("config", "user.name", "test")

	// Create and commit a file
	testFile := filepath.Join(dir, "pages", "Hello.md")
	os.MkdirAll(filepath.Join(dir, "pages"), 0o755)
	os.WriteFile(testFile, []byte("# Hello"), 0o644)
	run("add", "pages/Hello.md")
	run("commit", "-m", "add hello")

	// Get dates
	created, modified, err := GitDates(dir, "pages/Hello.md")
	if err != nil {
		t.Fatalf("GitDates() error: %v", err)
	}

	now := time.Now()
	if now.Sub(created) > 10*time.Second {
		t.Errorf("created time too old: %v", created)
	}
	if now.Sub(modified) > 10*time.Second {
		t.Errorf("modified time too old: %v", modified)
	}
}

func TestGitDatesNoRepo(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.md")
	os.WriteFile(testFile, []byte("# Test"), 0o644)

	created, modified, err := GitDates(dir, "test.md")
	if err != nil {
		t.Fatalf("GitDates() error: %v", err)
	}

	// Should fall back to file system times
	if created.IsZero() || modified.IsZero() {
		t.Error("should return non-zero times from filesystem")
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/content/ -run TestGitDates -v
```

Expected: FAIL

- [ ] **Step 3: gitdate.go 구현**

```go
// internal/content/gitdate.go
package content

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GitDates(repoDir, relPath string) (created, modified time.Time, err error) {
	// Try git log first
	created, errC := gitFirstCommitDate(repoDir, relPath)
	modified, errM := gitLastCommitDate(repoDir, relPath)

	if errC != nil || errM != nil || created.IsZero() || modified.IsZero() {
		return fileDates(filepath.Join(repoDir, relPath))
	}

	return created, modified, nil
}

func gitFirstCommitDate(repoDir, relPath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "--diff-filter=A", "--follow",
		"--format=%aI", "--", relPath)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return time.Time{}, nil
	}
	// Last line = first commit (oldest)
	return time.Parse(time.RFC3339, lines[len(lines)-1])
}

func gitLastCommitDate(repoDir, relPath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%aI", "--", relPath)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}

func fileDates(absPath string) (created, modified time.Time, err error) {
	info, err := os.Stat(absPath)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	mod := info.ModTime()
	return mod, mod, nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/content/ -run TestGitDates -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/content/gitdate.go internal/content/gitdate_test.go
git commit -m "feat: add git-based date extraction with filesystem fallback"
```

---

## Task 5: 페이지 로더 (디렉토리 → []Page)

**Files:**
- Create: `internal/content/loader.go`
- Create: `internal/content/loader_test.go`

- [ ] **Step 1: 로더 테스트 작성**

```go
// internal/content/loader_test.go
package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPages(t *testing.T) {
	dir := t.TempDir()
	pagesDir := filepath.Join(dir, "pages")
	os.MkdirAll(pagesDir, 0o755)

	// Public page
	os.WriteFile(filepath.Join(pagesDir, "Hello.md"), []byte(`---
title: Hello World
titleKo: 안녕 세상
type: Article
tags:
  - greeting
---

# Hello

Welcome to my wiki.
`), 0o644)

	// Private page
	os.WriteFile(filepath.Join(pagesDir, "Secret.md"), []byte(`---
title: Secret
private: true
---

Hidden content.
`), 0o644)

	// No frontmatter
	os.WriteFile(filepath.Join(pagesDir, "Simple.md"), []byte(`# Simple

Just text.
`), 0o644)

	pages, err := LoadPages(pagesDir)
	if err != nil {
		t.Fatalf("LoadPages() error: %v", err)
	}

	if len(pages) != 3 {
		t.Fatalf("len(pages) = %d, want 3", len(pages))
	}

	// Check page by name
	byName := make(map[string]*Page)
	for i := range pages {
		byName[pages[i].Name] = &pages[i]
	}

	hello := byName["Hello"]
	if hello == nil {
		t.Fatal("Hello page not found")
	}
	if hello.Title != "Hello World" {
		t.Errorf("title = %q, want %q", hello.Title, "Hello World")
	}
	if hello.TitleKo != "안녕 세상" {
		t.Errorf("titleKo = %q, want %q", hello.TitleKo, "안녕 세상")
	}

	secret := byName["Secret"]
	if secret == nil || !secret.Private {
		t.Error("Secret page should exist and be private")
	}

	simple := byName["Simple"]
	if simple == nil {
		t.Fatal("Simple page not found")
	}
	if simple.Title != "Simple" {
		t.Errorf("title = %q, want %q (from filename)", simple.Title, "Simple")
	}
	if simple.Brief == "" {
		t.Error("brief should be extracted from first paragraph")
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/content/ -run TestLoadPages -v
```

Expected: FAIL

- [ ] **Step 3: loader.go 구현**

```go
// internal/content/loader.go
package content

import (
	"os"
	"path/filepath"
	"strings"
)

func LoadPages(pagesDir string) ([]Page, error) {
	entries, err := os.ReadDir(pagesDir)
	if err != nil {
		return nil, err
	}

	repoDir := filepath.Dir(pagesDir) // assume pages/ is one level below repo root

	var pages []Page
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(e.Name(), ".md")
		absPath := filepath.Join(pagesDir, e.Name())

		data, err := os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}

		fm, body, err := ParseFrontmatter(data)
		if err != nil {
			return nil, err
		}

		title := fm.Title
		if title == "" {
			title = name
		}

		relPath := filepath.Join("pages", e.Name())
		created, modified, _ := GitDates(repoDir, relPath)

		page := Page{
			Name:       name,
			Title:      title,
			TitleKo:    fm.TitleKo,
			Type:       fm.Type,
			Private:    fm.Private,
			Aliases:    fm.Aliases,
			Tags:       fm.Tags,
			CreatedAt:  created,
			ModifiedAt: modified,
			RawBody:    body,
			RawSource:  data,
			Brief:      extractBrief(body),
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func extractBrief(body []byte) string {
	text := string(body)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip headings, empty lines, images, links-only lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "![") {
			continue
		}
		// First real paragraph line
		if len(trimmed) > 200 {
			return trimmed[:200]
		}
		return trimmed
	}
	return ""
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/content/ -run TestLoadPages -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/content/loader.go internal/content/loader_test.go
git commit -m "feat: add page loader with frontmatter and brief extraction"
```

---

## Task 6: 위키링크 goldmark 확장

**Files:**
- Create: `internal/wiki/wikilink.go`
- Create: `internal/wiki/wikilink_test.go`

- [ ] **Step 1: goldmark 의존성 추가**

```bash
go get github.com/yuin/goldmark@latest
```

- [ ] **Step 2: 위키링크 파서 테스트 작성**

```go
// internal/wiki/wikilink_test.go
package wiki

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestWikilinkParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHTML string
	}{
		{
			name:     "basic wikilink",
			input:    "See [[Hello World]] for details.",
			wantHTML: `<p>See <a href="/pages/Hello%20World" class="wikilink">Hello World</a> for details.</p>`,
		},
		{
			name:     "wikilink with display text",
			input:    "Check [[Hello World|the greeting]] page.",
			wantHTML: `<p>Check <a href="/pages/Hello%20World" class="wikilink">the greeting</a> page.</p>`,
		},
		{
			name:     "multiple wikilinks",
			input:    "See [[Foo]] and [[Bar]].",
			wantHTML: `<p>See <a href="/pages/Foo" class="wikilink">Foo</a> and <a href="/pages/Bar" class="wikilink">Bar</a>.</p>`,
		},
		{
			name:     "no wikilink",
			input:    "Normal text without links.",
			wantHTML: `<p>Normal text without links.</p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(NewWikilinkExtension("/pages")),
			)
			var buf bytes.Buffer
			reader := text.NewReader([]byte(tt.input))
			doc := md.Parser().Parse(reader)
			md.Renderer().Render(&buf, []byte(tt.input), doc)

			got := bytes.TrimSpace(buf.Bytes())
			want := []byte(tt.wantHTML)
			if !bytes.Equal(got, want) {
				t.Errorf("\ngot:  %s\nwant: %s", got, want)
			}
		})
	}
}
```

- [ ] **Step 3: 테스트 실패 확인**

```bash
go test ./internal/wiki/ -run TestWikilinkParser -v
```

Expected: FAIL

- [ ] **Step 4: wikilink.go 구현 — goldmark 인라인 파서**

```go
// internal/wiki/wikilink.go
package wiki

import (
	"fmt"
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// AST node

type Wikilink struct {
	ast.BaseInline
	Target  string
	Display string
}

func (n *Wikilink) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Target":  n.Target,
		"Display": n.Display,
	}, nil)
}

var KindWikilink = ast.NewNodeKind("Wikilink")

func (n *Wikilink) Kind() ast.NodeKind { return KindWikilink }

// Parser

type wikilinkParser struct {
}

func (p *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

func (p *wikilinkParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 4 || line[0] != '[' || line[1] != '[' {
		return nil
	}

	// Find closing ]]
	end := -1
	for i := 2; i < len(line)-1; i++ {
		if line[i] == ']' && line[i+1] == ']' {
			end = i
			break
		}
	}
	if end == -1 {
		return nil
	}

	content := string(line[2:end])
	block.Advance(end + 2)

	target := content
	display := content
	for i := 0; i < len(content); i++ {
		if content[i] == '|' {
			target = content[:i]
			display = content[i+1:]
			break
		}
	}

	return &Wikilink{Target: target, Display: display}
}

// Renderer

type wikilinkRenderer struct {
	pageRoute string
}

func (r *wikilinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikilink, r.renderWikilink)
}

func (r *wikilinkRenderer) renderWikilink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*Wikilink)
	href := fmt.Sprintf("%s/%s", r.pageRoute, url.PathEscape(n.Target))
	_, _ = fmt.Fprintf(w, `<a href="%s" class="wikilink">%s</a>`, href, n.Display)
	return ast.WalkContinue, nil
}

// Extension

type wikilinkExtension struct {
	pageRoute string
}

func NewWikilinkExtension(pageRoute string) goldmark.Extender {
	return &wikilinkExtension{pageRoute: pageRoute}
}

func (e *wikilinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&wikilinkParser{}, 199),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&wikilinkRenderer{pageRoute: e.pageRoute}, 199),
		),
	)
}
```

- [ ] **Step 5: 테스트 통과 확인**

```bash
go test ./internal/wiki/ -run TestWikilinkParser -v
```

Expected: PASS

- [ ] **Step 6: 커밋**

```bash
git add internal/wiki/ go.mod go.sum
git commit -m "feat: add goldmark wikilink parser extension"
```

---

## Task 7: 링크맵 + 백링크맵

**Files:**
- Create: `internal/wiki/linkmap.go`
- Create: `internal/wiki/linkmap_test.go`

- [ ] **Step 1: 링크맵 테스트 작성**

```go
// internal/wiki/linkmap_test.go
package wiki

import (
	"testing"
)

func TestExtractWikilinks(t *testing.T) {
	body := []byte("See [[Hello]] and [[World|the world]] here. Also [[Hello]] again.")
	links := ExtractWikilinks(body)

	if len(links) != 2 {
		t.Fatalf("len(links) = %d, want 2 (deduplicated)", len(links))
	}

	targets := make(map[string]bool)
	for _, l := range links {
		targets[l] = true
	}
	if !targets["Hello"] || !targets["World"] {
		t.Errorf("links = %v, want Hello and World", links)
	}
}

func TestBuildLinkMap(t *testing.T) {
	pages := map[string][]byte{
		"Home":    []byte("Welcome. See [[About]] and [[Blog]]."),
		"About":   []byte("About page. Back to [[Home]]."),
		"Blog":    []byte("Blog page. See [[About]]."),
		"Orphan":  []byte("No links here."),
	}

	linkMap, backlinks := BuildLinkMaps(pages)

	// Home links to About and Blog
	if len(linkMap["Home"]) != 2 {
		t.Errorf("Home links = %v, want 2", linkMap["Home"])
	}

	// About is linked from Home and Blog
	if len(backlinks["About"]) != 2 {
		t.Errorf("About backlinks = %v, want 2", backlinks["About"])
	}

	// Home is linked from About
	if len(backlinks["Home"]) != 1 || backlinks["Home"][0] != "About" {
		t.Errorf("Home backlinks = %v, want [About]", backlinks["Home"])
	}

	// Orphan has no backlinks
	if len(backlinks["Orphan"]) != 0 {
		t.Errorf("Orphan backlinks = %v, want none", backlinks["Orphan"])
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/wiki/ -run "TestExtractWikilinks|TestBuildLinkMap" -v
```

Expected: FAIL

- [ ] **Step 3: linkmap.go 구현**

```go
// internal/wiki/linkmap.go
package wiki

import "regexp"

var wikilinkRe = regexp.MustCompile(`\[\[([^\]|]+)(?:\|[^\]]+)?\]\]`)

func ExtractWikilinks(body []byte) []string {
	matches := wikilinkRe.FindAllSubmatch(body, -1)
	seen := make(map[string]bool)
	var links []string
	for _, m := range matches {
		target := string(m[1])
		if !seen[target] {
			seen[target] = true
			links = append(links, target)
		}
	}
	return links
}

// BuildLinkMaps returns (page→targets, page→sources).
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
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/wiki/ -run "TestExtractWikilinks|TestBuildLinkMap" -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/wiki/linkmap.go internal/wiki/linkmap_test.go
git commit -m "feat: add wikilink extraction and bidirectional link maps"
```

---

## Task 8: TOC 추출

**Files:**
- Create: `internal/wiki/toc.go`
- Create: `internal/wiki/toc_test.go`

- [ ] **Step 1: TOC 테스트 작성**

```go
// internal/wiki/toc_test.go
package wiki

import (
	"testing"
)

func TestExtractTOC(t *testing.T) {
	markdown := []byte(`# Title

Some text.

## Section One

Content.

### Subsection

More content.

## Section Two

Final.
`)
	toc := ExtractTOC(markdown)

	// h1 excluded (page title), only h2 and h3
	if len(toc) != 3 {
		t.Fatalf("len(toc) = %d, want 3", len(toc))
	}

	if toc[0].Text != "Section One" || toc[0].Level != 2 {
		t.Errorf("toc[0] = %+v, want Section One level 2", toc[0])
	}
	if toc[1].Text != "Subsection" || toc[1].Level != 3 {
		t.Errorf("toc[1] = %+v, want Subsection level 3", toc[1])
	}
	if toc[2].Text != "Section Two" || toc[2].Level != 2 {
		t.Errorf("toc[2] = %+v, want Section Two level 2", toc[2])
	}

	// IDs should be non-empty
	for _, h := range toc {
		if h.ID == "" {
			t.Errorf("heading %q has empty ID", h.Text)
		}
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/wiki/ -run TestExtractTOC -v
```

Expected: FAIL

- [ ] **Step 3: toc.go 구현**

```go
// internal/wiki/toc.go
package wiki

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
)

type Heading struct {
	Level int
	Text  string
	ID    string
}

var headingRe = regexp.MustCompile(`(?m)^(#{2,6})\s+(.+)$`)

func ExtractTOC(markdown []byte) []Heading {
	matches := headingRe.FindAllSubmatch(markdown, -1)
	var headings []Heading

	for _, m := range matches {
		level := len(m[1])
		text := strings.TrimSpace(string(m[2]))
		id := headingID(text)
		headings = append(headings, Heading{
			Level: level,
			Text:  text,
			ID:    id,
		})
	}

	return headings
}

func headingID(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("h%x", hash[:4])
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/wiki/ -run TestExtractTOC -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/wiki/toc.go internal/wiki/toc_test.go
git commit -m "feat: add TOC extraction from markdown headings"
```

---

## Task 9: TF-IDF 관련 콘텐츠

**Files:**
- Create: `internal/search/tfidf.go`
- Create: `internal/search/tfidf_test.go`

- [ ] **Step 1: TF-IDF 테스트 작성**

```go
// internal/search/tfidf_test.go
package search

import (
	"testing"
)

func TestTFIDF(t *testing.T) {
	docs := map[string]string{
		"go-intro":    "Go is a programming language. Go is fast and concurrent.",
		"rust-intro":  "Rust is a programming language. Rust focuses on safety.",
		"go-advanced": "Advanced Go programming with goroutines and channels.",
		"cooking":     "Cooking pasta requires boiling water and sauce.",
	}

	index := NewTFIDFIndex(docs)

	// go-intro should be most similar to go-advanced
	similar := index.MostSimilar("go-intro", 3)
	if len(similar) == 0 {
		t.Fatal("no similar documents found")
	}
	if similar[0].Name != "go-advanced" {
		t.Errorf("most similar to go-intro = %q, want go-advanced", similar[0].Name)
	}

	// cooking should not be highly similar to go-intro
	for _, s := range similar {
		if s.Name == "cooking" && s.Score > 0.5 {
			t.Errorf("cooking too similar to go-intro: score=%f", s.Score)
		}
	}
}

func TestTFIDFEmpty(t *testing.T) {
	docs := map[string]string{
		"only": "single document here",
	}
	index := NewTFIDFIndex(docs)
	similar := index.MostSimilar("only", 5)
	if len(similar) != 0 {
		t.Errorf("expected no similar docs, got %d", len(similar))
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/search/ -run TestTFIDF -v
```

Expected: FAIL

- [ ] **Step 3: tfidf.go 구현**

```go
// internal/search/tfidf.go
package search

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

type SimilarDoc struct {
	Name  string
	Score float64
}

type TFIDFIndex struct {
	docs    map[string]map[string]float64 // doc → term → tf-idf
	idf     map[string]float64
	docList []string
}

func NewTFIDFIndex(docs map[string]string) *TFIDFIndex {
	idx := &TFIDFIndex{
		docs: make(map[string]map[string]float64),
		idf:  make(map[string]float64),
	}

	// Document frequency
	df := make(map[string]int)
	termFreqs := make(map[string]map[string]int)

	for name, text := range docs {
		idx.docList = append(idx.docList, name)
		tokens := tokenize(text)
		tf := make(map[string]int)
		seen := make(map[string]bool)
		for _, tok := range tokens {
			tf[tok]++
			if !seen[tok] {
				df[tok]++
				seen[tok] = true
			}
		}
		termFreqs[name] = tf
	}

	n := float64(len(docs))
	for term, count := range df {
		idx.idf[term] = math.Log(1 + n/float64(count))
	}

	for name, tf := range termFreqs {
		vec := make(map[string]float64)
		total := 0
		for _, c := range tf {
			total += c
		}
		for term, count := range tf {
			vec[term] = (float64(count) / float64(total)) * idx.idf[term]
		}
		idx.docs[name] = vec
	}

	return idx
}

func (idx *TFIDFIndex) MostSimilar(name string, limit int) []SimilarDoc {
	vec, ok := idx.docs[name]
	if !ok {
		return nil
	}

	var results []SimilarDoc
	for _, other := range idx.docList {
		if other == name {
			continue
		}
		score := cosineSimilarity(vec, idx.docs[other])
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

func cosineSimilarity(a, b map[string]float64) float64 {
	var dot, normA, normB float64
	for term, va := range a {
		if vb, ok := b[term]; ok {
			dot += va * vb
		}
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

func tokenize(text string) []string {
	var tokens []string
	for _, word := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if len(word) > 1 {
			tokens = append(tokens, word)
		}
	}
	return tokens
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/search/ -run TestTFIDF -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/search/tfidf.go internal/search/tfidf_test.go
git commit -m "feat: add TF-IDF similarity engine for related content"
```

---

## Task 10: 검색 인덱스 생성

**Files:**
- Create: `internal/search/index.go`
- Create: `internal/search/index_test.go`

- [ ] **Step 1: 인덱스 생성 테스트 작성**

```go
// internal/search/index_test.go
package search

import (
	"encoding/json"
	"testing"

	"github.com/kyungw00k/akwiki/internal/content"
)

func TestBuildSearchIndex(t *testing.T) {
	pages := []content.Page{
		{Name: "Hello", Title: "Hello World", TitleKo: "안녕 세상", Brief: "A greeting", Tags: []string{"intro"}, Aliases: []string{"greeting"}},
		{Name: "About", Title: "About Me", Brief: "About page"},
	}

	data := BuildSearchIndex(pages)

	var entries []SearchEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2", len(entries))
	}

	if entries[0].Title != "Hello World" || entries[0].TitleKo != "안녕 세상" {
		t.Errorf("entry[0] = %+v", entries[0])
	}
	if len(entries[0].Tags) != 1 || entries[0].Tags[0] != "intro" {
		t.Errorf("tags = %v", entries[0].Tags)
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/search/ -run TestBuildSearchIndex -v
```

Expected: FAIL

- [ ] **Step 3: index.go 구현**

```go
// internal/search/index.go
package search

import (
	"encoding/json"

	"github.com/kyungw00k/akwiki/internal/content"
)

type SearchEntry struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	TitleKo string   `json:"titleKo,omitempty"`
	Brief   string   `json:"brief,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
}

func BuildSearchIndex(pages []content.Page) []byte {
	var entries []SearchEntry
	for _, p := range pages {
		entries = append(entries, SearchEntry{
			Name:    p.Name,
			Title:   p.Title,
			TitleKo: p.TitleKo,
			Brief:   p.Brief,
			Tags:    p.Tags,
			Aliases: p.Aliases,
		})
	}

	data, _ := json.MarshalIndent(entries, "", "  ")
	return data
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/search/ -run TestBuildSearchIndex -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/search/index.go internal/search/index_test.go
git commit -m "feat: add search index JSON generator"
```

---

## Task 11: 마크다운 → HTML 렌더링

**Files:**
- Create: `internal/render/markdown.go`
- Create: `internal/render/markdown_test.go`

- [ ] **Step 1: 마크다운 렌더링 테스트 작성**

```go
// internal/render/markdown_test.go
package render

import (
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	input := []byte(`# Hello World

This is a paragraph with a [[Wiki Link]] and a [[Another|display text]].

## Section

- List item 1
- List item 2

> A blockquote.

` + "```go\nfmt.Println(\"hello\")\n```")

	html, err := RenderMarkdown(input, "/pages")
	if err != nil {
		t.Fatalf("RenderMarkdown() error: %v", err)
	}

	result := string(html)

	// Check wikilinks are rendered
	if !strings.Contains(result, `class="wikilink"`) {
		t.Error("wikilinks not rendered")
	}
	if !strings.Contains(result, `/pages/Wiki%20Link`) {
		t.Error("wikilink href not correct")
	}
	if !strings.Contains(result, `>display text</a>`) {
		t.Error("wikilink display text not correct")
	}

	// Check standard markdown elements
	if !strings.Contains(result, "<h1") {
		t.Error("h1 not rendered")
	}
	if !strings.Contains(result, "<blockquote>") {
		t.Error("blockquote not rendered")
	}
}

func TestRenderMarkdownHeadingIDs(t *testing.T) {
	input := []byte("## My Section\n\nContent.\n\n## Another Section\n")
	html, err := RenderMarkdown(input, "/pages")
	if err != nil {
		t.Fatalf("RenderMarkdown() error: %v", err)
	}

	result := string(html)
	if !strings.Contains(result, `id="`) {
		t.Error("headings should have IDs for TOC anchors")
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/render/ -v
```

Expected: FAIL

- [ ] **Step 3: markdown.go 구현**

```go
// internal/render/markdown.go
package render

import (
	"bytes"

	"github.com/kyungw00k/akwiki/internal/wiki"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func RenderMarkdown(source []byte, pageRoute string) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			wiki.NewWikilinkExtension(pageRoute),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/render/ -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/render/markdown.go internal/render/markdown_test.go
git commit -m "feat: add markdown to HTML renderer with wikilink support"
```

---

## Task 12: JSON-LD 생성

**Files:**
- Create: `internal/render/jsonld.go`
- Create: `internal/render/jsonld_test.go`

- [ ] **Step 1: JSON-LD 테스트 작성**

```go
// internal/render/jsonld_test.go
package render

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kyungw00k/akwiki/internal/content"
)

func TestGenerateJSONLD(t *testing.T) {
	page := content.Page{
		Title:      "Test Page",
		Brief:      "A test page",
		CreatedAt:  time.Date(2024, 9, 8, 10, 0, 0, 0, time.UTC),
		ModifiedAt: time.Date(2026, 3, 15, 6, 0, 0, 0, time.UTC),
		Type:       "Article",
	}

	result := GenerateJSONLD(page)

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if data["@context"] != "https://schema.org/" {
		t.Errorf("context = %v", data["@context"])
	}
	if data["@type"] != "Article" {
		t.Errorf("type = %v", data["@type"])
	}
	if data["name"] != "Test Page" {
		t.Errorf("name = %v", data["name"])
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/render/ -run TestGenerateJSONLD -v
```

Expected: FAIL

- [ ] **Step 3: jsonld.go 구현**

```go
// internal/render/jsonld.go
package render

import (
	"encoding/json"

	"github.com/kyungw00k/akwiki/internal/content"
)

type jsonLD struct {
	Context      string `json:"@context"`
	Type         string `json:"@type"`
	Name         string `json:"name"`
	Abstract     string `json:"abstract,omitempty"`
	DateCreated  string `json:"dateCreated"`
	DateModified string `json:"dateModified"`
}

func GenerateJSONLD(page content.Page) string {
	ld := jsonLD{
		Context:      "https://schema.org/",
		Type:         page.Type,
		Name:         page.Title,
		Abstract:     page.Brief,
		DateCreated:  page.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		DateModified: page.ModifiedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	data, _ := json.Marshal(ld)
	return string(data)
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/render/ -run TestGenerateJSONLD -v
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/render/jsonld.go internal/render/jsonld_test.go
git commit -m "feat: add Schema.org JSON-LD generator"
```

---

## Task 13: 테마 템플릿 시스템

**Files:**
- Create: `internal/render/template.go`
- Create: `internal/render/template_test.go`
- Create: `theme/default/templates/page.html`
- Create: `theme/default/templates/partials/header.html`
- Create: `theme/default/templates/partials/toc.html`
- Create: `theme/default/templates/partials/backlinks.html`
- Create: `theme/default/templates/partials/related.html`
- Create: `theme/default/templates/partials/search.html`
- Create: `theme/default/templates/partials/footer.html`

- [ ] **Step 1: 기본 테마 HTML 템플릿 작성 — page.html**

```html
<!-- theme/default/templates/page.html -->
<!DOCTYPE html>
<html lang="{{ .Site.Language }}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Page.Title }}{{ if .Site.Title }} — {{ .Site.Title }}{{ end }}</title>
  <link rel="stylesheet" href="{{ .Site.BaseURL }}/assets/style.css">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Noto+Serif+KR:wght@200..900&display=swap" rel="stylesheet">
  {{ if .Site.Analytics.GA }}<script async src="https://www.googletagmanager.com/gtag/js?id={{ .Site.Analytics.GA }}"></script>
  <script>window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','{{ .Site.Analytics.GA }}');</script>{{ end }}
  <script type="application/ld+json">{{ .JSONLD }}</script>
  <style>{{ .ThemeCSS }}</style>
</head>
<body>
  <a href="#content" class="skip-link">본문으로 건너뛰기</a>
  <div data-layout="wiki-page">
    {{ template "header" . }}
    {{ template "toc" . }}
    <article id="content" data-content>
      {{ .Content }}
    </article>
    {{ template "related" . }}
    {{ template "footer" . }}
  </div>
  {{ if .Site.Theme.Layout.Search }}{{ template "search" . }}{{ end }}
</body>
</html>
```

- [ ] **Step 2: partials 작성**

```html
<!-- theme/default/templates/partials/header.html -->
{{ define "header" }}
<nav data-area="nav">
  <a href="{{ .Site.BaseURL }}/pages/Home">위키 홈</a>
  {{ if .Site.Theme.Edit.URL }}<a href="{{ editURL .Site.Theme.Edit.URL .Page.Name }}">Edit</a>{{ end }}
</nav>
<header data-area="header">
  <h1>{{ if .Page.TitleKo }}{{ .Page.TitleKo }}{{ else }}{{ .Page.Title }}{{ end }}</h1>
  <time>{{ formatDate .Page.CreatedAt }}{{ if not (sameDay .Page.CreatedAt .Page.ModifiedAt) }} (modified: {{ formatDate .Page.ModifiedAt }}){{ end }}</time>
  {{ if .Page.Aliases }}<div class="aliases">{{ range .Page.Aliases }}<span>{{ . }}</span>{{ end }}</div>{{ end }}
</header>
{{ end }}
```

```html
<!-- theme/default/templates/partials/toc.html -->
{{ define "toc" }}
{{ if and .Site.Theme.Layout.TOC .TOC }}
<nav aria-label="목차" data-area="toc">
  <h2>목차</h2>
  <ol>
    <li><a href="#top">(맨 위로)</a></li>
    {{ range .TOC }}<li class="toc-level-{{ .Level }}"><a href="#{{ .ID }}">{{ .Text }}</a></li>{{ end }}
  </ol>
</nav>
{{ end }}
{{ end }}
```

```html
<!-- theme/default/templates/partials/backlinks.html -->
{{ define "backlinks" }}
{{ if and .Site.Theme.Layout.Backlinks .Backlinks }}
<aside data-area="backlinks">
  <h2>역링크</h2>
  <ul>
    {{ range .Backlinks }}<li><a href="{{ $.Site.BaseURL }}/pages/{{ urlPathEscape .Name }}">{{ .Title }}</a></li>{{ end }}
  </ul>
</aside>
{{ end }}
{{ end }}
```

```html
<!-- theme/default/templates/partials/related.html -->
{{ define "related" }}
{{ if and .Site.Theme.Layout.Related .Related }}
<aside data-area="related">
  {{ range $type, $pages := .Related }}
  {{ if $pages }}
  <section>
    <h2>{{ relatedLabel $type }}</h2>
    <ul>
      {{ range $pages }}<li><a href="{{ $.Site.BaseURL }}/pages/{{ urlPathEscape .Name }}">{{ .Title }}</a></li>{{ end }}
    </ul>
  </section>
  {{ end }}
  {{ end }}
</aside>
{{ end }}
{{ end }}
```

```html
<!-- theme/default/templates/partials/search.html -->
{{ define "search" }}
<div id="search-container" style="display:none">
  <input type="search" id="search-input" placeholder="검색..." aria-label="검색">
  <ul id="search-results"></ul>
</div>
<script src="{{ .Site.BaseURL }}/assets/search.js"></script>
{{ end }}
```

```html
<!-- theme/default/templates/partials/footer.html -->
{{ define "footer" }}
<footer data-area="footer">
  <p>
    {{ if .Site.Theme.Footer.Copyright }}{{ .Site.Theme.Footer.Copyright }}{{ end }}
    {{ range .Site.Theme.Footer.Links }} | <a href="{{ .URL }}">{{ .Label }}</a>{{ end }}
    | <a href="{{ .Page.RawURL }}" title="Markdown Version">markdown</a>
  </p>
  <p class="credits">
    Inspired by <a href="https://github.com/akngs">akngs</a>'s <a href="https://wiki.g15e.com/pages/Home">wiki</a>.
    Powered by <a href="https://github.com/kyungw00k/akwiki">akwiki</a>.
  </p>
</footer>
{{ end }}
```

- [ ] **Step 3: 템플릿 로더 테스트 작성**

```go
// internal/render/template_test.go
package render

import (
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/content"
	wikiPkg "github.com/kyungw00k/akwiki/internal/wiki"
)

func TestRenderPage(t *testing.T) {
	cfg := config.Config{
		Site: config.SiteConfig{
			Title:    "Test Wiki",
			Language: "ko",
		},
		Theme: config.ThemeConfig{
			Layout: config.ThemeLayout{
				TOC:       true,
				Backlinks: true,
				Related:   true,
				Search:    true,
			},
		},
	}

	ctx := &TemplateContext{
		Site:    &cfg,
		Page:    content.Page{Name: "Hello", Title: "Hello World", CreatedAt: time.Now(), ModifiedAt: time.Now()},
		Content: template.HTML("<p>Hello</p>"),
		TOC:     []wikiPkg.Heading{{Level: 2, Text: "Section", ID: "h1234"}},
		JSONLD:  template.JS(`{"@type":"Article"}`),
	}

	engine, err := NewTemplateEngine("", "") // use embedded defaults
	if err != nil {
		t.Fatalf("NewTemplateEngine() error: %v", err)
	}

	html, err := engine.RenderPage(ctx)
	if err != nil {
		t.Fatalf("RenderPage() error: %v", err)
	}

	result := string(html)
	if !strings.Contains(result, "Hello World") {
		t.Error("page title not in output")
	}
	if !strings.Contains(result, "위키 홈") {
		t.Error("wiki home link not in output")
	}
	if !strings.Contains(result, "본문으로 건너뛰기") {
		t.Error("skip link not in output")
	}
	if !strings.Contains(result, "akngs") {
		t.Error("credits not in output")
	}
	if !strings.Contains(result, "Section") {
		t.Error("TOC entry not in output")
	}
}
```

- [ ] **Step 4: template.go 구현**

```go
// internal/render/template.go
package render

import (
	"bytes"
	"embed"
	"fmt"
	htmltemplate "html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/content"
	"github.com/kyungw00k/akwiki/internal/search"
	"github.com/kyungw00k/akwiki/internal/wiki"
)

//go:embed ../../../theme/default/templates/*.html ../../../theme/default/templates/partials/*.html
var defaultTemplates embed.FS

type TemplateContext struct {
	Site      *config.Config
	Page      content.Page
	Content   htmltemplate.HTML
	TOC       []wiki.Heading
	Links     []PageRef
	Backlinks []PageRef
	Related   map[string][]PageRef
	JSONLD    htmltemplate.JS
	ThemeCSS  htmltemplate.CSS
}

type PageRef struct {
	Name  string
	Title string
	Brief string
	Type  string
	Score float64
}

func NewPageRef(name, title, brief, typ string, score float64) PageRef {
	return PageRef{Name: name, Title: title, Brief: brief, Type: typ, Score: score}
}

func PageRefsFromSimilar(docs []search.SimilarDoc, pages map[string]content.Page) []PageRef {
	var refs []PageRef
	for _, d := range docs {
		if p, ok := pages[d.Name]; ok {
			refs = append(refs, PageRef{
				Name: p.Name, Title: p.Title, Brief: p.Brief,
				Type: p.Type, Score: d.Score,
			})
		}
	}
	return refs
}

type TemplateEngine struct {
	tmpl *htmltemplate.Template
}

var funcMap = htmltemplate.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"sameDay": func(a, b time.Time) bool {
		return a.Format("2006-01-02") == b.Format("2006-01-02")
	},
	"editURL": func(pattern, pagename string) string {
		return strings.ReplaceAll(pattern, "{{pagename}}", pagename)
	},
	"urlPathEscape": func(s string) string {
		return url.PathEscape(s)
	},
	"relatedLabel": func(typ string) string {
		switch typ {
		case "Book":
			return "관련 책"
		case "Article":
			return "관련 글"
		case "ScholarlyArticle":
			return "관련 논문"
		default:
			return "관련 " + typ
		}
	},
}

func NewTemplateEngine(themeOverrideDir, embeddedBase string) (*TemplateEngine, error) {
	tmpl := htmltemplate.New("page.html").Funcs(funcMap)

	// Load embedded defaults
	partials := []string{"header", "toc", "backlinks", "related", "search", "footer"}
	for _, name := range partials {
		embPath := fmt.Sprintf("theme/default/templates/partials/%s.html", name)
		overridePath := ""
		if themeOverrideDir != "" {
			overridePath = filepath.Join(themeOverrideDir, "templates", "partials", name+".html")
		}

		data, err := loadTemplate(embPath, overridePath)
		if err != nil {
			return nil, fmt.Errorf("load partial %s: %w", name, err)
		}
		if _, err := tmpl.Parse(string(data)); err != nil {
			return nil, fmt.Errorf("parse partial %s: %w", name, err)
		}
	}

	// Load main page template
	mainOverride := ""
	if themeOverrideDir != "" {
		mainOverride = filepath.Join(themeOverrideDir, "templates", "page.html")
	}
	mainData, err := loadTemplate("theme/default/templates/page.html", mainOverride)
	if err != nil {
		return nil, fmt.Errorf("load page template: %w", err)
	}
	if _, err := tmpl.Parse(string(mainData)); err != nil {
		return nil, fmt.Errorf("parse page template: %w", err)
	}

	return &TemplateEngine{tmpl: tmpl}, nil
}

func loadTemplate(embeddedPath, overridePath string) ([]byte, error) {
	if overridePath != "" {
		data, err := os.ReadFile(overridePath)
		if err == nil {
			return data, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	// Normalize embedded path to match embed.FS structure
	// embed.FS uses paths relative to the go:embed directive location
	// Since we embed from internal/render/, paths are relative to that
	return defaultTemplates.ReadFile(embeddedPath)
}

func (e *TemplateEngine) RenderPage(ctx *TemplateContext) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.tmpl.Execute(&buf, ctx); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
```

- [ ] **Step 5: embed 경로 수정 — theme 디렉토리를 올바르게 embed하기 위해 render 패키지에 embed.go 배치**

`internal/render/template.go`의 embed 지시문은 상대 경로 제약이 있으므로, 프로젝트 루트에 embed 패키지를 두는 것이 더 깔끔합니다:

```go
// theme/embed.go
package theme

import "embed"

//go:embed default/templates/*.html default/templates/partials/*.html default/static/*
var DefaultTheme embed.FS
```

그리고 `internal/render/template.go`에서 `defaultTemplates`를 외부에서 주입받도록 수정:

```go
// internal/render/template.go 상단의 embed 지시문 제거하고, NewTemplateEngine에 embed.FS 파라미터 추가
func NewTemplateEngine(defaultFS embed.FS, themeOverrideDir string) (*TemplateEngine, error) {
```

`loadTemplate`도 수정:

```go
func loadTemplate(defaultFS embed.FS, embeddedPath, overridePath string) ([]byte, error) {
	if overridePath != "" {
		data, err := os.ReadFile(overridePath)
		if err == nil {
			return data, nil
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return defaultFS.ReadFile(embeddedPath)
}
```

- [ ] **Step 6: 테스트 통과 확인**

```bash
go test ./internal/render/ -v
```

Expected: PASS

- [ ] **Step 7: 커밋**

```bash
git add theme/ internal/render/template.go internal/render/template_test.go
git commit -m "feat: add template engine with embedded theme and override support"
```

---

## Task 14: 기본 테마 CSS

**Files:**
- Create: `theme/default/static/style.css`

- [ ] **Step 1: 원본 재현 CSS 작성**

```css
/* theme/default/static/style.css */

/* === Reset === */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* === Color System (OKLCh) === */
:root {
  --hue-primary: 90;
  --hue-complement: 270;
  --hue-triadic: 210;
  --hue-warning: 30;
  --hue-info: 250;

  --l-offset: 0;
  --l-sign: 1;

  --c-bg: oklch(calc((97 - var(--l-offset)) * 1%) 0.02 var(--hue-primary));
  --c-bg-code: oklch(calc((92 - var(--l-offset)) * 1%) 0.02 var(--hue-primary));
  --c-bg-highlight: oklch(calc((85 - var(--l-offset)) * 1%) 0.05 var(--hue-primary));
  --c-text: oklch(calc((30 + var(--l-offset) * 0.4) * 1%) 0.007 var(--hue-primary));
  --c-text-muted: oklch(calc((45 + var(--l-offset) * 0.1) * 1%) 0.007 var(--hue-primary));
  --c-text-highlight: oklch(calc((20 + var(--l-offset) * 0.6) * 1%) 0.007 var(--hue-primary));
  --c-text-link: oklch(calc((25 + var(--l-offset) * 0.5) * 1%) 0.12 var(--hue-triadic));
  --c-text-link-ext: oklch(calc((25 + var(--l-offset) * 0.5) * 1%) 0.18 var(--hue-triadic));
  --c-accent: oklch(50% 0.28 var(--hue-complement));

  /* Fonts */
  --font-headings: "Noto Serif", "Noto Serif KR", serif;
  --font-text: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif;
  --font-text-alt: "Noto Serif", "Noto Serif KR", serif;
  --font-code: monospace;

  /* Font sizes */
  --text-xs: 0.75em;
  --text-s: 0.85rem;
  --text-l: 1.2em;
  --text-xl: 1.35em;
  --text-2xl: 1.65rem;
  --text-3xl: 1.8em;
  --text-4xl: max(1.8rem, 2.5dvw);

  /* Line heights */
  --leading-tight: 1.3;
  --leading-snug: 1.4;
  --leading-normal: 1.5;
  --leading-relaxed: 1.6;
  --leading-loose: 1.65;

  /* Spacing */
  --space-xs: 0.25rem;
  --space-s: 0.5rem;
  --space-m: 1rem;
  --space-l: 2rem;
  --space-xl: 3rem;
  --space-2xl: 6rem;

  /* Layout */
  --measure-content: 50rem;
  --measure-page: 100rem;
  --content-indent: 1.5em;
}

/* === Dark Mode === */
@media (prefers-color-scheme: dark) {
  :root {
    --l-offset: 100;
    --l-sign: -1;
  }
}

/* === Base === */
html {
  font-family: var(--font-text);
  line-height: var(--leading-loose);
  color: var(--c-text);
  background-color: var(--c-bg);
  word-break: keep-all;
  overflow-wrap: break-word;
}

/* === Skip Link === */
.skip-link {
  position: absolute;
  left: -9999px;
  top: 0;
  z-index: 100;
  padding: var(--space-s) var(--space-m);
  background: var(--c-bg);
  color: var(--c-text-link);
}
.skip-link:focus { left: 0; }

/* === Layout === */
[data-layout="wiki-page"] {
  max-width: var(--measure-content);
  margin: 0 auto;
  padding-inline: var(--space-m);
  display: grid;
  grid-template-areas:
    "nav"
    "header"
    "toc"
    "content"
    "related"
    "footer";
}

@media (min-width: 72rem) {
  [data-layout="wiki-page"] {
    max-width: var(--measure-page);
    grid-template-columns: 2fr 5fr 2fr;
    grid-template-areas:
      ".   nav     ."
      ".   header  ."
      "toc content related"
      ".   footer  .";
    column-gap: var(--space-xl);
  }
}

[data-area="nav"]       { grid-area: nav; }
[data-area="header"]    { grid-area: header; }
[data-area="toc"]       { grid-area: toc; }
[data-content]          { grid-area: content; }
[data-area="related"]   { grid-area: related; }
[data-area="footer"]    { grid-area: footer; }

/* === Nav === */
[data-area="nav"] {
  display: flex;
  gap: var(--space-m);
  padding: var(--space-m) 0;
}
[data-area="nav"] a {
  color: var(--c-text-muted);
  text-decoration: none;
  font-size: var(--text-s);
}
[data-area="nav"] a:hover { color: var(--c-text); }

/* === Header === */
[data-area="header"] {
  padding-bottom: var(--space-l);
}
[data-area="header"] h1 {
  font-family: var(--font-headings);
  font-size: var(--text-4xl);
  font-weight: 700;
  line-height: var(--leading-snug);
  color: var(--c-text);
}
[data-area="header"] time {
  font-size: var(--text-s);
  color: var(--c-text-muted);
}
.aliases {
  font-size: var(--text-s);
  color: var(--c-text-muted);
  margin-top: var(--space-xs);
}
.aliases span::after { content: ", "; }
.aliases span:last-child::after { content: ""; }

/* === TOC === */
[data-area="toc"] h2 {
  font-family: var(--font-headings);
  font-size: var(--text-2xl);
  font-weight: 600;
  margin-bottom: var(--space-s);
}
[data-area="toc"] ol {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-s);
}
[data-area="toc"] a {
  text-decoration: none;
  color: var(--c-text-link);
}
[data-area="toc"] a:hover { text-decoration: underline; }
.toc-level-3 { padding-left: var(--content-indent); }
.toc-level-4 { padding-left: calc(var(--content-indent) * 2); }

@media (min-width: 72rem) {
  [data-area="toc"] {
    position: sticky;
    top: 0;
    align-self: start;
    max-height: 100vh;
    overflow-y: auto;
  }
  [data-area="toc"] h2 { font-size: 1.25rem; }
  [data-area="toc"] ol { font-size: var(--text-s); line-height: var(--leading-normal); }
}

/* === Content === */
[data-content] {
  line-height: var(--leading-relaxed);
}
[data-content] h2, [data-content] h3, [data-content] h4,
[data-content] h5, [data-content] h6 {
  font-family: var(--font-headings);
  line-height: var(--leading-tight);
  margin-top: var(--space-l);
  margin-bottom: var(--space-s);
}
[data-content] h2 { font-size: var(--text-2xl); font-weight: 600; }
[data-content] h3 { font-size: var(--text-xl); font-weight: 600; }
[data-content] h4 { font-size: var(--text-l); font-weight: 600; }
[data-content] p { margin-bottom: var(--space-m); }
[data-content] ul, [data-content] ol {
  padding-left: var(--content-indent);
  margin-bottom: var(--space-m);
}
[data-content] li { margin-bottom: var(--space-xs); }
[data-content] img {
  max-width: 100%;
  height: auto;
}

/* === Links === */
a {
  color: var(--c-text-link);
  text-decoration: underline;
  text-underline-offset: 0.25em;
  text-decoration-thickness: 0.5px;
}
a[href^="http"]::after,
a[href^="https"]::after {
  content: " ↗";
  font-size: 0.5em;
  vertical-align: super;
}
/* Don't add arrow to internal absolute URLs */
a[href^="http"].wikilink::after { content: ""; }
a.wikilink::after { content: ""; }
a[data-link="private"] {
  color: var(--c-text-muted);
  text-decoration-style: dashed;
  text-decoration-thickness: 0.5px;
}
a.wikilink-missing {
  color: var(--c-text-muted);
  text-decoration-style: dotted;
}

/* === Inline elements === */
strong {
  font-weight: inherit;
  color: var(--c-text-highlight);
  background-color: var(--c-bg-highlight);
  padding-inline: 0.125em;
}
code {
  font-family: var(--font-code);
  background-color: var(--c-bg-code);
  border-radius: 0.25em;
  padding: 0.125em 0.25em;
}
pre {
  background-color: var(--c-bg-code);
  margin-block: 1.5em;
  overflow-x: auto;
}
pre code {
  display: block;
  padding: 1em 1.5em;
  background: none;
  letter-spacing: -0.01em;
}

/* === Blockquote === */
blockquote {
  border-inline-start: 0.25em solid var(--c-bg-highlight);
  padding: 0.5em calc(1.5em - 0.25em);
  font-family: var(--font-text-alt);
  color: var(--c-text);
}

/* === HR === */
hr {
  border: none;
  text-align: center;
  margin-block: var(--space-l);
}
hr::after {
  content: "- - - § - - -";
  color: var(--c-text-muted);
  font-size: var(--text-s);
}

/* === Table === */
table {
  border-collapse: collapse;
  margin-bottom: var(--space-m);
  width: 100%;
}
th, td {
  padding: 0.25em 1em;
  border: 1px solid var(--c-text-muted);
  text-align: left;
}
th {
  font-weight: 600;
  background-color: var(--c-bg-code);
}

/* === Related / Backlinks === */
[data-area="related"], [data-area="backlinks"] {
  margin-top: var(--space-l);
}
[data-area="related"] h2, [data-area="backlinks"] h2 {
  font-family: var(--font-headings);
  font-size: var(--text-l);
  font-weight: 600;
  margin-bottom: var(--space-s);
}
[data-area="related"] ul, [data-area="backlinks"] ul {
  list-style: none;
  padding: 0;
}
[data-area="related"] li, [data-area="backlinks"] li {
  margin-bottom: var(--space-xs);
}

@media (min-width: 72rem) {
  [data-area="related"], [data-area="backlinks"] {
    position: sticky;
    top: 0;
    align-self: start;
    font-size: var(--text-s);
  }
}

/* === Footer === */
[data-area="footer"] {
  margin-top: var(--space-2xl);
  padding: var(--space-l) 0;
  border-top: 1px solid var(--c-bg-highlight);
  font-size: var(--text-s);
  color: var(--c-text-muted);
}
[data-area="footer"] a { color: var(--c-text-muted); }
.credits { margin-top: var(--space-s); font-size: var(--text-xs); }

/* === Search === */
#search-container {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 50;
  background: var(--c-bg);
  padding: var(--space-m);
  border-bottom: 1px solid var(--c-bg-highlight);
}
#search-input {
  width: 100%;
  max-width: var(--measure-content);
  margin: 0 auto;
  display: block;
  padding: var(--space-s) var(--space-m);
  border: 1px solid var(--c-text-muted);
  border-radius: 0.25em;
  font-size: 1rem;
  font-family: var(--font-text);
  background: var(--c-bg);
  color: var(--c-text);
}
#search-results {
  max-width: var(--measure-content);
  margin: var(--space-s) auto 0;
  list-style: none;
  padding: 0;
}
#search-results li {
  padding: var(--space-xs) var(--space-m);
}
#search-results a {
  text-decoration: none;
}

/* === Progress bar === */
.progress-bar {
  position: fixed;
  top: 0; left: 0;
  width: 100%;
  height: 3px;
  background: var(--c-accent);
  animation: progress 3s cubic-bezier(.19,1.07,.23,.94);
  z-index: 100;
}
@keyframes progress { from { width: 0; } to { width: 100%; } }

/* === Theme variable overrides from config === */
```

- [ ] **Step 2: 커밋**

```bash
git add theme/default/static/style.css
git commit -m "feat: add default theme CSS with OKLCh color system"
```

---

## Task 15: 클라이언트 검색 JS

**Files:**
- Create: `theme/default/static/search.js`

- [ ] **Step 1: 인라인 퍼지 검색 JS 작성**

```javascript
// theme/default/static/search.js
(function() {
  var index = null;
  var container = document.getElementById('search-container');
  var input = document.getElementById('search-input');
  var results = document.getElementById('search-results');
  if (!container || !input) return;

  // Toggle search with Ctrl+K / Cmd+K
  document.addEventListener('keydown', function(e) {
    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      if (container.style.display === 'none') {
        container.style.display = 'block';
        input.focus();
      } else {
        container.style.display = 'none';
      }
    }
    if (e.key === 'Escape') {
      container.style.display = 'none';
    }
  });

  // Load index on first focus
  input.addEventListener('focus', function() {
    if (index) return;
    var base = document.querySelector('link[rel="stylesheet"]');
    var baseURL = base ? base.href.replace(/\/assets\/style\.css$/, '') : '';
    fetch(baseURL + '/search-index.json')
      .then(function(r) { return r.json(); })
      .then(function(data) { index = data; });
  });

  // Simple fuzzy search
  function search(query) {
    if (!index || !query) return [];
    var q = query.toLowerCase();
    var scored = [];
    for (var i = 0; i < index.length; i++) {
      var entry = index[i];
      var score = 0;
      var fields = [entry.title, entry.titleKo, entry.brief].concat(entry.tags || []).concat(entry.aliases || []);
      for (var j = 0; j < fields.length; j++) {
        if (fields[j] && fields[j].toLowerCase().indexOf(q) !== -1) {
          score += (j === 0 || j === 1) ? 10 : 1;
        }
      }
      if (score > 0) scored.push({ entry: entry, score: score });
    }
    scored.sort(function(a, b) { return b.score - a.score; });
    return scored.slice(0, 10);
  }

  input.addEventListener('input', function() {
    var matches = search(input.value);
    results.innerHTML = '';
    for (var i = 0; i < matches.length; i++) {
      var m = matches[i].entry;
      var li = document.createElement('li');
      var a = document.createElement('a');
      a.href = (document.querySelector('meta[name="base-url"]') || {}).content || '';
      a.href += '/pages/' + encodeURIComponent(m.name);
      a.textContent = m.titleKo || m.title;
      if (m.brief) {
        var span = document.createElement('span');
        span.textContent = ' — ' + m.brief;
        span.style.color = 'var(--c-text-muted)';
        span.style.fontSize = 'var(--text-s)';
        a.appendChild(span);
      }
      li.appendChild(a);
      results.appendChild(li);
    }
  });
})();
```

- [ ] **Step 2: 커밋**

```bash
git add theme/default/static/search.js
git commit -m "feat: add client-side fuzzy search with keyboard shortcut"
```

---

## Task 16: 빌드 파이프라인 오케스트레이션

**Files:**
- Create: `internal/builder/builder.go`
- Create: `internal/builder/builder_test.go`

- [ ] **Step 1: 빌드 통합 테스트 작성**

```go
// internal/builder/builder_test.go
package builder

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()

	// Set up wiki structure
	os.MkdirAll(filepath.Join(dir, "pages"), 0o755)
	os.MkdirAll(filepath.Join(dir, "public"), 0o755)
	os.MkdirAll(filepath.Join(dir, ".akwiki"), 0o755)

	os.WriteFile(filepath.Join(dir, ".akwiki", "config.yml"), []byte(`
site:
  title: "Test Wiki"
  author: "tester"
`), 0o644)

	os.WriteFile(filepath.Join(dir, "pages", "Home.md"), []byte(`---
title: Home
titleKo: 위키 홈
---

# Home

Welcome. See [[About]] for more.
`), 0o644)

	os.WriteFile(filepath.Join(dir, "pages", "About.md"), []byte(`---
title: About
type: Article
---

# About

This is the about page. Back to [[Home]].
`), 0o644)

	os.WriteFile(filepath.Join(dir, "pages", "Secret.md"), []byte(`---
title: Secret
private: true
---

Hidden.
`), 0o644)

	// Init git for date extraction
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Run()
	}
	run("init")
	run("config", "user.email", "t@t.com")
	run("config", "user.name", "t")
	run("add", ".")
	run("commit", "-m", "init")

	outDir := filepath.Join(dir, "dist")
	err := Build(dir, outDir)
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	// Check output files exist
	checks := []string{
		"index.html",
		"pages/Home/index.html",
		"pages/About/index.html",
		"pages/Home.txt",
		"pages/About.txt",
		"search-index.json",
		"assets/style.css",
		"assets/search.js",
	}
	for _, f := range checks {
		path := filepath.Join(outDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("missing output file: %s", f)
		}
	}

	// Secret page should NOT exist
	if _, err := os.Stat(filepath.Join(outDir, "pages/Secret/index.html")); !os.IsNotExist(err) {
		t.Error("private page Secret should not be in output")
	}

	// Check Home HTML contains wikilink to About
	homeHTML, _ := os.ReadFile(filepath.Join(outDir, "pages/Home/index.html"))
	if len(homeHTML) == 0 {
		t.Fatal("Home HTML is empty")
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/builder/ -v -timeout 30s
```

Expected: FAIL

- [ ] **Step 3: builder.go 구현**

```go
// internal/builder/builder.go
package builder

import (
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/content"
	"github.com/kyungw00k/akwiki/internal/render"
	"github.com/kyungw00k/akwiki/internal/search"
	"github.com/kyungw00k/akwiki/internal/wiki"
	"github.com/kyungw00k/akwiki/theme"
)

func Build(rootDir, outDir string) error {
	// 1. Load config
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	pagesDir := filepath.Join(rootDir, "pages")

	// 2. Load all pages
	allPages, err := content.LoadPages(pagesDir)
	if err != nil {
		return fmt.Errorf("load pages: %w", err)
	}

	// 3. Filter private pages
	var pages []content.Page
	privateNames := make(map[string]bool)
	for _, p := range allPages {
		if p.Private {
			privateNames[p.Name] = true
		} else {
			pages = append(pages, p)
		}
	}

	// Build alias map (alias → page name)
	aliasMap := make(map[string]string)
	for _, p := range pages {
		for _, a := range p.Aliases {
			aliasMap[a] = p.Name
		}
	}

	// 4. Build link maps
	bodyMap := make(map[string][]byte)
	for _, p := range pages {
		bodyMap[p.Name] = p.RawBody
	}
	linkMap, backlinkMap := wiki.BuildLinkMaps(bodyMap)

	// 5. TF-IDF
	docTexts := make(map[string]string)
	for _, p := range pages {
		docTexts[p.Name] = string(p.RawBody)
	}
	tfidfIndex := search.NewTFIDFIndex(docTexts)

	// Page lookup
	pageByName := make(map[string]content.Page)
	for _, p := range pages {
		pageByName[p.Name] = p
	}

	// 6. Prepare output directory
	os.RemoveAll(outDir)
	os.MkdirAll(outDir, 0o755)
	os.MkdirAll(filepath.Join(outDir, "pages"), 0o755)
	os.MkdirAll(filepath.Join(outDir, "assets"), 0o755)

	// 7. Copy static assets from embedded theme
	copyEmbeddedStatic(theme.DefaultTheme, outDir)

	// Copy user public/ directory
	publicDir := filepath.Join(rootDir, "public")
	if info, err := os.Stat(publicDir); err == nil && info.IsDir() {
		copyDir(publicDir, filepath.Join(outDir, "assets"))
	}

	// Copy user theme static overrides
	userStaticDir := filepath.Join(rootDir, ".akwiki", "theme", "static")
	if info, err := os.Stat(userStaticDir); err == nil && info.IsDir() {
		copyDir(userStaticDir, filepath.Join(outDir, "assets"))
	}

	// 8. Template engine
	themeOverrideDir := filepath.Join(rootDir, ".akwiki", "theme")
	engine, err := render.NewTemplateEngine(theme.DefaultTheme, themeOverrideDir)
	if err != nil {
		return fmt.Errorf("template engine: %w", err)
	}

	// 9. Render each page
	for _, page := range pages {
		pageRoute := cfg.Build.PageRoute

		// Render markdown to HTML
		htmlContent, err := render.RenderMarkdown(page.RawBody, pageRoute)
		if err != nil {
			return fmt.Errorf("render %s: %w", page.Name, err)
		}

		// TOC
		toc := wiki.ExtractTOC(page.RawBody)

		// Backlinks
		var backlinks []render.PageRef
		for _, blName := range backlinkMap[page.Name] {
			if p, ok := pageByName[blName]; ok {
				backlinks = append(backlinks, render.NewPageRef(p.Name, p.Title, p.Brief, p.Type, 0))
			}
		}

		// Related content (grouped by type)
		similarDocs := tfidfIndex.MostSimilar(page.Name, 10)
		related := make(map[string][]render.PageRef)
		for _, sd := range similarDocs {
			if p, ok := pageByName[sd.Name]; ok {
				ref := render.NewPageRef(p.Name, p.Title, p.Brief, p.Type, sd.Score)
				related[p.Type] = append(related[p.Type], ref)
			}
		}

		// Links from this page
		var links []render.PageRef
		for _, target := range linkMap[page.Name] {
			if p, ok := pageByName[target]; ok {
				links = append(links, render.NewPageRef(p.Name, p.Title, p.Brief, p.Type, 0))
			}
		}

		// JSON-LD
		jsonLD := render.GenerateJSONLD(page)

		baseURL := strings.TrimRight(cfg.Site.URL, "/")
		rawURL := fmt.Sprintf("%s/pages/%s.txt", baseURL, page.Name)
		page.RawURL = rawURL

		ctx := &render.TemplateContext{
			Site:      cfg,
			Page:      page,
			Content:   htmltemplate.HTML(htmlContent),
			TOC:       toc,
			Links:     links,
			Backlinks: backlinks,
			Related:   related,
			JSONLD:    htmltemplate.JS(jsonLD),
		}

		output, err := engine.RenderPage(ctx)
		if err != nil {
			return fmt.Errorf("render template %s: %w", page.Name, err)
		}

		// Write HTML
		pageDir := filepath.Join(outDir, "pages", page.Name)
		os.MkdirAll(pageDir, 0o755)
		if err := os.WriteFile(filepath.Join(pageDir, "index.html"), output, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", page.Name, err)
		}

		// Write .txt (raw markdown source)
		txtPath := filepath.Join(outDir, "pages", page.Name+".txt")
		if err := os.WriteFile(txtPath, page.RawSource, 0o644); err != nil {
			return fmt.Errorf("write txt %s: %w", page.Name, err)
		}
	}

	// 10. Search index
	searchData := search.BuildSearchIndex(pages)
	os.WriteFile(filepath.Join(outDir, "search-index.json"), searchData, 0o644)

	// 11. Index redirect
	indexHTML := fmt.Sprintf(`<!DOCTYPE html><html><head><meta http-equiv="refresh" content="0;url=%s/pages/Home"></head></html>`, strings.TrimRight(cfg.Site.URL, "/"))
	os.WriteFile(filepath.Join(outDir, "index.html"), []byte(indexHTML), 0o644)

	return nil
}

func copyEmbeddedStatic(fsys embed.FS, outDir string) {
	fs.WalkDir(fsys, "default/static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}
		name := filepath.Base(path)
		return os.WriteFile(filepath.Join(outDir, "assets", name), data, 0o644)
	})
}

func copyDir(src, dst string) {
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)
		os.MkdirAll(filepath.Dir(dstPath), 0o755)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, 0o644)
	})
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/builder/ -v -timeout 30s
```

Expected: PASS

- [ ] **Step 5: 커밋**

```bash
git add internal/builder/ theme/embed.go
git commit -m "feat: add build pipeline orchestrator"
```

---

## Task 17: init 명령 구현

**Files:**
- Create: `scaffold/Home.md`
- Create: `scaffold/config.yml`
- Create: `scaffold/deploy.yml`
- Modify: `cmd/init_cmd.go`

- [ ] **Step 1: 스캐폴드 파일 작성**

```markdown
<!-- scaffold/Home.md -->
---
title: Home
titleKo: 위키 홈
---

# Home

> 환영합니다 :)

나의 개인 위키입니다. 마크다운으로 자유롭게 작성하세요.

## 시작하기

- 이 파일(`pages/Home.md`)을 수정하세요
- `pages/` 디렉토리에 새 `.md` 파일을 추가하면 위키 페이지가 됩니다
- `[[페이지 이름]]` 문법으로 페이지끼리 연결하세요
- `akwiki dev`로 로컬에서 미리보기, `akwiki build`로 정적 사이트 생성
```

```yaml
# scaffold/config.yml
site:
  title: "나의 위키"
  author: ""
  url: ""
  language: "ko"
```

```yaml
# scaffold/deploy.yml
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
          fetch-depth: 0
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

- [ ] **Step 2: scaffold을 embed으로 포함**

```go
// scaffold/embed.go
package scaffold

import "embed"

//go:embed Home.md config.yml deploy.yml
var Files embed.FS
```

- [ ] **Step 3: init 명령 구현**

```go
// cmd/init_cmd.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyungw00k/akwiki/scaffold"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create a new wiki",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}
		return runInit(dir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(dir string) error {
	dirs := []string{
		filepath.Join(dir, "pages"),
		filepath.Join(dir, "public"),
		filepath.Join(dir, ".akwiki"),
		filepath.Join(dir, ".github", "workflows"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	files := []struct {
		src string
		dst string
	}{
		{"Home.md", filepath.Join(dir, "pages", "Home.md")},
		{"config.yml", filepath.Join(dir, ".akwiki", "config.yml")},
		{"deploy.yml", filepath.Join(dir, ".github", "workflows", "deploy.yml")},
	}

	for _, f := range files {
		if _, err := os.Stat(f.dst); err == nil {
			fmt.Printf("  skip %s (already exists)\n", f.dst)
			continue
		}
		data, err := scaffold.Files.ReadFile(f.src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(f.dst, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("  create %s\n", f.dst)
	}

	fmt.Printf("\nWiki initialized in %s\n", dir)
	fmt.Println("Next steps:")
	fmt.Println("  cd " + dir)
	fmt.Println("  akwiki dev")
	return nil
}
```

- [ ] **Step 4: 테스트**

```bash
go build -o akwiki . && ./akwiki init /tmp/test-wiki && ls -R /tmp/test-wiki
```

Expected: 디렉토리 구조와 파일이 올바르게 생성됨.

- [ ] **Step 5: 커밋**

```bash
git add scaffold/ cmd/init_cmd.go
git commit -m "feat: add init command with scaffold files and GitHub Actions"
```

---

## Task 18: build 명령 연결

**Files:**
- Modify: `cmd/build.go`

- [ ] **Step 1: build 명령을 builder 패키지에 연결**

```go
// cmd/build.go
package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/kyungw00k/akwiki/internal/builder"
	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static site",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "."

		cfg, err := config.Load(rootDir)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		outDir := filepath.Join(rootDir, cfg.Build.OutDir)

		start := time.Now()
		fmt.Println("Building wiki...")

		if err := builder.Build(rootDir, outDir); err != nil {
			return err
		}

		fmt.Printf("Done in %s → %s/\n", time.Since(start).Round(time.Millisecond), outDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
```

- [ ] **Step 2: 통합 테스트**

```bash
cd /tmp/test-wiki
git init && git add . && git commit -m "init"
/Users/humphrey.park/Sandbox/akwiki/akwiki build
ls -R dist/
```

Expected: `dist/` 에 index.html, pages/Home/, assets/ 등 생성.

- [ ] **Step 3: 커밋**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
git add cmd/build.go
git commit -m "feat: wire build command to builder pipeline"
```

---

## Task 19: dev 명령 (개발 서버 + 라이브 리로드)

**Files:**
- Modify: `cmd/dev.go`

- [ ] **Step 1: fsnotify 의존성 추가**

```bash
go get github.com/fsnotify/fsnotify@latest
```

- [ ] **Step 2: dev 명령 구현**

```go
// cmd/dev.go
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kyungw00k/akwiki/internal/builder"
	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/spf13/cobra"
)

var devPort string

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with live reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "."
		cfg, err := config.Load(rootDir)
		if err != nil {
			return err
		}
		outDir := filepath.Join(rootDir, cfg.Build.OutDir)

		// Initial build
		fmt.Println("Building...")
		if err := builder.Build(rootDir, outDir); err != nil {
			return err
		}

		// Watch for changes
		go watchAndRebuild(rootDir, outDir)

		// Serve
		addr := ":" + devPort
		fmt.Printf("Serving at http://localhost%s\n", addr)
		return http.ListenAndServe(addr, http.FileServer(http.Dir(outDir)))
	},
}

func init() {
	devCmd.Flags().StringVarP(&devPort, "port", "p", "3000", "port to serve on")
	rootCmd.AddCommand(devCmd)
}

func watchAndRebuild(rootDir, outDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("watch error: %v", err)
		return
	}
	defer watcher.Close()

	watcher.Add(filepath.Join(rootDir, "pages"))
	watcher.Add(filepath.Join(rootDir, ".akwiki"))

	var debounce <-chan time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				debounce = time.After(300 * time.Millisecond)
			}
		case <-debounce:
			fmt.Println("Rebuilding...")
			start := time.Now()
			if err := builder.Build(rootDir, outDir); err != nil {
				log.Printf("build error: %v", err)
			} else {
				fmt.Printf("Rebuilt in %s\n", time.Since(start).Round(time.Millisecond))
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("watch error: %v", err)
		}
	}
}
```

- [ ] **Step 3: 수동 테스트**

```bash
go build -o akwiki . && cd /tmp/test-wiki && /Users/humphrey.park/Sandbox/akwiki/akwiki dev
# 브라우저에서 http://localhost:3000 확인 후 Ctrl+C
```

Expected: 위키 페이지가 브라우저에 표시됨.

- [ ] **Step 4: 커밋**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
git add cmd/dev.go go.mod go.sum
git commit -m "feat: add dev server with file watching and auto-rebuild"
```

---

## Task 20: serve 명령

**Files:**
- Modify: `cmd/serve.go`

- [ ] **Step 1: serve 명령 구현**

```go
// cmd/serve.go
package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/spf13/cobra"
)

var servePort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the built site",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "."
		cfg, err := config.Load(rootDir)
		if err != nil {
			return err
		}
		outDir := filepath.Join(rootDir, cfg.Build.OutDir)
		addr := ":" + servePort
		fmt.Printf("Serving %s at http://localhost%s\n", outDir, addr)
		return http.ListenAndServe(addr, http.FileServer(http.Dir(outDir)))
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "3000", "port to serve on")
	rootCmd.AddCommand(serveCmd)
}
```

- [ ] **Step 2: 커밋**

```bash
git add cmd/serve.go
git commit -m "feat: add serve command for previewing built site"
```

---

## Task 21: 위키링크 → private/missing 처리

**Files:**
- Modify: `internal/wiki/wikilink.go`
- Create: `internal/wiki/wikilink_resolve_test.go`

이 태스크는 위키링크 렌더링 시 private 페이지와 존재하지 않는 페이지를 처리합니다.

- [ ] **Step 1: 렌더링 컨텍스트 테스트 작성**

```go
// internal/wiki/wikilink_resolve_test.go
package wiki

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestWikilinkPrivate(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(NewWikilinkExtensionWithResolver("/pages", func(target string) LinkStatus {
			if target == "Secret" {
				return LinkPrivate
			}
			return LinkExists
		})),
	)
	var buf bytes.Buffer
	source := []byte("See [[Secret]] page.")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)
	md.Renderer().Render(&buf, source, doc)

	got := buf.String()
	want := `<p>See <a href="#private-link" data-link="private" class="wikilink">Secret</a> page.</p>`
	if got != want+"\n" {
		t.Errorf("\ngot:  %s\nwant: %s", got, want)
	}
}

func TestWikilinkMissing(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(NewWikilinkExtensionWithResolver("/pages", func(target string) LinkStatus {
			return LinkMissing
		})),
	)
	var buf bytes.Buffer
	source := []byte("See [[NonExistent]].")
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)
	md.Renderer().Render(&buf, source, doc)

	got := buf.String()
	want := `<p>See <a class="wikilink-missing">NonExistent</a>.</p>`
	if got != want+"\n" {
		t.Errorf("\ngot:  %s\nwant: %s", got, want)
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

```bash
go test ./internal/wiki/ -run "TestWikilinkPrivate|TestWikilinkMissing" -v
```

Expected: FAIL

- [ ] **Step 3: wikilink.go에 resolver 지원 추가**

```go
// wikilink.go에 추가할 타입과 함수

type LinkStatus int

const (
	LinkExists  LinkStatus = iota
	LinkPrivate
	LinkMissing
)

type LinkResolver func(target string) LinkStatus

type wikilinkResolverRenderer struct {
	pageRoute string
	resolver  LinkResolver
}

func (r *wikilinkResolverRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikilink, r.renderWikilink)
}

func (r *wikilinkResolverRenderer) renderWikilink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*Wikilink)

	status := LinkExists
	if r.resolver != nil {
		status = r.resolver(n.Target)
	}

	switch status {
	case LinkPrivate:
		_, _ = fmt.Fprintf(w, `<a href="#private-link" data-link="private" class="wikilink">%s</a>`, n.Display)
	case LinkMissing:
		_, _ = fmt.Fprintf(w, `<a class="wikilink-missing">%s</a>`, n.Display)
	default:
		href := fmt.Sprintf("%s/%s", r.pageRoute, url.PathEscape(n.Target))
		_, _ = fmt.Fprintf(w, `<a href="%s" class="wikilink">%s</a>`, href, n.Display)
	}
	return ast.WalkContinue, nil
}

type wikilinkExtensionWithResolver struct {
	pageRoute string
	resolver  LinkResolver
}

func NewWikilinkExtensionWithResolver(pageRoute string, resolver LinkResolver) goldmark.Extender {
	return &wikilinkExtensionWithResolver{pageRoute: pageRoute, resolver: resolver}
}

func (e *wikilinkExtensionWithResolver) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&wikilinkParser{}, 199),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&wikilinkResolverRenderer{
				pageRoute: e.pageRoute,
				resolver:  e.resolver,
			}, 199),
		),
	)
}
```

- [ ] **Step 4: 테스트 통과 확인**

```bash
go test ./internal/wiki/ -v
```

Expected: PASS (all tests)

- [ ] **Step 5: builder.go의 RenderMarkdown 호출에 resolver 연결**

`internal/render/markdown.go`에 resolver 파라미터 추가:

```go
func RenderMarkdownWithResolver(source []byte, pageRoute string, resolver wiki.LinkResolver) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			wiki.NewWikilinkExtensionWithResolver(pageRoute, resolver),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
```

`internal/builder/builder.go`에서 렌더링 시 resolver 전달:

```go
// Build 함수 내, 각 페이지 렌더링 부분
resolver := func(target string) wiki.LinkStatus {
    if privateNames[target] {
        return wiki.LinkPrivate
    }
    if _, ok := pageByName[target]; !ok {
        // Check aliases
        if resolved, ok := aliasMap[target]; ok {
            if _, ok := pageByName[resolved]; ok {
                return wiki.LinkExists
            }
        }
        return wiki.LinkMissing
    }
    return wiki.LinkExists
}

htmlContent, err := render.RenderMarkdownWithResolver(page.RawBody, pageRoute, resolver)
```

- [ ] **Step 6: 테스트 통과 확인**

```bash
go test ./... -v -timeout 30s
```

Expected: ALL PASS

- [ ] **Step 7: 커밋**

```bash
git add internal/wiki/wikilink.go internal/wiki/wikilink_resolve_test.go internal/render/markdown.go internal/builder/builder.go
git commit -m "feat: add private/missing link resolution in wikilinks"
```

---

## Task 22: 최종 통합 테스트 + 정리

**Files:**
- Modify: 필요 시 각 파일 미세 조정

- [ ] **Step 1: 전체 테스트 통과 확인**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
go test ./... -v -timeout 60s
```

Expected: ALL PASS

- [ ] **Step 2: 실제 위키로 E2E 테스트**

```bash
rm -rf /tmp/e2e-wiki
go build -o akwiki .
./akwiki init /tmp/e2e-wiki
cd /tmp/e2e-wiki

# 추가 페이지 생성
cat > pages/About.md << 'EOF'
---
title: About
titleKo: 소개
type: Article
---

# About

이 위키는 [[Home|홈]]에서 시작합니다.

## 참고

자세한 내용은 추후 추가됩니다.
EOF

cat > pages/Secret.md << 'EOF'
---
title: Secret Notes
private: true
---

# Secret

이 페이지는 보이지 않아야 합니다.
EOF

git init && git add . && git commit -m "init"
/Users/humphrey.park/Sandbox/akwiki/akwiki build

# 검증
test -f dist/pages/Home/index.html && echo "✓ Home page"
test -f dist/pages/About/index.html && echo "✓ About page"
test ! -f dist/pages/Secret/index.html && echo "✓ Secret excluded"
test -f dist/search-index.json && echo "✓ Search index"
test -f dist/pages/Home.txt && echo "✓ Raw markdown"
test -f dist/assets/style.css && echo "✓ CSS"
test -f dist/assets/search.js && echo "✓ Search JS"
grep "akngs" dist/pages/Home/index.html && echo "✓ Credits"
```

Expected: 모든 ✓ 출력.

- [ ] **Step 3: 크로스 컴파일 확인**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
GOOS=linux GOARCH=amd64 go build -o akwiki-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o akwiki-darwin-arm64 .
ls -lh akwiki-*
```

Expected: 두 바이너리 생성.

- [ ] **Step 4: 빌드 산출물 정리 + .gitignore**

```gitignore
# .gitignore
dist/
akwiki
akwiki-*
```

- [ ] **Step 5: 최종 커밋**

```bash
cd /Users/humphrey.park/Sandbox/akwiki
rm -f akwiki-linux-amd64 akwiki-darwin-arm64
git add .gitignore
git add -A
git commit -m "feat: akwiki v0.1.0 — personal wiki static site generator"
```
