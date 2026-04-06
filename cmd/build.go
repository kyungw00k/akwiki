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
