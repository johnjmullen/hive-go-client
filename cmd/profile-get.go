package cmd

import (
	"fmt"
	"os"

	rest "github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get profile details",
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			profile, err = restClient.GetProfile(id)
		case cmd.Flags().Changed("name"):
			profile, err = restClient.GetProfileByName(name)
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

var id, name string

func init() {
	profileCmd.AddCommand(profileGetCmd)
	profileGetCmd.Flags().StringVarP(&id, "id", "i", "", "Storage Pool Id")
	profileGetCmd.Flags().StringVarP(&name, "name", "n", "", "Storage Pool Id")
}
