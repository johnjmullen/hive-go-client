package cmd

import (
	"fmt"
	"io"
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
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		defer file.Close()
		data, _ := io.ReadAll(file)
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
		data, _ := io.ReadAll(file)
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

var poolAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "assign user or group to a stndalone pool",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
		viper.BindPFlag("assign-realm", cmd.Flags().Lookup("assign-realm"))
		viper.BindPFlag("assign-user", cmd.Flags().Lookup("assign-user"))
		viper.BindPFlag("assign-group", cmd.Flags().Lookup("assign-group"))
		cmd.MarkFlagRequired("assign-realm")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var pool *rest.Pool
		var err error
		if !cmd.Flags().Changed("assign-user") && !cmd.Flags().Changed("assign-group") {
			cmd.Usage()
			os.Exit(1)
		}

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
		err = pool.Assign(restClient, viper.GetString("assign-realm"), viper.GetString("assign-user"), viper.GetString("assign-group"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var poolDeleteAssignmentCmd = &cobra.Command{
	Use:   "delete-assignment",
	Short: "delete the assignment for a standalone pool",
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
		err = pool.DeleteAssignment(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var poolSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "snapshot creates disk snapshots for running guests and backs up pool state",
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
		err = pool.Snapshot(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var poolMergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merges snapshots back into the main disk files",
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
		err = pool.Merge(restClient)
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
	poolDeleteCmd.Flags().StringP("id", "i", "", "pool Id")
	poolDeleteCmd.Flags().StringP("name", "n", "", "pool Name")

	poolCmd.AddCommand(poolDiffCmd)

	poolCmd.AddCommand(poolGetCmd)
	poolGetCmd.Flags().StringP("id", "i", "", "pool id")
	poolGetCmd.Flags().StringP("name", "n", "", "pool name")

	poolCmd.AddCommand(poolListCmd)
	addListFlags(poolListCmd)

	poolCmd.AddCommand(poolUpdateCmd)

	poolCmd.AddCommand(poolAssignCmd)
	poolAssignCmd.Flags().StringP("id", "i", "", "pool Id")
	poolAssignCmd.Flags().StringP("name", "n", "", "pool Name")
	poolAssignCmd.Flags().String("assign-realm", "", "realm to assign")
	poolAssignCmd.Flags().String("assign-user", "", "user to assign")
	poolAssignCmd.Flags().String("assign-group", "", "group to assign")

	poolCmd.AddCommand(poolDeleteAssignmentCmd)
	poolDeleteAssignmentCmd.Flags().StringP("id", "i", "", "pool Id")
	poolDeleteAssignmentCmd.Flags().StringP("name", "n", "", "pool Name")

	poolCmd.AddCommand(poolSnapshotCmd)
	poolSnapshotCmd.Flags().StringP("id", "i", "", "pool pool Id")
	poolSnapshotCmd.Flags().StringP("name", "n", "", "pool pool Name")

	poolCmd.AddCommand(poolMergeCmd)
	poolMergeCmd.Flags().StringP("id", "i", "", "pool pool Id")
	poolMergeCmd.Flags().StringP("name", "n", "", "pool pool Name")
}
