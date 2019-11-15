package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageMoveFileCmd = &cobra.Command{
	Use:   "move-file",
	Short: "move a storage pool file",
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
		handleTask(restClient.MoveFile(viper.GetString("srcStorageId"), viper.GetString("srcFilePath"), viper.GetString("destStorageId"), viper.GetString("destFilePath")))
	},
}

func init() {
	storageCmd.AddCommand(storageMoveFileCmd)
	storageMoveFileCmd.Flags().String("srcStorageId", "", "Source storage pool id")
	storageMoveFileCmd.Flags().String("srcFilePath", "", "path to file in the source storage pool")
	storageMoveFileCmd.Flags().String("destStorageId", "", "Destination storage pool id")
	storageMoveFileCmd.Flags().String("destFilePath", "", "path to file in the destination storage pool")
	addTaskFlags(storageMoveFileCmd)
}
