package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var taskForceCompleteCmd = &cobra.Command{
	Use:   "force-complete",
	Short: "force task state to completed",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var task *rest.Task
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			task, err = restClient.GetTask(viper.GetString("id"))
		case cmd.Flags().Changed("name"):
			task, err = restClient.GetTaskByName(viper.GetString("name"))
		default:
			cmd.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = task.ForceComplete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println(formatString(task))
	},
}

func init() {
	taskCmd.AddCommand(taskForceCompleteCmd)
	taskForceCompleteCmd.Flags().StringP("id", "i", "", "task id")
	taskForceCompleteCmd.Flags().StringP("name", "n", "", "task name")
}
