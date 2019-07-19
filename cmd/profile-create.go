package cmd

import (
	"fmt"
	"os"

	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		fmt.Println(formatString(profile))
		msg, err := profile.Create(restClient)
		fmt.Println(msg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var profilePoolFile string

func init() {
	profileCmd.AddCommand(profileCreateCmd)
	profileCreateCmd.Flags().StringP("id", "i", "", "profile id")
	profileCreateCmd.Flags().StringP("name", "n", "", "profile name")
	profileCreateCmd.Flags().String("timezone", "disabled", "timezone to inject")

	profileCreateCmd.Flags().Bool("enable-ad", false, "enable active directory options")
	profileCreateCmd.Flags().String("ad-domain", "", "AD realm")
	profileCreateCmd.Flags().String("ad-username", "", "user to join guests to AD")
	profileCreateCmd.Flags().String("ad-password", "", "Password for the join user")
	profileCreateCmd.Flags().String("ad-user-group", "", "User group")

	profileCreateCmd.Flags().Bool("enable-broker-options", false, "enable broker options")
	profileCreateCmd.Flags().Bool("broker-allow-desktop-composition", true, "Allow Destop compositing")
	profileCreateCmd.Flags().Bool("broker-audio-capture", true, "audio capture")
	profileCreateCmd.Flags().Bool("broker-redirect-cssp", true, "CredSSP redirection")
	profileCreateCmd.Flags().Bool("broker-redirect-clipboard", true, "Clipboard redirection")
	profileCreateCmd.Flags().Bool("broker-redirect-pnp", true, "Plug-and-Play redirection")
	profileCreateCmd.Flags().Bool("broker-redirect-usb", true, "USB redirection")
	profileCreateCmd.Flags().Bool("broker-redirect-printer", true, "Printer redirection")
	profileCreateCmd.Flags().Bool("broker-redirect-smartcard", true, "Smartcard redirection")
	profileCreateCmd.Flags().Bool("broker-smart-resize", true, "Smart screen resize")
	profileCreateCmd.Flags().Bool("broker-hide-certificate-warnings", false, "Hide certificate warnings")

	profileCreateCmd.Flags().Bool("enable-uv", false, "Enable user volumes")
	profileCreateCmd.Flags().Int("uv-backup-schedule", 28800, "User volume backup schedule (s)")
	profileCreateCmd.Flags().String("uv-repository", "", "Storage pool for storing the user volume")
	profileCreateCmd.Flags().Int("uv-size", 10, "User volume size (GB)")
	profileCreateCmd.Flags().String("uv-target", "disk", "Local cache (disk/ram)")
}
