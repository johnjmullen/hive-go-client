package cmd

import (
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

var poolDiffCmd = &cobra.Command{
	Use:   "diff [pool1 id] [pool2 id]",
	Short: "compare 2 pools",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pool1, err := restClient.GetPool(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pool2, err := restClient.GetPool(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(cmp.Diff(pool1, pool2))
	},
}

func init() {
	poolCmd.AddCommand(poolDiffCmd)
}
