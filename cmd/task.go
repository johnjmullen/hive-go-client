package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "task operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

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
	},
}

var taskGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get task details",
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
		fmt.Println(formatString(task))
	},
}

// listCmd represents the list command
var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "list tasks",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := restClient.ListTasks(listFlagsToQuery())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(tasks))
		} else {
			list := []map[string]string{}
			for _, task := range tasks {
				var info = map[string]string{"id": task.ID, "name": task.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

var taskWaitCmd = &cobra.Command{
	Use:   "wait",
	Short: "wait for a task to complete",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("progress", cmd.Flags().Lookup("progress"))
		viper.BindPFlag("progress-bar", cmd.Flags().Lookup("progress-bar"))
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
		newTask, err := task.WaitForTask(restClient, viper.GetBool("progress"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if newTask.State == "completed" {
			fmt.Println(formatString("Task Complete"))
		}
		if newTask.State == "failed" {
			fmt.Println(formatString("Task Failed: " + task.Message))
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskForceCompleteCmd)
	taskForceCompleteCmd.Flags().StringP("id", "i", "", "task id")
	taskForceCompleteCmd.Flags().StringP("name", "n", "", "task name")

	taskCmd.AddCommand(taskGetCmd)
	taskGetCmd.Flags().StringP("id", "i", "", "task id")
	taskGetCmd.Flags().StringP("name", "n", "", "task name")

	taskCmd.AddCommand(taskListCmd)
	addListFlags(taskListCmd)

	taskCmd.AddCommand(taskWaitCmd)
	taskWaitCmd.Flags().StringP("id", "i", "", "task id")
	taskWaitCmd.Flags().StringP("name", "n", "", "task name")
	taskWaitCmd.Flags().Bool("progress", false, "print progress")
	taskWaitCmd.Flags().Bool("progress-bar", false, "print progress-bar")
}

func addTaskFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("wait", false, "wait for task to complete")
	cmd.Flags().Bool("raw-progress", false, "print progress as a number with --wait")
	cmd.Flags().Bool("progress-bar", false, "show a progress bar with --wait")
}
func bindTaskFlags(cmd *cobra.Command) {
	viper.BindPFlag("wait", cmd.Flags().Lookup("wait"))
	viper.BindPFlag("raw-progress", cmd.Flags().Lookup("raw-progress"))
	viper.BindPFlag("progress-bar", cmd.Flags().Lookup("progress-bar"))
}

func handleTask(task *rest.Task, err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if task == nil {
		fmt.Println("Error reading task")
		os.Exit(1)
	}
	if viper.GetBool("wait") {
		if err := waitForTask(task, viper.GetBool("raw-progress"), viper.GetBool("progress-bar")); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !viper.GetBool("raw-progress") && !viper.GetBool("progress-bar") {
			fmt.Println(formatString("Task Complete"))
		}
	} else {
		fmt.Println(formatString(map[string]string{"taskId": task.ID}))
	}

}

func waitForTask(task *rest.Task, rawProgress, progressBar bool) error {
	if task == nil {
		return fmt.Errorf("error reading task")
	}
	if progressBar {
		taskProgressBar(task)
	} else {
		taskVal, err := task.WaitForTask(restClient, rawProgress)
		if err != nil {
			return err
		}
		if taskVal.State == "failed" {
			return fmt.Errorf("%s", formatString("Task Failed: "+task.Message))
		}
	}
	return nil
}

func taskProgressBar(task *rest.Task) error {
	bar := progressbar.NewOptions(100,
		progressbar.OptionFullWidth(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetDescription(task.Description))
	errChannel := make(chan error)
	taskData := make(chan rest.Task)
	var newVal rest.Task
	bar.Set(int(task.Progress))
	go task.WatchTask(restClient, taskData, errChannel)
	for {
		select {
		case newVal = <-taskData:
			bar.Set(int(newVal.Progress))
			if newVal.State == "completed" {
				bar.Set(100)
				bar.Finish()
				fmt.Println("")
				return nil
			} else if newVal.State == "failed" {
				bar.Set(0)
				return fmt.Errorf("%s", formatString("Task Failed: "+newVal.Message))
			}
		case err := <-errChannel:
			return err
		}
	}
}
