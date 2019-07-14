package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var hostListCmd = &cobra.Command{
	Use:   "list",
	Short: "list hosts",
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := restClient.ListHosts()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(hosts))
		} else {
			for _, host := range hosts {
				fmt.Printf("hostid: %s\t hostname: %s\n", host.Hostid, host.Hostname)
			}
		}
	},
}

func init() {
	hostCmd.AddCommand(hostListCmd)
	hostListCmd.Flags().Bool("details", false, "show details")
}
