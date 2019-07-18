package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var poolGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get pool details",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var pool *rest.Pool
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			pool, err = restClient.GetPool(viper.GetString("id"))
		case cmd.Flags().Changed("name"):
			pool, err = restClient.GetPoolByName(viper.GetString("name"))
		default:
			cmd.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(pool))
	},
}

func init() {
	poolCmd.AddCommand(poolGetCmd)
	poolGetCmd.Flags().StringP("id", "i", "", "pool id")
	poolGetCmd.Flags().StringP("name", "n", "", "pool name")
}
