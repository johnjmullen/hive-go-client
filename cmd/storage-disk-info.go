package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageDiskInfoCmd = &cobra.Command{
	Use:   "disk-info [filename]",
	Short: "get information for a disk in a storage pool",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
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

		info, err := pool.DiskInfo(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString((info)))
	},
}

func init() {
	storageCmd.AddCommand(storageDiskInfoCmd)
	storageDiskInfoCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageDiskInfoCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
	addTaskFlags(storageDiskInfoCmd)
}
