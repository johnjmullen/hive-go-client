package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var restClient *rest.Client
var RootCmd = &cobra.Command{
	Use:              "hioctl",
	Short:            "hiveio cli",
	PersistentPreRun: connectRest,
	TraverseChildren: true,
}

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
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("hioctl")
		viper.AddConfigPath("/etc/hive/")
		viper.AddConfigPath("/etc/hiveio/")
		viper.AddConfigPath("$HOME/.hiveio")
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
	restClient = &rest.Client{Host: viper.GetString("host"), Port: viper.GetUint("port"), AllowInsecure: viper.GetBool("insecure")}
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
