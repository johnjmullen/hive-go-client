package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var userCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new user",
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
		var user rest.User
		err = unmarshal(data, &user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := user.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user, err := restClient.GetUser(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = user.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var userDiffCmd = &cobra.Command{
	Use:   "diff [user1 id] [user2 id]",
	Short: "compare 2 users",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		user1, err := restClient.GetUser(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		user2, err := restClient.GetUser(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(cmp.Diff(user1, user2))
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get user details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user, err := restClient.GetUser(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(user))
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "list users",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		users, err := restClient.ListUsers(listFlagsToQuery())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(users))
		} else {
			list := []map[string]string{}
			for _, user := range users {
				var info = map[string]string{"id": user.ID}
				if user.Username != "" {
					info["username"] = user.Username
				}
				if user.GroupName != "" {
					info["groupname"] = user.GroupName
				}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

var userUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "update a user",
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
		var user rest.User
		err = unmarshal(data, &user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := user.Update(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)

	userCmd.AddCommand(userDeleteCmd)

	userCmd.AddCommand(userDiffCmd)

	userCmd.AddCommand(userGetCmd)

	userCmd.AddCommand(userListCmd)
	addListFlags(userListCmd)

	userCmd.AddCommand(userUpdateCmd)
}
