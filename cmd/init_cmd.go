package cmd

import (
	"fmt"
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
		fmt.Printf("Initializing wiki in %s\n", dir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
