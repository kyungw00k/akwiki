package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the built site",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Serving dist/...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
