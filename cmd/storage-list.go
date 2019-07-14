package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var storageListCmd = &cobra.Command{
	Use:   "list",
	Short: "list storage pools",
	Run: func(cmd *cobra.Command, args []string) {
		pools, err := restClient.ListStoragePools()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(pools))
		} else {
			for _, pool := range pools {
				fmt.Printf("ID: %s\t Name: %s\n", pool.ID, pool.Name)
			}
		}
	},
}

func init() {
	storageCmd.AddCommand(storageListCmd)
	storageListCmd.Flags().Bool("details", false, "show details")
}
