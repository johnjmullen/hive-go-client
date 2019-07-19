package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var templateLoadCmd = &cobra.Command{
	Use:   "load [Name]",
	Short: "load template to all hosts",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("storage", cmd.Flags().Lookup("storage"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Load(restClient, viper.GetString("storage"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//var name string

func init() {
	templateCmd.AddCommand(templateLoadCmd)
	templateLoadCmd.Flags().StringP("storage", "s", "disk", "Location to load the template (disk or ram)")
}
