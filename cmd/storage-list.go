package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var storageListCmd = &cobra.Command{
	Use:   "list",
	Short: "list storage pools",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		pools, err := restClient.ListStoragePools(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(pools))
		} else {
			var list []map[string]string
			for _, pool := range pools {
				var info = map[string]string{"id": pool.ID, "name": pool.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

func init() {
	storageCmd.AddCommand(storageListCmd)
	storageListCmd.Flags().Bool("details", false, "show details")
	storageListCmd.Flags().String("filter", "", "filter query string")
}
