package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateUnloadCmd = &cobra.Command{
	Use:   "unload [Name]",
	Short: "unload template from all hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Unload(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateUnloadCmd)
}
