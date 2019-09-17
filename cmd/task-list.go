package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "list tasks",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := restClient.ListTasks(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(tasks))
		} else {
			var list []map[string]string
			for _, task := range tasks {
				var info = map[string]string{"id": task.ID, "name": task.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

func init() {
	taskCmd.AddCommand(taskListCmd)
	taskListCmd.Flags().Bool("details", false, "show details")
	taskListCmd.Flags().String("filter", "", "filter query string")
}
