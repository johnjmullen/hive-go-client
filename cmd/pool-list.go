package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var poolListCmd = &cobra.Command{
	Use:   "list",
	Short: "list pools",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pools, err := restClient.ListGuestPools(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(pools))
		} else {
			for _, pool := range pools {
				fmt.Println(pool.Name)
			}
		}
	},
}

func init() {
	poolCmd.AddCommand(poolListCmd)
	poolListCmd.Flags().Bool("details", false, "show details")
	poolListCmd.Flags().String("filter", "", "filter query string")
}
