package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var storageCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new storage pool",
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
		var sp rest.StoragePool
		err = unmarshal(data, &sp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := sp.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	storageCmd.AddCommand(storageCreateCmd)
}
