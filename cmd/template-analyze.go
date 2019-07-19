package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var templateAnalyzeCmd = &cobra.Command{
	Use:   "analyze [Name]",
	Short: "analyze template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = template.Analyze(restClient)
		if err != nil {
			log.Fatal(err)
		}
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateAnalyzeCmd)
}
