package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "profile operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var profileCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new profile",
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
		var profile rest.Profile
		err = unmarshal(data, &profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := profile.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete profile profile",
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			id, _ := cmd.Flags().GetString("id")
			profile, err = restClient.GetProfile(id)
		case cmd.Flags().Changed("name"):
			name, _ := cmd.Flags().GetString("name")
			profile, err = restClient.GetProfileByName(name)
		default:
			cmd.Usage()
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = profile.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get profile details",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			profile, err = restClient.GetProfile(viper.GetString("id"))
		case cmd.Flags().Changed("name"):
			profile, err = restClient.GetProfileByName(viper.GetString("name"))
		default:
			cmd.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(profile))
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := restClient.ListProfiles(listFlagsToQuery())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(profiles))
		} else {
			list := []map[string]string{}
			for _, profile := range profiles {
				var info = map[string]string{"id": profile.ID, "name": profile.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "update a profile",
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
		var profile rest.Profile
		err = unmarshal(data, &profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := profile.Update(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileCreateCmd)

	profileCmd.AddCommand(profileDeleteCmd)
	profileDeleteCmd.Flags().StringP("id", "i", "", "profile profile Id")
	profileDeleteCmd.Flags().StringP("name", "n", "", "profile profile Name")

	profileCmd.AddCommand(profileGetCmd)
	profileGetCmd.Flags().StringP("id", "i", "", "profile id")
	profileGetCmd.Flags().StringP("name", "n", "", "profile name")

	profileCmd.AddCommand(profileListCmd)
	addListFlags(profileListCmd)

	profileCmd.AddCommand(profileUpdateCmd)
}
