package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var realmListCmd = &cobra.Command{
	Use:   "list",
	Short: "list realms",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		realms, err := restClient.ListRealms(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(realms))
		} else {
			var list []string
			for _, realm := range realms {
				list = append(list, realm.Name)
			}
			fmt.Println(formatString(list))
		}
	},
}

func init() {
	realmCmd.AddCommand(realmListCmd)
	realmListCmd.Flags().Bool("details", false, "show details")
	realmListCmd.Flags().String("filter", "", "filter query string")
}
