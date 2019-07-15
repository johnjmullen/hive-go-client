package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var realmGetCmd = &cobra.Command{
	Use:   "get [Name]",
	Short: "get realm details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm, err := restClient.GetRealm(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(realm))
	},
}

//var name string

func init() {
	realmCmd.AddCommand(realmGetCmd)
}
