package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "list clusters",
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := restClient.ListClusters()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(clusters))
		} else {
			for _, cluster := range clusters {
				fmt.Println(cluster.ID)
			}
		}
	},
}

func init() {
	clusterCmd.AddCommand(clusterListCmd)
	clusterListCmd.Flags().Bool("details", false, "show details")
}
