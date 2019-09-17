package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostListCmd = &cobra.Command{
	Use:   "list",
	Short: "list hosts",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := restClient.ListHosts(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(hosts))
		} else {
			var list []map[string]string
			for _, host := range hosts {
				var hostInfo = map[string]string{"hostid": host.Hostid, "hostname": host.Hostname}
				list = append(list, hostInfo)
			}
			fmt.Println(formatString(list))
		}
	},
}

func init() {
	hostCmd.AddCommand(hostListCmd)
	hostListCmd.Flags().Bool("details", false, "show details")
	hostListCmd.Flags().String("filter", "", "filter query string")
}
