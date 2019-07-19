package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var templateAuthorCmd = &cobra.Command{
	Use:   "author [Name]",
	Short: "author template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = template.Author(restClient)
		if err != nil {
			log.Fatal(err)
		}
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateAuthorCmd)
}
