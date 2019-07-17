package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var hostRestartServicesCmd = &cobra.Command{
	Use:   "restart-services [hostid]",
	Short: "restart hive servies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.RestartServices(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//var name string

func init() {
	hostCmd.AddCommand(hostRestartServicesCmd)
}
