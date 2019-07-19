package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

func init() {
	realmCmd.AddCommand(realmCreateCmd)
	realmCreateCmd.Flags().StringP("name", "n", "", "Netbios Name")
	realmCreateCmd.Flags().String("fqdn", "", "FQDN")
}
