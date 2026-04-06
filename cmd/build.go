package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static site",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Building wiki...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
