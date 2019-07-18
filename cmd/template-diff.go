package cmd

import (
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

var templateDiffCmd = &cobra.Command{
	Use:   "diff [template1] [template2]",
	Short: "compare 2 templates",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		template1, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		template2, err := restClient.GetTemplate(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(cmp.Diff(template1, template2))
	},
}

func init() {
	templateCmd.AddCommand(templateDiffCmd)
}
