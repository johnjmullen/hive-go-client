package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "pool operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var poolCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new guest pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		var err error
		if args[0] == "-" {
			fmt.Println("reading stdin")
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		var pool rest.Pool
		err = unmarshal(data, &pool)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := pool.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var poolDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a pool",
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

var poolListCmd = &cobra.Command{
	Use:   "list",
	Short: "list pools",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		pools, err := restClient.ListGuestPools(listFlagsToQuery())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(pools))
		} else {
			list := []map[string]string{}
			for _, pool := range pools {
				var info = map[string]string{"id": pool.ID, "name": pool.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

var poolUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "update a guest pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		var err error
		if args[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		var pool rest.Pool
		err = unmarshal(data, &pool)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := pool.Update(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(poolCmd)
	poolCmd.AddCommand(poolCreateCmd)

	poolCmd.AddCommand(poolDeleteCmd)
	poolDeleteCmd.Flags().StringP("id", "i", "", "pool pool Id")
	poolDeleteCmd.Flags().StringP("name", "n", "", "pool pool Name")

	poolCmd.AddCommand(poolDiffCmd)

	poolCmd.AddCommand(poolGetCmd)
	poolGetCmd.Flags().StringP("id", "i", "", "pool id")
	poolGetCmd.Flags().StringP("name", "n", "", "pool name")

	poolCmd.AddCommand(poolListCmd)
	addListFlags(poolListCmd)

	poolCmd.AddCommand(poolUpdateCmd)
}
