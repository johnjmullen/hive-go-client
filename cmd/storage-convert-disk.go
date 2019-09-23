package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageConvertDiskCmd = &cobra.Command{
	Use:   "convert-disk",
	Short: "convert or copy a disk",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("src-storage")
		cmd.MarkFlagRequired("src-filename")
		cmd.MarkFlagRequired("dest-storage")
		cmd.MarkFlagRequired("dest-filename")
		viper.BindPFlag("src-storage", cmd.Flags().Lookup("src-storage"))
		viper.BindPFlag("src-filename", cmd.Flags().Lookup("src-filename"))
		viper.BindPFlag("dest-storage", cmd.Flags().Lookup("dest-storage"))
		viper.BindPFlag("dest-filename", cmd.Flags().Lookup("dest-filename"))
		viper.BindPFlag("dest-format", cmd.Flags().Lookup("dest-format"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		srcPool, err := restClient.GetStoragePoolByName(viper.GetString("src-storage"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		destPool, err := restClient.GetStoragePoolByName(viper.GetString("dest-storage"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = srcPool.ConvertDisk(restClient, viper.GetString("src-filename"), destPool.ID, viper.GetString("dest-filename"), viper.GetString("dest-format"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	storageCmd.AddCommand(storageConvertDiskCmd)
	storageConvertDiskCmd.Flags().String("src-storage", "", "Source storage pool name")
	storageConvertDiskCmd.Flags().String("src-filename", "", "Source filename")
	storageConvertDiskCmd.Flags().String("dest-storage", "", "Destination storage pool name")
	storageConvertDiskCmd.Flags().String("dest-filename", "", "Destination filename")
	storageConvertDiskCmd.Flags().String("dest-format", "qcow2", "Destination file format")
}
