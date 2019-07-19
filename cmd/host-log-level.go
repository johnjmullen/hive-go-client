package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostLogLevelCmd = &cobra.Command{
	Use:   "log-level [hostid]",
	Short: "get or set host log level",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("set", cmd.Flags().Lookup("set"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if viper.GetString("set") != "" {
			host.Appliance.Loglevel = viper.GetString("set")
			_, err := host.Update(restClient)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println(formatString(host.Appliance.Loglevel))
		}

	},
}

func init() {
	hostCmd.AddCommand(hostLogLevelCmd)
	hostLogLevelCmd.Flags().StringP("set", "s", "", "set log level (error/warn/info/debug)")
}
