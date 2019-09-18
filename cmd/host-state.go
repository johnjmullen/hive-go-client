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
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if viper.GetString("set") != "" {
			_, err := host.SetState(restClient, viper.GetString("set"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
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
}
