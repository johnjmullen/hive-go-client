package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/hive-io/hive-go-client/rest"
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
		newTask := task.WaitForTask(restClient, viper.GetBool("progress"))
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
	taskListCmd.Flags().Bool("details", false, "show details")
	taskListCmd.Flags().String("filter", "", "filter query string")

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
	if viper.GetBool("wait") {
		if viper.GetBool("raw-progress") {
			taskProgressNum(task)
		} else if viper.GetBool("progress-bar") {
			taskProgressBar(task)
		} else {
			taskVal := task.WaitForTask(restClient, false)
			if taskVal.State == "Complete" {
				fmt.Println(formatString("Task Complete"))
				os.Exit(1)
			}
			if taskVal.State == "failed" {
				fmt.Println(formatString("Task Failed: " + task.Message))
				os.Exit(1)
			}
		}
	} else {
		fmt.Println(formatString(map[string]string{"taskId": task.ID}))
	}
}

func taskProgressNum(task *rest.Task) {
	taskVal := task.WaitForTask(restClient, true)
	if taskVal.State == "failed" {
		fmt.Println(formatString("Task Failed: " + task.Message))
		os.Exit(1)
	}
}

func taskProgressBar(task *rest.Task) {
	uiprogress.Start()
	bar := uiprogress.AddBar(100)
	bar.AppendCompleted()
	done := make(chan struct{})
	taskData := make(chan rest.Task)
	var newVal rest.Task
	go task.WatchTask(restClient, taskData, done)
	for {
		select {
		case newVal = <-taskData:
			bar.Set(newVal.Progress)
		case <-done:
			if newVal.State == "completed" {
				bar.Set(100)
				uiprogress.Stop()
				time.Sleep(time.Millisecond * 100)
			} else if newVal.State == "failed" {
				bar.Set(0)
				fmt.Println(formatString("Task Failed: " + newVal.Message))
			}
			return
		}
	}
}
