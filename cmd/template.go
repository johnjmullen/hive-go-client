package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "template operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var templateAnalyzeCmd = &cobra.Command{
	Use:   "analyze [Name]",
	Short: "analyze template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Analyze(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var templateAuthorCmd = &cobra.Command{
	Use:   "author [Name]",
	Short: "author template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Author(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var templateCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new template",
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

		msg, err := template.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete template pool",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var templateDiffCmd = &cobra.Command{
	Use:   "diff [template1] [template2]",
	Short: "compare 2 templates",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		template1, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		template2, err := restClient.GetTemplate(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(cmp.Diff(template1, template2))
	},
}

var templateGetCmd = &cobra.Command{
	Use:   "get [Name]",
	Short: "get template details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(template))
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "list templates",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := restClient.ListTemplates(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(templates))
		} else {
			list := []string{}
			for _, template := range templates {
				list = append(list, template.Name)
			}
			fmt.Println(formatString(list))
		}
	},
}

var templateLoadCmd = &cobra.Command{
	Use:   "load [Name]",
	Short: "load template to all hosts",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("storage", cmd.Flags().Lookup("storage"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Load(restClient, viper.GetString("storage"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var templateUnloadCmd = &cobra.Command{
	Use:   "unload [Name]",
	Short: "unload template from all hosts",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, err := restClient.GetTemplate(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = template.Unload(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

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
	RootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateAnalyzeCmd)
	templateCmd.AddCommand(templateAuthorCmd)
	templateCmd.AddCommand(templateCreateCmd)

	templateCmd.AddCommand(templateDeleteCmd)
	templateDeleteCmd.Flags().StringP("id", "i", "", "template Pool Id")
	templateDeleteCmd.Flags().StringP("name", "n", "", "template Pool Name")

	templateCmd.AddCommand(templateDiffCmd)
	templateCmd.AddCommand(templateGetCmd)

	templateCmd.AddCommand(templateListCmd)
	templateListCmd.Flags().Bool("details", false, "show details")
	templateListCmd.Flags().String("filter", "", "filter query string")

	templateCmd.AddCommand(templateLoadCmd)
	templateLoadCmd.Flags().StringP("storage", "s", "disk", "Location to load the template (disk or ram)")

	templateCmd.AddCommand(templateUnloadCmd)
	templateCmd.AddCommand(templateUpdateCmd)
}
