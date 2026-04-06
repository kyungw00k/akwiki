package cli

import (
	"github.com/kyungw00k/akwiki/internal/i18n"
	"github.com/spf13/cobra"
)

var Version = "dev" // injected via ldflags

var rootCmd = &cobra.Command{
	Use:     "akwiki",
	Short:   i18n.T(i18n.MsgRootShort),
	Long:    i18n.T(i18n.MsgRootLong),
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.SuggestionsMinimumDistance = 2
}

func Execute() error {
	return rootCmd.Execute()
}
