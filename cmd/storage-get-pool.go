package cmd

import (
	"fmt"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var storageGetPoolCmd = &cobra.Command{
	Use:   "get-pool",
	Short: "get storage pool details",
	Run:   run,
}

var id, name string

func init() {
	storageCmd.AddCommand(storageGetPoolCmd)
	storageGetPoolCmd.Flags().StringVarP(&id, "id", "i", "", "Storage Pool Id")
	storageGetPoolCmd.Flags().StringVarP(&name, "name", "n", "", "Storage Pool Id")
}

func run(cmd *cobra.Command, args []string) {
	var pool *rest.StoragePool
	var err error
	switch {
	case cmd.Flags().Changed("id"):
		pool, err = restClient.GetStoragePool(id)
	case cmd.Flags().Changed("name"):
		pool, err = restClient.GetStoragePoolByName(name)
	default:
		cmd.Usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(formatString(pool))
}
