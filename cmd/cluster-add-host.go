package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rUser, rPass string
var addHostCmd = &cobra.Command{
	Use:   "add-host [ipAddress]",
	Short: "add a host to the cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := restClient.JoinHost(rUser, rPass, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	clusterCmd.AddCommand(addHostCmd)
	addHostCmd.Flags().StringVar(&rUser, "remote-username", "admin", "username for the remote host")
	addHostCmd.Flags().StringVar(&rPass, "remote-password", "admin", "password for the remote host")
}
