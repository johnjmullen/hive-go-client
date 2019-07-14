package cmd

import (
	"fmt"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
	"github.com/spf13/cobra"
)

var storageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get storage pool details",
	Run: func(cmd *cobra.Command, args []string) {
		var pool *rest.StoragePool
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			pool, err = restClient.GetStoragePool(poolId)
		case cmd.Flags().Changed("name"):
			pool, err = restClient.GetStoragePoolByName(poolName)
		default:
			cmd.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(pool))
	},
}

var poolId, poolName string

func init() {
	storageCmd.AddCommand(storageGetCmd)
	storageGetCmd.Flags().StringVarP(&poolId, "id", "i", "", "Storage Pool Id")
	storageGetCmd.Flags().StringVarP(&poolName, "name", "n", "", "Storage Pool Id")
}
