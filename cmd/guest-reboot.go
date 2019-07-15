package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var guestRebootCmd = &cobra.Command{
	Use:   "reboot [Name]",
	Short: "reboot guest",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		guest, err := restClient.GetGuest(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		guest.Reboot(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//var name string

func init() {
	guestCmd.AddCommand(guestRebootCmd)
}
