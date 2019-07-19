package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var guestDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete guest pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		guest, err := restClient.GetGuest(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = guest.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	guestCmd.AddCommand(guestDeleteCmd)
}
