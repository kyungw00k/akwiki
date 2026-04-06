package content

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GitDates returns the creation and last modification times of a file
// based on its git commit history. Falls back to filesystem times if
// git is unavailable or the file is not tracked.
func GitDates(repoDir, relPath string) (created, modified time.Time, err error) {
	created, err = gitFirstCommitDate(repoDir, relPath)
	if err != nil || created.IsZero() {
		return fileDates(repoDir, relPath)
	}

	modified, err = gitLastCommitDate(repoDir, relPath)
	if err != nil || modified.IsZero() {
		return fileDates(repoDir, relPath)
	}

	return created, modified, nil
}

// gitFirstCommitDate returns the date of the first commit that introduced the file.
func gitFirstCommitDate(repoDir, relPath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "--diff-filter=A", "--follow", "--format=%aI", "--", relPath)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	// Take the last non-empty line (oldest commit)
	lines := strings.Split(strings.TrimSpace(string(bytes.TrimSpace(out))), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return time.Parse(time.RFC3339, line)
		}
	}
	return time.Time{}, nil
}

// gitLastCommitDate returns the date of the most recent commit that touched the file.
func gitLastCommitDate(repoDir, relPath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%aI", "--", relPath)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	line := strings.TrimSpace(string(out))
	if line == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, line)
}

// fileDates returns the filesystem modification time for both created and modified.
func fileDates(repoDir, relPath string) (created, modified time.Time, err error) {
	fullPath := filepath.Join(repoDir, relPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	t := info.ModTime()
	return t, t, nil
}
