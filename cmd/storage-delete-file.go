package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageDeleteFileCmd = &cobra.Command{
	Use:   "delete-file [file]",
	Short: "delete a file from the storage pool",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
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
		err = pool.DeleteFile(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	storageCmd.AddCommand(storageDeleteFileCmd)
	storageDeleteFileCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageDeleteFileCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
}
