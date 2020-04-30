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
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := restClient.ListHosts(listFlagsToQuery())
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

var hostListSoftwareCmd = &cobra.Command{
	Use:   "list-software [hostid]",
	Short: "list available software packages on a host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		software, err := host.ListSoftware(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(software))
	},
}

var hostUploadSoftware = &cobra.Command{
	Use:   "upload-software [file]",
	Short: "upload a software pkg file to a host",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		hostid, err := restClient.HostID()
		if err != nil {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		host, err := restClient.GetHost(hostid)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.UploadSoftware(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostDeploySoftware = &cobra.Command{
	Use:   "deploy-software [hostid]",
	Short: "deploy a software package",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))
		cmd.MarkFlagRequired("package")
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		name := viper.GetString("package")
		handleTask(host.DeploySoftware(restClient, name))
	},
}

var hostDeleteSoftware = &cobra.Command{
	Use:   "delete-software [hostid]",
	Short: "delete a software package",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))
		cmd.MarkFlagRequired("package")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		name := viper.GetString("package")
		err = host.DeleteSoftware(restClient, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostListFirmwareCmd = &cobra.Command{
	Use:   "list-firmware [hostid]",
	Short: "list available firmware images on a host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		firmware, err := host.ListFirmware(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(firmware))
	},
}

var hostUploadFirmware = &cobra.Command{
	Use:   "upload-firmware [file]",
	Args:  cobra.ExactArgs(1),
	Short: "upload a firmware pkg file",
	Run: func(cmd *cobra.Command, args []string) {
		hostid, err := restClient.HostID()
		if err != nil {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		host, err := restClient.GetHost(hostid)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.UploadFirmware(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostStageFirmware = &cobra.Command{
	Use:   "stage-firmware [hostid]",
	Short: "stage a firmware image for the next boot",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))
		cmd.MarkFlagRequired("package")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		name := viper.GetString("package")
		err = host.StageFirmware(restClient, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostDestageFirmware = &cobra.Command{
	Use:   "destage-firmware [hostid]",
	Short: "destage a firmware package for the next boot",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))
		cmd.MarkFlagRequired("package")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHost(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		name := viper.GetString("package")
		err = host.DestageFirmware(restClient, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(hostCmd)
	hostCmd.AddCommand(hostGetCmd)
	hostCmd.AddCommand(hostInfoCmd)

	hostCmd.AddCommand(hostListCmd)
	addListFlags(hostListCmd)

	hostCmd.AddCommand(hostLogLevelCmd)
	hostLogLevelCmd.Flags().StringP("set", "s", "", "set log level (error/warn/info/debug)")

	hostCmd.AddCommand(hostRestartServicesCmd)
	hostCmd.AddCommand(hostUnjoinCmd)

	hostCmd.AddCommand(hostStateCmd)
	hostStateCmd.Flags().StringP("set", "s", "", "set host state (available/maintenance)")
	addTaskFlags(hostStateCmd)

	hostCmd.AddCommand(hostListSoftwareCmd)
	hostCmd.AddCommand(hostUploadSoftware)
	hostCmd.AddCommand(hostDeploySoftware)
	addTaskFlags(hostDeploySoftware)
	hostDeploySoftware.Flags().String("package", "", "package to deploy")
	hostCmd.AddCommand(hostDeleteSoftware)
	hostDeleteSoftware.Flags().String("package", "", "package to delete")

	hostCmd.AddCommand(hostListFirmwareCmd)
	hostCmd.AddCommand(hostUploadFirmware)
	hostCmd.AddCommand(hostStageFirmware)
	addTaskFlags(hostStageFirmware)
	hostStageFirmware.Flags().String("package", "", "package to stage")
	hostCmd.AddCommand(hostDestageFirmware)
	hostDestageFirmware.Flags().String("package", "", "package to destage")
}
