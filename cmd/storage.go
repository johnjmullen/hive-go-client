package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: `storage operations`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(storageCmd)
}
