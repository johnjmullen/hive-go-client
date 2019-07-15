package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var storageBrowseCmd = &cobra.Command{
	Use:   "browse [storagePoolId]",
	Short: "list storage pool contents",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := restClient.GetStoragePool(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		files, err := pool.Browse(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(files))
	},
}

func init() {
	storageCmd.AddCommand(storageBrowseCmd)
}
