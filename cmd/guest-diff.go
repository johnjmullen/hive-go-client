package cmd

import (
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

var guestDiffCmd = &cobra.Command{
	Use:   "diff [guest1] [guest2]",
	Short: "compare 2 guests",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		guest1, err := restClient.GetGuest(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		guest2, err := restClient.GetGuest(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(cmp.Diff(guest1, guest2))
	},
}

func init() {
	guestCmd.AddCommand(guestDiffCmd)
}
