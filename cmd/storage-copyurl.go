package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageCopyUrlCmd = &cobra.Command{
	Use:   "copy-url",
	Short: "copy a url to the storage pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("filename")
		cmd.MarkFlagRequired("url")
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("filename", cmd.Flags().Lookup("filename"))
		viper.BindPFlag("url", cmd.Flags().Lookup("url"))
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
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Println("\nDownloading " + viper.GetString("url"))
		}
		handleTask(pool.CopyUrl(restClient, viper.GetString("url"), viper.GetString("filename")))
	},
}

func init() {
	storageCmd.AddCommand(storageCopyUrlCmd)
	storageCopyUrlCmd.Flags().StringP("id", "i", "", "Storage Pool Id")
	storageCopyUrlCmd.Flags().StringP("name", "n", "", "Storage Pool Name")
	storageCopyUrlCmd.Flags().String("filename", "", "filename for the disk")
	storageCopyUrlCmd.Flags().String("url", "", "url to download")
	addTaskFlags(storageCopyUrlCmd)
}
