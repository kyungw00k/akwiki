package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with live reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting dev server...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
