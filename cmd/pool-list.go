package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var poolListCmd = &cobra.Command{
	Use:   "list",
	Short: "list pools",
	Run: func(cmd *cobra.Command, args []string) {
		pools, err := restClient.ListGuestPools()
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
}
