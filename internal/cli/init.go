package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyungw00k/akwiki/internal/i18n"
	"github.com/kyungw00k/akwiki/scaffold"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: i18n.T(i18n.MsgInitShort),
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

	files := []struct{ src, dst string }{
		{"Home.md", filepath.Join(dir, "pages", "Home.md")},
		{"config.yml", filepath.Join(dir, ".akwiki", "config.yml")},
		{"deploy.yml", filepath.Join(dir, ".github", "workflows", "deploy.yml")},
		{"AGENTS.md", filepath.Join(dir, "AGENTS.md")},
	}
	for _, f := range files {
		if _, err := os.Stat(f.dst); err == nil {
			fmt.Println(i18n.Tf(i18n.MsgInitSkip, f.dst))
			continue
		}
		data, err := scaffold.Files.ReadFile(f.src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(f.dst, data, 0o644); err != nil {
			return err
		}
		fmt.Println(i18n.Tf(i18n.MsgInitCreate, f.dst))
	}

	// Create CLAUDE.md symlink to AGENTS.md
	claudePath := filepath.Join(dir, "CLAUDE.md")
	if _, err := os.Lstat(claudePath); err != nil {
		if err := os.Symlink("AGENTS.md", claudePath); err != nil {
			return err
		}
		fmt.Println(i18n.Tf(i18n.MsgInitCreate, claudePath))
	} else {
		fmt.Println(i18n.Tf(i18n.MsgInitSkip, claudePath))
	}

	fmt.Println(i18n.Tf(i18n.MsgInitDone, dir))
	fmt.Println(i18n.Tf(i18n.MsgInitNext, dir))
	return nil
}
