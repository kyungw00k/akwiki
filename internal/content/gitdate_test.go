package content

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestGitDates(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize a git repo
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test User"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = tmpDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("cmd %v failed: %v\n%s", args, err, out)
		}
	}

	// Create and commit a file
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	addCmd := exec.Command("git", "add", "test.md")
	addCmd.Dir = tmpDir
	if out, err := addCmd.CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v\n%s", err, out)
	}

	commitCmd := exec.Command("git", "commit", "-m", "initial commit")
	commitCmd.Dir = tmpDir
	if out, err := commitCmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v\n%s", err, out)
	}

	created, modified, err := GitDates(tmpDir, "test.md")
	if err != nil {
		t.Fatalf("GitDates() error = %v", err)
	}

	now := time.Now()
	if created.IsZero() {
		t.Error("created time is zero")
	}
	if modified.IsZero() {
		t.Error("modified time is zero")
	}
	// Times should be reasonably recent (within last minute)
	if created.After(now) {
		t.Errorf("created time %v is in the future", created)
	}
	if modified.After(now) {
		t.Errorf("modified time %v is in the future", modified)
	}
}

func TestGitDatesNoRepo(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a file without git
	testFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	created, modified, err := GitDates(tmpDir, "test.md")
	if err != nil {
		t.Fatalf("GitDates() error = %v", err)
	}

	// Should fallback to filesystem times (non-zero)
	if created.IsZero() {
		t.Error("created time is zero (expected filesystem fallback)")
	}
	if modified.IsZero() {
		t.Error("modified time is zero (expected filesystem fallback)")
	}
}
