package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var rUser, rPass string
var addHostCmd = &cobra.Command{
	Use:   "add-host [ipAddress]",
	Short: "add a host to the cluster",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("wait") && viper.GetBool("progress-bar") {
			fmt.Printf("Adding %s to cluster\n", args[0])
		}
		handleTask(restClient.JoinHost(rUser, rPass, args[0]))
	},
}

var clusterGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get cluster details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cluster, err := restClient.GetCluster(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(cluster))
	},
}

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "list clusters",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := restClient.ListClusters(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(clusters))
		} else {
			list := []string{}
			for _, cluster := range clusters {
				list = append(list, cluster.ID)
			}
			fmt.Println(formatString(list))
		}
	},
}

var setLicenseCmd = &cobra.Command{
	Use:   "set-license [license]",
	Short: "add a license for the cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterId()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cluster, err := restClient.GetCluster(clusterID)
		cluster.SetLicense(restClient, args[0])
	},
}

func init() {
	RootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(addHostCmd)
	addHostCmd.Flags().StringVar(&rUser, "remote-username", "admin", "username for the remote host")
	addHostCmd.Flags().StringVar(&rPass, "remote-password", "admin", "password for the remote host")
	addTaskFlags(addHostCmd)

	clusterCmd.AddCommand(clusterGetCmd)

	clusterCmd.AddCommand(clusterListCmd)
	clusterListCmd.Flags().Bool("details", false, "show details")
	clusterListCmd.Flags().String("filter", "", "filter query string")

	clusterCmd.AddCommand(setLicenseCmd)
}
