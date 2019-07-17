package cmd

import (
	"fmt"
	"os"

	rest "github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get profile details",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			profile, err = restClient.GetProfile(viper.GetString("id"))
		case cmd.Flags().Changed("name"):
			profile, err = restClient.GetProfileByName(viper.GetString("name"))
		default:
			cmd.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(profile))
	},
}

func init() {
	profileCmd.AddCommand(profileGetCmd)
	profileGetCmd.Flags().StringP("id", "i", "", "profile id")
	profileGetCmd.Flags().StringP("name", "n", "", "profile name")
}
