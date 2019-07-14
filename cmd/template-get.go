package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateGetCmd = &cobra.Command{
	Use:   "get [Name]",
	Short: "get template details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		template, err := restClient.GetTemplate(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(template))
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateGetCmd)
}
