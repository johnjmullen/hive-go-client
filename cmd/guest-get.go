package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var guestGetCmd = &cobra.Command{
	Use:   "get [Name]",
	Short: "get guest details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		guest, err := restClient.GetGuest(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(guest))
	},
}

//var name string

func init() {
	guestCmd.AddCommand(guestGetCmd)
}
