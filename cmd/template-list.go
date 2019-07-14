package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "list templates",
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := restClient.ListTemplates()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(templates))
		} else {
			for _, template := range templates {
				fmt.Println(template.Name)
			}
		}
	},
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	templateListCmd.Flags().Bool("details", false, "show details")
}
