package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hive-io/hive-go-client/rest"
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

func addHostIDFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("id", "i", "", "hostid")
	cmd.Flags().StringP("name", "n", "", "hostname")
	cmd.Flags().String("ip", "", "host ip address")
}

func bindHostIDFlags(cmd *cobra.Command, args []string) {
	viper.BindPFlag("id", cmd.Flags().Lookup("id"))
	viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	viper.BindPFlag("ip", cmd.Flags().Lookup("ip"))
}

func getHost(cmd *cobra.Command, args []string) (*rest.Host, error) {
	switch {
	case cmd.Flags().Changed("id"):
		host, err := restClient.GetHost(viper.GetString("id"))
		return &host, err
	case cmd.Flags().Changed("name"):
		return restClient.GetHostByName(viper.GetString("name"))
	case cmd.Flags().Changed("ip"):
		return restClient.GetHostByIP(viper.GetString("ip"))
	case len(args) == 1 && strings.Contains(cmd.Use, "| hostid}"):
		host, err := restClient.GetHost(args[0])
		return &host, err
	default:
		cmd.Usage()
		os.Exit(1)
	}
	return nil, fmt.Errorf("error getting host")
}

var hostGetCmd = &cobra.Command{
	Use:    "get {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "get host details",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
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

var hostGetIDCmd = &cobra.Command{
	Use:   "get-id name",
	Short: "get hostid from hostname",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := restClient.GetHostByName(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(host.Hostid)
	},
}

var hostLogLevelCmd = &cobra.Command{
	Use:   "log-level {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short: "get or set host log level",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("set", cmd.Flags().Lookup("set"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
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
	Use:    "restart-services {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "restart hive servies",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
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

var hostRebootCmd = &cobra.Command{
	Use:    "reboot {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "reboot a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.Reboot(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostShutdownCmd = &cobra.Command{
	Use:    "shutdown {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "shutdown a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.Shutdown(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostUnjoinCmd = &cobra.Command{
	Use:   "unjoin {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short: "remove host from cluster",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Printf("Removing %s from cluster for\n", host.Hostname)
		}
		handleTask(host.UnjoinCluster(restClient))
	},
}

var hostStateCmd = &cobra.Command{
	Use:   "state {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short: "get or set host state",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("set", cmd.Flags().Lookup("set"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
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

var hostEnableGatewayCmd = &cobra.Command{
	Use:    "enable-gateway-mode {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "Convert the host into a gateway appliance",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if host.Appliance.Role == "gateway" {
			return
		}
		err = host.ChangeGatewayMode(restClient, true)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostDisableGatewayCmd = &cobra.Command{
	Use:    "disable-gateway-mode {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "Convert the host from a gateway appliance to a regular fabric host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if host.Appliance.Role != "gateway" {
			return
		}
		err = host.ChangeGatewayMode(restClient, false)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostListSoftwareCmd = &cobra.Command{
	Use:    "list-software {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "list available software packages on a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
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

var hostDeleteSoftware = &cobra.Command{
	Use:   "delete-software {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short: "delete a software package",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("package", cmd.Flags().Lookup("package"))
		cmd.MarkFlagRequired("package")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
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

var hostEnableCRSCmd = &cobra.Command{
	Use:    "enable-crs {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "enable crs on a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.EnableCRS(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostDisableCRSCmd = &cobra.Command{
	Use:    "disable-crs {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "disable crs on a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.DisableCRS(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostDeleteCmd = &cobra.Command{
	Use:    "delete {-i hostid | -n hostname | --ip ip_address | hostid}",
	Short:  "delete a host record from the host table",
	Hidden: true,
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var hostIscsiDiscoverCmd = &cobra.Command{
	Use:   "iscsi-discover",
	Short: "discover iscsi targets",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("portal", cmd.Flags().Lookup("portal"))
		cmd.MarkFlagRequired("portal")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result, err := host.IscsiDiscover(restClient, viper.GetString("portal"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(result))
	},
}

var hostIscsiLoginCmd = &cobra.Command{
	Use:   "iscsi-login",
	Short: "login to an iscsi target",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("portal", cmd.Flags().Lookup("portal"))
		viper.BindPFlag("target", cmd.Flags().Lookup("target"))
		viper.BindPFlag("iscsiUsername", cmd.Flags().Lookup("iscsi-username"))
		viper.BindPFlag("iscsiPassword", cmd.Flags().Lookup("iscsi-password"))
		cmd.MarkFlagRequired("portal")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		authMethod := "None"
		if viper.GetString("iscsiUsername") != "" {
			authMethod = "CHAP"
		}
		result, err := host.IscsiLogin(restClient, viper.GetString("portal"), viper.GetString("target"), authMethod, viper.GetString("iscsiUsername"), viper.GetString("iscsiPassword"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(result))
	},
}

var hostIscsiSessionsCmd = &cobra.Command{
	Use:   "iscsi-sessions",
	Short: "list iscsi sessions",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("portal", cmd.Flags().Lookup("portal"))
		viper.BindPFlag("target", cmd.Flags().Lookup("target"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result, err := host.IscsiSessions(restClient, viper.GetString("portal"), viper.GetString("target"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(result))
	},
}

var hostIscsiLogoutCmd = &cobra.Command{
	Use:   "iscsi-logout",
	Short: "logout from an iscsi target",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindHostIDFlags(cmd, args)
		viper.BindPFlag("portal", cmd.Flags().Lookup("portal"))
		viper.BindPFlag("target", cmd.Flags().Lookup("target"))
		cmd.MarkFlagRequired("portal")
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.IscsiLogout(restClient, viper.GetString("portal"), viper.GetString("target"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(hostCmd)
	hostCmd.AddCommand(hostInfoCmd)
	hostCmd.AddCommand(hostGetIDCmd)
	hostCmd.AddCommand(hostGetCmd)
	addHostIDFlags(hostGetCmd)

	hostCmd.AddCommand(hostListCmd)
	addListFlags(hostListCmd)

	hostCmd.AddCommand(hostLogLevelCmd)
	hostLogLevelCmd.Flags().StringP("set", "s", "", "set log level (error/warn/info/debug)")
	addHostIDFlags(hostLogLevelCmd)

	hostCmd.AddCommand(hostRestartServicesCmd)
	addHostIDFlags(hostRestartServicesCmd)
	hostCmd.AddCommand(hostRebootCmd)
	addHostIDFlags(hostRebootCmd)
	hostCmd.AddCommand(hostShutdownCmd)
	addHostIDFlags(hostShutdownCmd)
	hostCmd.AddCommand(hostUnjoinCmd)
	addHostIDFlags(hostUnjoinCmd)
	addTaskFlags(hostUnjoinCmd)

	hostCmd.AddCommand(hostStateCmd)
	addHostIDFlags(hostStateCmd)
	hostStateCmd.Flags().StringP("set", "s", "", "set host state (available/maintenance)")
	addTaskFlags(hostStateCmd)

	hostCmd.AddCommand(hostListSoftwareCmd)
	addHostIDFlags(hostListSoftwareCmd)
	hostCmd.AddCommand(hostUploadSoftware)
	hostCmd.AddCommand(hostDeleteSoftware)
	addHostIDFlags(hostDeleteSoftware)
	hostDeleteSoftware.Flags().String("package", "", "package to delete")
	hostCmd.AddCommand(hostEnableCRSCmd)
	addHostIDFlags(hostEnableCRSCmd)
	hostCmd.AddCommand(hostDisableCRSCmd)
	addHostIDFlags(hostDisableCRSCmd)
	hostCmd.AddCommand(hostEnableGatewayCmd)
	addHostIDFlags(hostEnableGatewayCmd)
	hostCmd.AddCommand(hostDisableGatewayCmd)
	addHostIDFlags(hostDisableGatewayCmd)
	hostCmd.AddCommand(hostDeleteCmd)
	addHostIDFlags(hostDeleteCmd)

	hostCmd.AddCommand(hostIscsiDiscoverCmd)
	addHostIDFlags(hostIscsiDiscoverCmd)
	hostIscsiDiscoverCmd.Flags().String("portal", "", "portal")
	hostCmd.AddCommand(hostIscsiLoginCmd)
	addHostIDFlags(hostIscsiLoginCmd)
	hostIscsiLoginCmd.Flags().String("portal", "", "portal")
	hostIscsiLoginCmd.Flags().String("target", "", "target")
	hostIscsiLoginCmd.Flags().String("iscsi-username", "", "iscsi username")
	hostIscsiLoginCmd.Flags().String("iscsi-password", "", "iscsi password")
	hostCmd.AddCommand(hostIscsiSessionsCmd)
	addHostIDFlags(hostIscsiSessionsCmd)
	hostIscsiSessionsCmd.Flags().String("portal", "", "filter sessions by portal")
	hostIscsiSessionsCmd.Flags().String("target", "", "filter sessions by target")
	hostCmd.AddCommand(hostIscsiLogoutCmd)
	addHostIDFlags(hostIscsiLogoutCmd)
	hostIscsiLogoutCmd.Flags().String("portal", "", "portal")
	hostIscsiLogoutCmd.Flags().String("target", "", "target")
}
