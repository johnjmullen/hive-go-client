package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	rest "bitbucket.org/johnmullen/hiveio-go-client/rest"
	"github.com/ghodss/yaml"
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
	RootCmd.PersistentFlags().String("host", "", "Server to connect to")
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
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/hiveio/")
		viper.AddConfigPath("$HOME/.hiveio")
	}
	viper.SetEnvPrefix("hio")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func connectRest(cmd *cobra.Command, args []string) {
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

func jsonString(obj interface{}) string {
	json, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(json)
}

func formatString(obj interface{}) string {
	switch viper.GetString("format") {
	case "yaml":
		return yamlString(obj)
	case "json":
		return jsonString(obj)
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
