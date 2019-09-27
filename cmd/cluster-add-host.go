package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rUser, rPass string
var addHostCmd = &cobra.Command{
	Use:   "add-host [ipAddress]",
	Short: "add a host to the cluster",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Printf("Adding %s to cluster\n", args[0])
		}
		handleTask(restClient.JoinHost(rUser, rPass, args[0]))
	},
}

func init() {
	clusterCmd.AddCommand(addHostCmd)
	addHostCmd.Flags().StringVar(&rUser, "remote-username", "admin", "username for the remote host")
	addHostCmd.Flags().StringVar(&rPass, "remote-password", "admin", "password for the remote host")
	addTaskFlags(addHostCmd)
}
