package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
	"github.com/spf13/cobra"
)

var storageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a new storage pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		var sp *rest.StoragePool
		err = unmarshal(data, sp)
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

var storagePoolFile string

func init() {
	storageCmd.AddCommand(storageCreateCmd)
}
