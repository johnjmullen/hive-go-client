package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var realmCmd = &cobra.Command{
	Use:   "realm",
	Short: "realm operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(realmCmd)
}
