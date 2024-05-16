package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var hostNetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "host network operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var hostListNetworksCmd = &cobra.Command{
	Use:    "list {-i hostid | -n hostname | --ip ip_address }",
	Short:  "list networks on a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		networks, err := host.ListNetworks(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(networks))

	},
}

var hostListInterfacesCmd = &cobra.Command{
	Use:    "interfaces {-i hostid | -n hostname | --ip ip_address }",
	Short:  "list network interfaces on a host",
	PreRun: bindHostIDFlags,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		interfaces, err := host.ListInterfaces(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(interfaces))

	},
}

var hostGetNetworkCmd = &cobra.Command{
	Use:    "get {-i hostid | -n hostname | --ip ip_address } network ",
	Short:  "get network settings",
	PreRun: bindHostIDFlags,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		network, err := host.GetNetwork(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(network))

	},
}

var hostSetNetworkCmd = &cobra.Command{
	Use:    "set {-i hostid | -n hostname | --ip ip_address } network file",
	Short:  "create or edit a network",
	PreRun: bindHostIDFlags,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		var err error
		if args[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		defer file.Close()
		data, _ := io.ReadAll(file)
		network := rest.HostNetwork{
			Name: args[0],
		}
		err = unmarshal(data, &network)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.SetNetwork(restClient, network)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

var hostDeleteNetworkCmd = &cobra.Command{
	Use:    "delete {-i hostid | -n hostname | --ip ip_address } network ",
	Short:  "delete a host network",
	PreRun: bindHostIDFlags,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := getHost(cmd, args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = host.DeleteNetwork(restClient, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	hostCmd.AddCommand(hostNetworkCmd)
	hostNetworkCmd.AddCommand(hostListNetworksCmd)
	addHostIDFlags(hostListNetworksCmd)
	hostNetworkCmd.AddCommand(hostGetNetworkCmd)
	addHostIDFlags(hostGetNetworkCmd)
	hostNetworkCmd.AddCommand(hostSetNetworkCmd)
	addHostIDFlags(hostSetNetworkCmd)
	hostNetworkCmd.AddCommand(hostDeleteNetworkCmd)
	addHostIDFlags(hostDeleteNetworkCmd)
	hostNetworkCmd.AddCommand(hostListInterfacesCmd)
	addHostIDFlags(hostListInterfacesCmd)
}
