package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hive-io/hive-go-client/rest"
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
		cluster.Broker = nil //Hide broker settings because of base64 images
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
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := restClient.ListClusters(listFlagsToQuery())
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
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cluster, err := restClient.GetCluster(clusterID)
		cluster.SetLicense(restClient, args[0])
	},
}

var enableBackupCmd = &cobra.Command{
	Use:   "enable-backup",
	Short: "Enable Backup",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("start-window", cmd.Flags().Lookup("start-window"))
		viper.BindPFlag("end-window", cmd.Flags().Lookup("end-window"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cluster, err := restClient.GetCluster(clusterID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = cluster.EnableBackup(restClient, viper.GetString("start-window"), viper.GetString("end-window"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var disableBackupCmd = &cobra.Command{
	Use:   "disable-backup",
	Short: "Disable Backup",
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cluster, err := restClient.GetCluster(clusterID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = cluster.DisableBackup(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var enableSharedStorageCmd = &cobra.Command{
	Use:   "enable-shared-storage",
	Short: "Enable Shared Storage",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("storage-utilization", cmd.Flags().Lookup("storage-utilization"))
		viper.BindPFlag("set-size", cmd.Flags().Lookup("set-size"))
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cluster, err := restClient.GetCluster(clusterID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		handleTask(cluster.EnableSharedStorage(restClient, viper.GetInt("storage-utilization"), viper.GetInt("set-size")))
	},
}

var disableSharedStorageCmd = &cobra.Command{
	Use:   "disable-shared-storage",
	Short: "Disable Shared Storage",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindTaskFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cluster, err := restClient.GetCluster(clusterID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		handleTask(cluster.DisableSharedStorage(restClient))
	},
}

var clusterGetBrokerCmd = &cobra.Command{
	Use:   "get-broker",
	Short: "get broker settings",
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		broker, err := restClient.GetBroker(clusterID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(broker))
	},
}

var clusterSetBrokerCommand = &cobra.Command{
	Use:   "set-broker [file]",
	Short: "set broker settings",
	Args:  cobra.ExactArgs(1),
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
		data, _ := ioutil.ReadAll(file)
		var broker rest.Broker
		err = unmarshal(data, &broker)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = restClient.SetBroker(clusterID, broker)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var clusterResetBrokerCommand = &cobra.Command{
	Use:   "reset-broker [file]",
	Short: "reset broker settings to the defaults",
	Run: func(cmd *cobra.Command, args []string) {
		clusterID, err := restClient.ClusterID()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = restClient.ResetBroker(clusterID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
	addListFlags(clusterListCmd)

	clusterCmd.AddCommand(setLicenseCmd)
	clusterCmd.AddCommand(enableBackupCmd)
	enableBackupCmd.Flags().String("start-window", "00:00:00", "Time to start running backups")
	enableBackupCmd.Flags().String("end-window", "04:00:00", "Time to stop running backups")
	clusterCmd.AddCommand(disableBackupCmd)

	clusterCmd.AddCommand(enableSharedStorageCmd)
	enableSharedStorageCmd.Flags().IntP("storage-utilization", "s", 75, "Percentage of disk to allocate to shared storage")
	enableSharedStorageCmd.Flags().Int("set-size", 3, "minimum number of hosts to increase the shared storage by")
	addTaskFlags(enableSharedStorageCmd)
	clusterCmd.AddCommand(disableSharedStorageCmd)
	addTaskFlags(disableSharedStorageCmd)

	clusterCmd.AddCommand(clusterGetBrokerCmd)
	clusterCmd.AddCommand(clusterSetBrokerCommand)
	clusterCmd.AddCommand(clusterResetBrokerCommand)
}
