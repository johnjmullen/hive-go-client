package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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

//var name string

func init() {
	clusterCmd.AddCommand(clusterGetCmd)
}
