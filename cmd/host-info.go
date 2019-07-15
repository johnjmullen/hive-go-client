package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var hostInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "hostid and version",
	Run: func(cmd *cobra.Command, args []string) {
		hostid, err := restClient.HostId()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		host, err := restClient.GetHost(hostid)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(host))
		} else {
			data := make(map[string]string)
			data["hostid"] = host.Hostid
			data["hostname"] = host.Hostname
			data["ip"] = host.IP
			fmt.Println(formatString(data))
		}
	},
}

func init() {
	hostCmd.AddCommand(hostInfoCmd)
	//hostListCmd.Flags().Bool("details", false, "show details")
}
