package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var poolDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete pool pool",
	Run: func(cmd *cobra.Command, args []string) {
		var pool *rest.Pool
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			id, _ := cmd.Flags().GetString("id")
			pool, err = restClient.GetPool(id)
		case cmd.Flags().Changed("name"):
			name, _ := cmd.Flags().GetString("name")
			pool, err = restClient.GetPoolByName(name)
		default:
			cmd.Usage()
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = pool.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	poolCmd.AddCommand(poolDeleteCmd)
	poolDeleteCmd.Flags().StringP("id", "i", "", "pool pool Id")
	poolDeleteCmd.Flags().StringP("name", "n", "", "pool pool Name")
}
