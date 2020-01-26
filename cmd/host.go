package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "host operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var hostGetCmd = &cobra.Command{
	Use:   "get [hostid]",
	Short: "get host details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(host))
	},
}

var hostInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "hostid and version",
	Run: func(cmd *cobra.Command, args []string) {
		hostid, err := restClient.HostID()
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
			list := []map[string]string{}
			for _, host := range hosts {
				var hostInfo = map[string]string{"hostid": host.Hostid, "hostname": host.Hostname}
				list = append(list, hostInfo)
			}
			fmt.Println(formatString(list))
		}
	},
}

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
			_, err := host.UpdateAppliance(restClient)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println(formatString(host.Appliance.Loglevel))
		}

	},
}

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

var hostUnjoinCmd = &cobra.Command{
	Use:   "unjoin [hostid]",
	Short: "remove host from cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.UnjoinCluster(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

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
	RootCmd.AddCommand(hostCmd)
	hostCmd.AddCommand(hostGetCmd)
	hostCmd.AddCommand(hostInfoCmd)

	hostCmd.AddCommand(hostListCmd)
	hostListCmd.Flags().Bool("details", false, "show details")
	hostListCmd.Flags().String("filter", "", "filter query string")

	hostCmd.AddCommand(hostLogLevelCmd)
	hostLogLevelCmd.Flags().StringP("set", "s", "", "set log level (error/warn/info/debug)")

	hostCmd.AddCommand(hostRestartServicesCmd)
	hostCmd.AddCommand(hostUnjoinCmd)

	hostCmd.AddCommand(hostStateCmd)
	hostStateCmd.Flags().StringP("set", "s", "", "set host state (available/maintenance)")
	addTaskFlags(hostStateCmd)
}
