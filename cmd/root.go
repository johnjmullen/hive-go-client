package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var restClient *rest.Client

var (
	version = "dev"
	commit  = ""
	date    = ""
)

//RootCmd root command for hioctl
var RootCmd = &cobra.Command{
	Use:              "hioctl",
	Short:            "hive fabric rest api client",
	Version:          fmt.Sprintf("%s ", version),
	PersistentPreRun: connectRest,
	TraverseChildren: true,
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "hioctl version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", version)
		if commit != "" {
			fmt.Printf("commit: %s\n", commit)
		}
		if date != "" {
			fmt.Printf("commit date: %s\n", date)
		}
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			fmt.Printf("module version: %s, checksum: %s\n", info.Main.Version, info.Main.Sum)
		}
		fmt.Printf("https://github.com/hive-io/hive-go-client\n")
	},
}

//Execute run root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "config file")
	RootCmd.PersistentFlags().BoolP("insecure", "k", false, "ignore certificate errors")
	RootCmd.PersistentFlags().String("host", "", "Hostname or ip address")
	RootCmd.PersistentFlags().Uint("port", 8443, "port")
	RootCmd.PersistentFlags().StringP("user", "u", "admin", "Admin username")
	RootCmd.PersistentFlags().StringP("password", "p", "", "Admin user password")
	RootCmd.PersistentFlags().StringP("realm", "r", "local", "Admin user realm")
	RootCmd.PersistentFlags().StringP("format", "", "json", "format (json/yaml)")

	viper.BindPFlag("host", RootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("insecure", RootCmd.PersistentFlags().Lookup("insecure"))
	viper.BindPFlag("user", RootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("realm", RootCmd.PersistentFlags().Lookup("realm"))
	viper.BindPFlag("format", RootCmd.PersistentFlags().Lookup("format"))

	RootCmd.AddCommand(VersionCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("hioctl")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.hiveio")
		viper.AddConfigPath("/etc/hive/")
		viper.AddConfigPath("/etc/hiveio/")
	}
	viper.SetEnvPrefix("hio")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func connectRest(cmd *cobra.Command, args []string) {
	cmd.MarkPersistentFlagRequired("host")
	if viper.GetString("host") == "" {
		fmt.Println("Error: Host was not provided.")
		cmd.Usage()
		os.Exit(1)
	}
	restClient = &rest.Client{
		Host:          viper.GetString("host"),
		Port:          viper.GetUint("port"),
		AllowInsecure: viper.GetBool("insecure"),
		UserAgent:     "hioctl/" + version,
	}
	err := restClient.Login(viper.GetString("user"), viper.GetString("password"), viper.GetString("realm"))
	if err != nil {
		fmt.Printf("Error: Failed to connect %v", err)
		os.Exit(1)
	}
}

func yamlString(obj interface{}) string {
	yaml, err := yaml.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(yaml)
}

func jsonString(obj interface{}, compact bool) string {
	var data []byte
	var err error
	if compact {
		data, err = json.Marshal(obj)
	} else {
		data, err = json.MarshalIndent(obj, "", "  ")
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(data)
}

func formatString(obj interface{}) string {
	switch viper.GetString("format") {
	case "yaml":
		return yamlString(obj)
	case "json":
		return jsonString(obj, false)
	case "json-compact":
		return jsonString(obj, true)
	default:
		fmt.Println("Error: Unsupported format")
		os.Exit(1)
	}
	return ""
}

func unmarshal(data []byte, obj interface{}) error {
	var err error
	format := viper.GetString("format")
	switch format {
	case "yaml":
		return yaml.Unmarshal(data, obj)
	case "json":
		return json.Unmarshal(data, obj)
	default:
		err = (fmt.Errorf("Error: Unsupported format %s", format))
	}
	return err
}

func addListFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("details", false, "show details")
	cmd.Flags().String("filter", "", "filter results based on a field.")
	cmd.Flags().Int("count", 1000, "number of results to show")
	cmd.Flags().Int("offset", 0, "first result to show")
}

func bindListFlags(cmd *cobra.Command) {
	viper.BindPFlag("details", cmd.Flags().Lookup("details"))
	viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	viper.BindPFlag("count", cmd.Flags().Lookup("count"))
	viper.BindPFlag("offset", cmd.Flags().Lookup("offset"))
}

func listFlagsToQuery() string {
	Values, err := url.ParseQuery(viper.GetString("filter"))
	if err != nil {
		fmt.Println("Error: Unable to parse filter")
		os.Exit(1)
	}
	Values.Add("count", strconv.Itoa(viper.GetInt("count")))
	Values.Add("offset", strconv.Itoa(viper.GetInt("offset")))
	return Values.Encode()
}
