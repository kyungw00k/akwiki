package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyungw00k/akwiki/internal/i18n"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:     "note [text]",
	Short:   i18n.T(i18n.MsgNoteShort),
	Aliases: []string{"n", "log"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := strings.TrimSpace(args[0])
		if text == "" {
			return fmt.Errorf("%s", i18n.T(i18n.ErrNoteNoText))
		}
		return runNote(".", text)
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}

func runNote(rootDir, text string) error {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	monthName := fmt.Sprintf("%d-%02d", year, month)
	dateStr := now.Format("2006-01-02")
	timeStr := now.Format("15:04")

	pagesDir := filepath.Join(rootDir, "pages")
	if err := os.MkdirAll(pagesDir, 0o755); err != nil {
		return err
	}

	monthPath := filepath.Join(pagesDir, monthName+".md")
	if err := upsertMonthPage(monthPath, monthName, dateStr, timeStr, text); err != nil {
		return err
	}

	yearPath := filepath.Join(pagesDir, fmt.Sprintf("Journal %d.md", year))
	if err := upsertYearPage(yearPath, year, month, monthName); err != nil {
		return err
	}

	fmt.Println(i18n.Tf(i18n.MsgNoteAdded, monthName))
	return nil
}

func upsertMonthPage(path, monthName, dateStr, timeStr, text string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		content := fmt.Sprintf("---\ntitle: %q\n---\n\n# %s\n\n## %s\n\n### %s\n\n%s\n",
			monthName, monthName, dateStr, timeStr, text)
		return os.WriteFile(path, []byte(content), 0o644)
	}

	content := string(data)
	dateHeader := "## " + dateStr
	timeHeader := "### " + timeStr
	entry := "\n\n" + timeHeader + "\n\n" + text + "\n"

	lines := strings.Split(content, "\n")
	h1Idx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "## ") {
			h1Idx = i
			break
		}
	}

	if strings.Contains(content, dateHeader) {
		// Date section exists — append to it
		dateIdx := strings.Index(content, dateHeader)
		insertIdx := dateIdx + len(dateHeader)
		remaining := content[insertIdx:]
		nextDate := strings.Index(remaining, "\n## ")
		if nextDate != -1 {
			insertIdx += nextDate
		} else {
			insertIdx = len(content)
		}
		content = content[:insertIdx] + entry + content[insertIdx:]
	} else if h1Idx != -1 {
		// New date — insert after h1 block
		insertPos := h1Idx + 1
		for insertPos < len(lines) && strings.TrimSpace(lines[insertPos]) == "" {
			insertPos++
		}
		prefix := strings.Join(lines[:insertPos], "\n")
		suffix := strings.Join(lines[insertPos:], "\n")
		content = prefix + "\n\n" + dateHeader + entry + "\n" + suffix
	} else {
		content += "\n\n" + dateHeader + entry
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

func upsertYearPage(path string, year, month int, monthName string) error {
	monthLabel := koreanMonth(month)
	linkLine := fmt.Sprintf("- [[%s|%s]]", monthName, monthLabel)

	data, err := os.ReadFile(path)
	if err != nil {
		prevYear := year - 1
		content := fmt.Sprintf("---\ntitle: Journal %d\naliases: [%d 일지]\n---\n\n# Journal %d\n\n%s\n\n← [[Journal %d]]",
			year, year, year, linkLine, prevYear)
		return os.WriteFile(path, []byte(content), 0o644)
	}

	content := string(data)
	if strings.Contains(content, linkLine) {
		return nil
	}

	lines := strings.Split(content, "\n")
	h1Idx := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			h1Idx = i
			break
		}
	}

	insertPos := h1Idx + 1
	if insertPos >= len(lines) {
		insertPos = len(lines)
	}
	// Skip blank lines after heading
	for insertPos < len(lines) && strings.TrimSpace(lines[insertPos]) == "" {
		insertPos++
	}

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:insertPos]...)
	newLines = append(newLines, linkLine)
	newLines = append(newLines, lines[insertPos:]...)

	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0o644)
}

func koreanMonth(month int) string {
	names := []string{"", "1월", "2월", "3월", "4월", "5월", "6월",
		"7월", "8월", "9월", "10월", "11월", "12월"}
	if month >= 1 && month <= 12 {
		return names[month]
	}
	return fmt.Sprintf("%d월", month)
}
