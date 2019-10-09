package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var templateUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "update a template",
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
		var template rest.Template
		err = unmarshal(data, &template)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := template.Update(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	templateCmd.AddCommand(templateUpdateCmd)
}
