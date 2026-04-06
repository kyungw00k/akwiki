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
