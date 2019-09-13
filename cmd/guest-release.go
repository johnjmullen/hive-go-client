package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var guestReleaseCmd = &cobra.Command{
	Use:   "release [Name]",
	Short: "release guest assignment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		guest, err := restClient.GetGuest(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = restClient.ReleaseGuest(guest.PoolID, guest.Username, guest.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	guestCmd.AddCommand(guestReleaseCmd)
}
