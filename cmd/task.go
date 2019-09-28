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

func init() {
	RootCmd.AddCommand(taskCmd)
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
