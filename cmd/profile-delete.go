package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete profile profile",
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			id, _ := cmd.Flags().GetString("id")
			profile, err = restClient.GetProfile(id)
		case cmd.Flags().Changed("name"):
			name, _ := cmd.Flags().GetString("name")
			profile, err = restClient.GetProfileByName(name)
		default:
			cmd.Usage()
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = profile.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	profileCmd.AddCommand(profileDeleteCmd)
	profileDeleteCmd.Flags().StringP("id", "i", "", "profile profile Id")
	profileDeleteCmd.Flags().StringP("name", "n", "", "profile profile Name")
}
