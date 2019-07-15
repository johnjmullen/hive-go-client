package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var hostGetCmd = &cobra.Command{
	Use:   "get [hostid]",
	Short: "get host details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		host, err := restClient.GetHost(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(host))
	},
}

//var name string

func init() {
	hostCmd.AddCommand(hostGetCmd)
}
