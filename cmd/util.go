package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "hioctl utilities",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
	Hidden: true,
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

 . <(hioctl completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

var docCmd = &cobra.Command{
	Use:   "doc [directory]",
	Args:  cobra.ExactArgs(1),
	Short: "Generates documentation",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(RootCmd, args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(completionCmd)
	utilCmd.AddCommand(docCmd)
}
