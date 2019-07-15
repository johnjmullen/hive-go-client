package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var realmDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete realm pool",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		realm, err := restClient.GetRealm(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = realm.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	realmCmd.AddCommand(realmDeleteCmd)
	realmDeleteCmd.Flags().StringP("id", "i", "", "realm Pool Id")
	realmDeleteCmd.Flags().StringP("name", "n", "", "realm Pool Name")
}
