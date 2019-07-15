package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete template pool",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	templateCmd.AddCommand(templateDeleteCmd)
	templateDeleteCmd.Flags().StringP("id", "i", "", "template Pool Id")
	templateDeleteCmd.Flags().StringP("name", "n", "", "template Pool Name")
}
