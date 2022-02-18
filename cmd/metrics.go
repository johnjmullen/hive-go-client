package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var metricsCmd = &cobra.Command{
	Use:   "metric",
	Short: "metrics operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var metricsLatestCmd = &cobra.Command{
	Use:   "latest [metric] [entity]",
	Short: "latest metric",
	Args:  cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("resolution", cmd.Flags().Lookup("resolution"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		metrics, err := restClient.MetricsLatest(args[0], args[1], viper.GetUint("resolution"))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(metrics))
	},
}

var metricsExportCmd = &cobra.Command{
	Use:   "export [metric] [entity]",
	Short: "export metric",
	Args:  cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("resolution", cmd.Flags().Lookup("resolution"))
		viper.BindPFlag("output", cmd.Flags().Lookup("output"))
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{"host", "guest", "storage", "pool", "cluster"}, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		metrics, err := restClient.MetricsExport(args[0], args[1], viper.GetUint("resolution"), viper.GetString("output"))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(metrics))
	},
}

/*var metricsLatestHostCmd = &cobra.Command{
	Use:   "host [hostid]",
	Short: "get host metrics",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("resolution", cmd.Flags().Lookup("resolution"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		metrics, err := restClient.MetricsLatest("host", args[0], viper.GetUint("resolution"))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(metrics))
	},
}*/

func init() {
	RootCmd.AddCommand(metricsCmd)
	metricsCmd.AddCommand(metricsLatestCmd)
	metricsCmd.AddCommand(metricsExportCmd)
	metricsLatestCmd.Flags().Uint("resolution", 20, "resolution")
	metricsExportCmd.Flags().Uint("resolution", 20, "resolution")
	metricsExportCmd.Flags().String("output", "json", "output")
}
