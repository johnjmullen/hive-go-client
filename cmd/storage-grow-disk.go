package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageGrowDiskCmd = &cobra.Command{
	Use:   "grow-disk",
	Short: "grow a disk in the storage pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("filename", cmd.Flags().Lookup("filename"))
		viper.BindPFlag("disk-size", cmd.Flags().Lookup("disk-size"))
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
		err = pool.GrowDisk(restClient, viper.GetString("filename"), uint(viper.GetInt("disk-size")))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	storageCmd.AddCommand(storageGrowDiskCmd)
	storageGrowDiskCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageGrowDiskCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
	storageGrowDiskCmd.Flags().String("filename", "", "filename for the disk")
	storageGrowDiskCmd.Flags().Int("disk-size", 0, "size to add in GB")
}
