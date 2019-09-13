package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var guestAssignCmd = &cobra.Command{
	Use:   "assign [GuestName]",
	Short: "assign guest to a user",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("guest-user")
		cmd.MarkFlagRequired("guest-realm")
	},
	Run: func(cmd *cobra.Command, args []string) {
		guest, err := restClient.GetGuest(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := restClient.AssignGuest(guest.PoolID, viper.GetString("guest-user"), viper.GetString("guest-realm"), guest.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(res))
	},
}

func init() {
	guestCmd.AddCommand(guestReleaseCmd)
	hostListCmd.Flags().String("guest-user", "", "user to assign to this guest")
	hostListCmd.Flags().String("guest-realm", "", "user's realm")

}
