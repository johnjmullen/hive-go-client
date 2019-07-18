package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var realmCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new realm",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		var realm rest.Realm
		err = unmarshal(data, &realm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := realm.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var realmPoolFile string

func init() {
	realmCmd.AddCommand(realmCreateCmd)
}
