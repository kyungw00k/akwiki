package cli

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/i18n"
	"github.com/spf13/cobra"
)

var servePort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: i18n.T(i18n.MsgServeShort),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir := "."
		cfg, err := config.Load(rootDir)
		if err != nil {
			return fmt.Errorf(i18n.T(i18n.ErrConfigLoad), err)
		}
		outDir := filepath.Join(rootDir, cfg.Build.OutDir)
		addr := ":" + servePort
		fmt.Println(i18n.Tf(i18n.MsgServeServing, outDir, addr))
		return http.ListenAndServe(addr, http.FileServer(http.Dir(outDir)))
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "3000", i18n.T(i18n.FlagPortUsage))
	rootCmd.AddCommand(serveCmd)
}
