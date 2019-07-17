package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var storageGetIDCmd = &cobra.Command{
	Use:   "get-id [name]",
	Short: "get storage pool id from name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := restClient.GetStoragePoolByName(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(pool.ID)
	},
}

func init() {
	storageCmd.AddCommand(storageGetIDCmd)
}
