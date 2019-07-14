package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := restClient.ListProfiles()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(profiles))
		} else {
			for _, profile := range profiles {
				fmt.Println(profile.Name)
			}
		}
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileListCmd.Flags().Bool("details", false, "show details")
}
