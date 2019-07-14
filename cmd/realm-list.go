package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var realmListCmd = &cobra.Command{
	Use:   "list",
	Short: "list realms",
	Run: func(cmd *cobra.Command, args []string) {
		realms, err := restClient.ListRealms()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(realms))
		} else {
			for _, realm := range realms {
				fmt.Println(realm.Name)
			}
		}
	},
}

func init() {
	realmCmd.AddCommand(realmListCmd)
	realmListCmd.Flags().Bool("details", false, "show details")
}
