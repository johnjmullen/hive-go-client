package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var realmCmd = &cobra.Command{
	Use:   "realm",
	Short: "realm operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var realmCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new realm",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.MarkFlagRequired("name")
			cmd.MarkFlagRequired("fqdn")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.New()
		if len(args) == 1 {
			v.SetConfigFile(args[0])
			v.ReadInConfig()
		}
		v.BindPFlag("fqdn", cmd.Flags().Lookup("fqdn"))
		v.BindPFlag("name", cmd.Flags().Lookup("name"))
		var realm rest.Realm
		err := v.Unmarshal(&realm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		realm.Verified = true
		msg, err := realm.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var realmDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete realm pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		realm, err := restClient.GetRealm(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = realm.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var realmGetCmd = &cobra.Command{
	Use:   "get [Name]",
	Short: "get realm details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm, err := restClient.GetRealm(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(realm))
	},
}

var realmListCmd = &cobra.Command{
	Use:   "list",
	Short: "list realms",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindListFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		realms, err := restClient.ListRealms(listFlagsToQuery())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(realms))
		} else {
			list := []string{}
			for _, realm := range realms {
				list = append(list, realm.Name)
			}
			fmt.Println(formatString(list))
		}
	},
}

func init() {
	RootCmd.AddCommand(realmCmd)

	realmCmd.AddCommand(realmCreateCmd)
	realmCreateCmd.Flags().StringP("name", "n", "", "Netbios Name")
	realmCreateCmd.Flags().String("fqdn", "", "FQDN")

	realmCmd.AddCommand(realmDeleteCmd)
	realmCmd.AddCommand(realmGetCmd)

	realmCmd.AddCommand(realmListCmd)
	addListFlags(realmListCmd)
}
