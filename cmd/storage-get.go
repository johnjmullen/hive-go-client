package cmd

import (
	"fmt"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get storage pool details",
	Run: func(cmd *cobra.Command, args []string) {
		var pool *rest.StoragePool
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			pool, err = restClient.GetStoragePool(viper.GetString("id"))
		case cmd.Flags().Changed("name"):
			pool, err = restClient.GetStoragePoolByName(viper.GetString("name"))
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

func init() {
	storageCmd.AddCommand(storageGetCmd)
	storageGetCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageGetCmd.Flags().StringP("name", "n", "", "Storage Pool Id")
	viper.BindPFlag("id", storageGetCmd.Flags().Lookup("id"))
	viper.BindPFlag("name", storageGetCmd.Flags().Lookup("name"))
}
