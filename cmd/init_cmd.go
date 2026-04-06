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

	files := []struct{ src, dst string }{
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
