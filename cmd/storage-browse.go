package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageBrowseCmd = &cobra.Command{
	Use:   "browse",
	Short: "list storage pool files",
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
	storageBrowseCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageBrowseCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
}
