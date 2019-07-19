package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var templateUnloadCmd = &cobra.Command{
	Use:   "unload [Name]",
	Short: "unload template from all hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = template.Unload(restClient)
		if err != nil {
			log.Fatal(err)
		}
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateUnloadCmd)
}
