package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var alertListCmd = &cobra.Command{
	Use:   "list",
	Short: "list alerts",
	Run: func(cmd *cobra.Command, args []string) {
		alerts, err := restClient.ListAlerts()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(alerts))
		} else {
			for _, alert := range alerts {
				fmt.Printf("Acknowledged: %v, Message: %s", alert.Acknowledged, alert.Message)
			}
		}
	},
}

func init() {
	alertCmd.AddCommand(alertListCmd)
	alertListCmd.Flags().Bool("details", false, "show details")
}
