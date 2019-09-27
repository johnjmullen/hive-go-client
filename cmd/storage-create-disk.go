package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageCreateDiskCmd = &cobra.Command{
	Use:   "create-disk",
	Short: "create a disk in the storage pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("filename", cmd.Flags().Lookup("filename"))
		viper.BindPFlag("disk-format", cmd.Flags().Lookup("disk-format"))
		viper.BindPFlag("disk-size", cmd.Flags().Lookup("disk-size"))
		bindTaskFlags(cmd)
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
		handleTask(pool.CreateDisk(restClient, viper.GetString("filename"), viper.GetString("disk-format"), uint(viper.GetInt("disk-size"))))
	},
}

func init() {
	storageCmd.AddCommand(storageCreateDiskCmd)
	storageCreateDiskCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageCreateDiskCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
	storageCreateDiskCmd.Flags().String("filename", "", "filename for the disk")
	storageCreateDiskCmd.Flags().String("disk-format", "qcow2", "disk format ()")
	storageCreateDiskCmd.Flags().Int("disk-size", 25, "size of the disk in GB")
	addTaskFlags(storageCreateDiskCmd)
}
