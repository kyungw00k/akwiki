package builder

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	// Create temp directory
	dir := t.TempDir()

	// Create pages directory
	pagesDir := filepath.Join(dir, "pages")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create .akwiki directory for config
	if err := os.MkdirAll(filepath.Join(dir, ".akwiki"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Home.md with frontmatter and wikilink to About
	homeContent := `---
title: Home
titleKo: 홈
---
# Home

Welcome to my wiki!

Check out [[About]] for more info.
`
	if err := os.WriteFile(filepath.Join(pagesDir, "Home.md"), []byte(homeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// About.md with wikilink back to Home
	aboutContent := `---
title: About
---
# About

This is the about page. Go back to [[Home]].
`
	if err := os.WriteFile(filepath.Join(pagesDir, "About.md"), []byte(aboutContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Secret.md (private page)
	secretContent := `---
title: Secret
private: true
---
# Secret

This page is private.
`
	if err := os.WriteFile(filepath.Join(pagesDir, "Secret.md"), []byte(secretContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Init git repo, add, and commit (needed for GitDates)
	for _, args := range [][]string{
		{"init"},
		{"add", "."},
		{"-c", "user.name=test", "-c", "user.email=test@test.com", "commit", "-m", "init"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %s\n%s", args, err, out)
		}
	}

	// Run Build
	outDir := filepath.Join(dir, "dist")
	if err := Build(dir, outDir); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify: index.html exists
	assertFileExists(t, filepath.Join(outDir, "index.html"))

	// Verify: pages/Home/index.html exists
	assertFileExists(t, filepath.Join(outDir, "pages", "Home", "index.html"))

	// Verify: pages/About/index.html exists
	assertFileExists(t, filepath.Join(outDir, "pages", "About", "index.html"))

	// Verify: pages/Home.txt exists
	assertFileExists(t, filepath.Join(outDir, "pages", "Home.txt"))

	// Verify: search-index.json exists
	assertFileExists(t, filepath.Join(outDir, "search-index.json"))

	// Verify: assets/style.css exists
	assertFileExists(t, filepath.Join(outDir, "assets", "style.css"))

	// Verify: assets/search.js exists
	assertFileExists(t, filepath.Join(outDir, "assets", "search.js"))

	// Verify: Secret page NOT in output
	secretDir := filepath.Join(outDir, "pages", "Secret")
	if _, err := os.Stat(secretDir); !os.IsNotExist(err) {
		t.Error("Secret page should not be in output, but directory exists")
	}
	secretTxt := filepath.Join(outDir, "pages", "Secret.txt")
	if _, err := os.Stat(secretTxt); !os.IsNotExist(err) {
		t.Error("Secret.txt should not be in output, but file exists")
	}

	// Verify: Home HTML is non-empty
	homeHTML, err := os.ReadFile(filepath.Join(outDir, "pages", "Home", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if len(homeHTML) == 0 {
		t.Error("Home HTML should be non-empty")
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}
