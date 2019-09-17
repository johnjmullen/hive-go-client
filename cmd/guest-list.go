package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var guestListCmd = &cobra.Command{
	Use:   "list",
	Short: "list guests",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		guests, err := restClient.ListGuests(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(guests))
		} else {
			var guestList []string
			for _, guest := range guests {
				guestList = append(guestList, guest.Name)
			}
			fmt.Println(formatString(guestList))
		}
	},
}

func init() {
	guestCmd.AddCommand(guestListCmd)
	guestListCmd.Flags().Bool("details", false, "show details")
	guestListCmd.Flags().String("filter", "", "filter query string")
}
