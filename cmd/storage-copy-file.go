package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageCopyFileCmd = &cobra.Command{
	Use:   "copy-file",
	Short: "copy a storage pool file",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired("srcStorageId")
		cmd.MarkFlagRequired("srcFilePath")
		cmd.MarkFlagRequired("destFilePath")
		cmd.MarkFlagRequired("destFilePath")
		viper.BindPFlag("srcStorageId", cmd.Flags().Lookup("srcStorageId"))
		viper.BindPFlag("srcFilePath", cmd.Flags().Lookup("srcFilePath"))
		viper.BindPFlag("destStorageId", cmd.Flags().Lookup("destStorageId"))
		viper.BindPFlag("destFilePath", cmd.Flags().Lookup("destFilePath"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleTask(restClient.CopyFile(viper.GetString("srcStorageId"), viper.GetString("srcFilePath"), viper.GetString("destStorageId"), viper.GetString("destFilePath")))
	},
}

func init() {
	storageCmd.AddCommand(storageCopyFileCmd)
	storageCopyFileCmd.Flags().String("srcStorageId", "", "Source storage pool id")
	storageCopyFileCmd.Flags().String("srcFilePath", "", "path to file in the source storage pool")
	storageCopyFileCmd.Flags().String("destStorageId", "", "Destination storage pool id")
	storageCopyFileCmd.Flags().String("destFilePath", "", "path to file in the destination storage pool")
	addTaskFlags(storageCopyFileCmd)
}
