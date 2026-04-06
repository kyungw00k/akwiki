package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "akwiki",
	Short: "Personal wiki static site generator",
	Long:  "akwiki generates static wiki sites from markdown files.\nInspired by akngs's wiki (https://wiki.g15e.com).",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
