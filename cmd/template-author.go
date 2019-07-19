package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateAuthorCmd = &cobra.Command{
	Use:   "author [Name]",
	Short: "author template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Author(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateAuthorCmd)
}
