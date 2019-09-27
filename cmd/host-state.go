package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostStateCmd = &cobra.Command{
	Use:   "state [hostid]",
	Short: "get or set host state",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("set", cmd.Flags().Lookup("set"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if viper.GetString("set") != "" {
			if viper.GetBool("wait") && viper.GetBool("progress-bar") {
				fmt.Printf("Setting state on %s to %s\n", host.Hostname, viper.GetString("set"))
			}
			handleTask(host.SetState(restClient, viper.GetString("set")))
		} else {
			state, err := host.GetState(restClient)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(formatString(state))
		}

	},
}

func init() {
	hostCmd.AddCommand(hostStateCmd)
	hostStateCmd.Flags().StringP("set", "s", "", "set host state (available/maintenance)")
	addTaskFlags(hostStateCmd)
}
