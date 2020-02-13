package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "profile operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(0)
	},
}

var profileCreateCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Add a new profile",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.MarkFlagRequired("name")
		}
		if cmd.Flags().Changed("enable-ad") {
			cmd.MarkFlagRequired("ad-username")
			cmd.MarkFlagRequired("ad-password")
			cmd.MarkFlagRequired("ad-domain")
			cmd.MarkFlagRequired("ad-user-group")
		}
		if cmd.Flags().Changed("enable-uv") {
			cmd.MarkFlagRequired("uv-repository")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.New()
		if len(args) == 1 {
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
			v.SetConfigType(viper.GetString("format"))
			v.SetConfigFile(args[0])
			v.ReadConfig(file)
		}
		v.BindPFlag("id", cmd.Flags().Lookup("id"))
		v.BindPFlag("name", cmd.Flags().Lookup("name"))
		v.BindPFlag("timezone", cmd.Flags().Lookup("timezone"))
		if cmd.Flags().Changed("enable-ad") {
			v.BindPFlag("adConfig.domain", cmd.Flags().Lookup("ad-domain"))
			v.BindPFlag("adConfig.username", cmd.Flags().Lookup("ad-username"))
			v.BindPFlag("adConfig.password", cmd.Flags().Lookup("ad-password"))
			v.BindPFlag("adConfig.userGroup", cmd.Flags().Lookup("ad-user-group"))
		}

		if cmd.Flags().Changed("enable-broker-options") {
			v.BindPFlag("brokerOptions.allowDesktopComposition", cmd.Flags().Lookup("broker-allow-desktop-composition"))
			v.BindPFlag("brokerOptions.audioCapture", cmd.Flags().Lookup("broker-audio-capture"))
			v.BindPFlag("brokerOptions.redirectCSSP", cmd.Flags().Lookup("broker-redirect-cssp"))
			v.BindPFlag("brokerOptions.redirectClipboard", cmd.Flags().Lookup("broker-redirect-clipboard"))
			v.BindPFlag("brokerOptions.redirectPNP", cmd.Flags().Lookup("broker-redirect-pnp"))
			v.BindPFlag("brokerOptions.redirectUSB", cmd.Flags().Lookup("broker-redirect-usb"))
			v.BindPFlag("brokerOptions.redirectPrinter", cmd.Flags().Lookup("broker-redirect-printer"))
			v.BindPFlag("brokerOptions.redirectSmartcard", cmd.Flags().Lookup("broker-redirect-smartcard"))
			v.BindPFlag("brokerOptions.smartResize", cmd.Flags().Lookup("broker-smart-resize"))
			v.BindPFlag("brokerOptions.hideCertificateWarnings", cmd.Flags().Lookup("broker-hide-certificate-warnings"))
		}
		if cmd.Flags().Changed("enable-uv") {
			v.BindPFlag("userVolumes.backupSchedule", cmd.Flags().Lookup("uv-backup-schedule"))
			v.BindPFlag("userVolumes.repository", cmd.Flags().Lookup("uv-repository"))
			v.BindPFlag("userVolumes.size", cmd.Flags().Lookup("uv-size"))
			v.BindPFlag("userVolumes.target", cmd.Flags().Lookup("uv-target"))
		}
		var profile rest.Profile
		err := v.Unmarshal(&profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		msg, err := profile.Create(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(msg))
	},
}

func addProfileFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("id", "i", "", "profile id")
	cmd.Flags().StringP("name", "n", "", "profile name")
	cmd.Flags().String("timezone", "disabled", "timezone to inject")

	cmd.Flags().Bool("enable-ad", false, "enable active directory options")
	cmd.Flags().String("ad-domain", "", "AD realm")
	cmd.Flags().String("ad-username", "", "user to join guests to AD")
	cmd.Flags().String("ad-password", "", "Password for the join user")
	cmd.Flags().String("ad-user-group", "", "User group")

	cmd.Flags().Bool("enable-broker-options", false, "enable broker options")
	cmd.Flags().Bool("broker-allow-desktop-composition", true, "Allow Destop compositing")
	cmd.Flags().Bool("broker-audio-capture", true, "audio capture")
	cmd.Flags().Bool("broker-redirect-cssp", true, "CredSSP redirection")
	cmd.Flags().Bool("broker-redirect-clipboard", true, "Clipboard redirection")
	cmd.Flags().Bool("broker-redirect-pnp", true, "Plug-and-Play redirection")
	cmd.Flags().Bool("broker-redirect-usb", true, "USB redirection")
	cmd.Flags().Bool("broker-redirect-printer", true, "Printer redirection")
	cmd.Flags().Bool("broker-redirect-smartcard", true, "Smartcard redirection")
	cmd.Flags().Bool("broker-smart-resize", true, "Smart screen resize")
	cmd.Flags().Bool("broker-hide-certificate-warnings", false, "Hide certificate warnings")

	cmd.Flags().Bool("enable-uv", false, "Enable user volumes")
	cmd.Flags().Int("uv-backup-schedule", 28800, "User volume backup schedule (s)")
	cmd.Flags().String("uv-repository", "", "Storage pool for storing the user volume")
	cmd.Flags().Int("uv-size", 10, "User volume size (GB)")
	cmd.Flags().String("uv-target", "disk", "Local cache (disk/ram)")
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete profile profile",
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			id, _ := cmd.Flags().GetString("id")
			profile, err = restClient.GetProfile(id)
		case cmd.Flags().Changed("name"):
			name, _ := cmd.Flags().GetString("name")
			profile, err = restClient.GetProfileByName(name)
		default:
			cmd.Usage()
			os.Exit(1)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = profile.Delete(restClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get profile details",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var profile *rest.Profile
		var err error
		switch {
		case cmd.Flags().Changed("id"):
			profile, err = restClient.GetProfile(viper.GetString("id"))
		case cmd.Flags().Changed("name"):
			profile, err = restClient.GetProfileByName(viper.GetString("name"))
		default:
			cmd.Usage()
			os.Exit(1)
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(formatString(profile))
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("filter", cmd.Flags().Lookup("filter"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := restClient.ListProfiles(viper.GetString("filter"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("details") {
			fmt.Println(formatString(profiles))
		} else {
			list := []map[string]string{}
			for _, profile := range profiles {
				var info = map[string]string{"id": profile.ID, "name": profile.Name}
				list = append(list, info)
			}
			fmt.Println(formatString(list))
		}
	},
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update [file]",
	Short: "update a profile",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.MarkFlagRequired("name")
		}
		if cmd.Flags().Changed("enable-ad") {
			cmd.MarkFlagRequired("ad-username")
			cmd.MarkFlagRequired("ad-password")
			cmd.MarkFlagRequired("ad-domain")
			cmd.MarkFlagRequired("ad-user-group")
		}
		if cmd.Flags().Changed("enable-uv") {
			cmd.MarkFlagRequired("uv-repository")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.New()
		if len(args) == 1 {
			v.SetConfigFile(args[0])
			v.ReadInConfig()
		}
		v.BindPFlag("id", cmd.Flags().Lookup("id"))
		v.BindPFlag("name", cmd.Flags().Lookup("name"))
		v.BindPFlag("timezone", cmd.Flags().Lookup("timezone"))
		if cmd.Flags().Changed("enable-ad") {
			v.BindPFlag("adConfig.domain", cmd.Flags().Lookup("ad-domain"))
			v.BindPFlag("adConfig.username", cmd.Flags().Lookup("ad-username"))
			v.BindPFlag("adConfig.password", cmd.Flags().Lookup("ad-password"))
			v.BindPFlag("adConfig.userGroup", cmd.Flags().Lookup("ad-user-group"))
		}

		if cmd.Flags().Changed("enable-broker-options") {
			v.BindPFlag("brokerOptions.allowDesktopComposition", cmd.Flags().Lookup("broker-allow-desktop-composition"))
			v.BindPFlag("brokerOptions.audioCapture", cmd.Flags().Lookup("broker-audio-capture"))
			v.BindPFlag("brokerOptions.redirectCSSP", cmd.Flags().Lookup("broker-redirect-cssp"))
			v.BindPFlag("brokerOptions.redirectClipboard", cmd.Flags().Lookup("broker-redirect-clipboard"))
			v.BindPFlag("brokerOptions.redirectPNP", cmd.Flags().Lookup("broker-redirect-pnp"))
			v.BindPFlag("brokerOptions.redirectUSB", cmd.Flags().Lookup("broker-redirect-usb"))
			v.BindPFlag("brokerOptions.redirectPrinter", cmd.Flags().Lookup("broker-redirect-printer"))
			v.BindPFlag("brokerOptions.redirectSmartcard", cmd.Flags().Lookup("broker-redirect-smartcard"))
			v.BindPFlag("brokerOptions.smartResize", cmd.Flags().Lookup("broker-smart-resize"))
			v.BindPFlag("brokerOptions.hideCertificateWarnings", cmd.Flags().Lookup("broker-hide-certificate-warnings"))
		}
		if cmd.Flags().Changed("enable-uv") {
			v.BindPFlag("userVolumes.backupSchedule", cmd.Flags().Lookup("uv-backup-schedule"))
			v.BindPFlag("userVolumes.repository", cmd.Flags().Lookup("uv-repository"))
			v.BindPFlag("userVolumes.size", cmd.Flags().Lookup("uv-size"))
			v.BindPFlag("userVolumes.target", cmd.Flags().Lookup("uv-target"))
		}
		var profile rest.Profile
		err := v.Unmarshal(&profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		msg, err := profile.Update(restClient)
		fmt.Println(formatString(msg))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileCreateCmd)
	addProfileFlags(profileCreateCmd)

	profileCmd.AddCommand(profileDeleteCmd)
	profileDeleteCmd.Flags().StringP("id", "i", "", "profile profile Id")
	profileDeleteCmd.Flags().StringP("name", "n", "", "profile profile Name")

	profileCmd.AddCommand(profileGetCmd)
	profileGetCmd.Flags().StringP("id", "i", "", "profile id")
	profileGetCmd.Flags().StringP("name", "n", "", "profile name")

	profileCmd.AddCommand(profileListCmd)
	profileListCmd.Flags().Bool("details", false, "show details")
	profileListCmd.Flags().String("filter", "", "filter query string")

	profileCmd.AddCommand(profileUpdateCmd)
	addProfileFlags(profileUpdateCmd)
}
