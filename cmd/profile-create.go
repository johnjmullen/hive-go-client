package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	rest "github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var profileCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
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

var profilePoolFile string

func init() {
	profileCmd.AddCommand(profileCreateCmd)
}
