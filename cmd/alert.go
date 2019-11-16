package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "alert operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(alertCmd)
	alertCmd.AddCommand(alertListCmd)
	alertListCmd.Flags().Bool("details", false, "show details")
	alertListCmd.Flags().String("filter", "", "filter query string")
}

// listCmd represents the list command
var alertListCmd = &cobra.Command{
	Use:   "list",
	Short: "list alerts",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		alerts, err := restClient.ListAlerts(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(alerts))
		} else {
			list := []map[string]string{}
			for _, alert := range alerts {
				var info = map[string]string{"id": alert.ID, "message": alert.Message}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}
