package cmd

import (
	"fmt"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
	"github.com/spf13/cobra"
)

var storageDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete storage pool",
	Run: func(cmd *cobra.Command, args []string) {
		var pool *rest.StoragePool
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			id, _ := cmd.Flags().GetString("id")
			pool, err = restClient.GetStoragePool(id)
		case cmd.Flags().Changed("name"):
			name, _ := cmd.Flags().GetString("name")
			pool, err = restClient.GetStoragePoolByName(name)
		default:
			cmd.Usage()
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = pool.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	storageCmd.AddCommand(storageDeleteCmd)
	storageDeleteCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageDeleteCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
}
