package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
	clusterCmd.AddCommand(setLicenseCmd)
}
