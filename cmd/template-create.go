package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var templateCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		var template rest.Template
		err = unmarshal(data, &template)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		msg, err := template.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var templatePoolFile string

func init() {
	templateCmd.AddCommand(templateCreateCmd)
}
