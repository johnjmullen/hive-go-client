package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var guestListCmd = &cobra.Command{
	Use:   "list",
	Short: "list guests",
	Run: func(cmd *cobra.Command, args []string) {
		guests, err := restClient.ListGuests()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(guests))
		} else {
			for _, guest := range guests {
				fmt.Println(guest.Name)
			}
		}
	},
}

func init() {
	guestCmd.AddCommand(guestListCmd)
	guestListCmd.Flags().Bool("details", false, "show details")
	viper.BindPFlag("details", guestListCmd.Flags().Lookup("details"))
}
