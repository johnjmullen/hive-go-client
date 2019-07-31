package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export data",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var data struct {
			Hosts        []rest.Host
			Realms       []rest.Realm
			Profiles     []rest.Profile
			StoragePools []rest.StoragePool
			Templates    []rest.Template
			Pools        []rest.Pool
		}

		data.Realms, err = restClient.ListRealms("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data.Profiles, err = restClient.ListProfiles("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data.StoragePools, err = restClient.ListStoragePools("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data.Templates, err = restClient.ListTemplates("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data.Pools, err = restClient.ListGuestPools("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data.Hosts, err = restClient.ListHosts("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(formatString(data))
	},
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
