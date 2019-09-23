package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

func init() {
	clusterCmd.AddCommand(clusterListCmd)
	clusterListCmd.Flags().Bool("details", false, "show details")
	clusterListCmd.Flags().String("filter", "", "filter query string")
}
