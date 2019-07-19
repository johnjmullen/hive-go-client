package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docCmd = &cobra.Command{
	Use:   "doc [directory]",
	Args:  cobra.ExactArgs(1),
	Short: "Generates documentation",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(RootCmd, args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	utilCmd.AddCommand(docCmd)
}
