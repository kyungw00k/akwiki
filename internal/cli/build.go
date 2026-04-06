package cli

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/kyungw00k/akwiki/internal/builder"
	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/i18n"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: i18n.T(i18n.MsgBuildShort),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "."
		cfg, err := config.Load(rootDir)
		if err != nil {
			return fmt.Errorf(i18n.T(i18n.ErrConfigLoad), err)
		}
		outDir := filepath.Join(rootDir, cfg.Build.OutDir)
		start := time.Now()
		fmt.Println(i18n.T(i18n.MsgBuildBuilding))
		if err := builder.Build(rootDir, outDir); err != nil {
			return fmt.Errorf(i18n.T(i18n.ErrBuildFail), err)
		}
		fmt.Println(i18n.Tf(i18n.MsgBuildDone, time.Since(start).Round(time.Millisecond), outDir))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
