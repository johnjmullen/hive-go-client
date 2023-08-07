package cmd

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hive-io/hive-go-client/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type ExportWriter struct {
	writer *tar.Writer
}

func NewExportWriter(file io.Writer) ExportWriter {
	return ExportWriter{
		writer: tar.NewWriter(file),
	}
}

func (t *ExportWriter) Close() error {
	return t.writer.Close()
}

func (t *ExportWriter) AddFile(filePath string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = t.writer.WriteHeader(&tar.Header{
		Name:  path.Join("export", filePath),
		Mode:  0644,
		Uname: "root",
		Gname: "root",
		Size:  int64(len(data)),
	})
	if err != nil {
		return err
	}

	_, err = t.writer.Write([]byte(data))
	return err
}

type ExportData struct {
	Clusters     []rest.Cluster
	Hosts        []rest.Host
	Realms       []rest.Realm
	Profiles     []rest.Profile
	StoragePools []rest.StoragePool
	Templates    []rest.Template
	Pools        []rest.Pool
	Guests       []rest.Guest
	Broker       rest.Broker
}

func (export ExportData) WriteFile(file *os.File) error {
	exportWriter := NewExportWriter(file)
	include := viper.GetStringSlice("include")

	if slices.Contains(include, "clusters") {
		for _, cluster := range export.Clusters {
			if err := exportWriter.AddFile(path.Join("clusters", cluster.ID), cluster); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "hosts") {
		for _, host := range export.Hosts {
			if err := exportWriter.AddFile(path.Join("hosts", host.Hostid), host); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "realms") {
		for _, realm := range export.Realms {
			if err := exportWriter.AddFile(path.Join("realms", realm.FQDN), realm); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "profiles") {
		for _, profile := range export.Profiles {
			if err := exportWriter.AddFile(path.Join("profiles", profile.ID), profile); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "storagePools") {
		for _, storagePool := range export.StoragePools {
			if err := exportWriter.AddFile(path.Join("storagePools", storagePool.ID), storagePool); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "templates") {
		for _, template := range export.Templates {
			if err := exportWriter.AddFile(path.Join("templates", template.Name), template); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "pools") {
		for _, pool := range export.Pools {
			if err := exportWriter.AddFile(path.Join("pools", pool.ID), pool); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "guests") {
		for _, guest := range export.Guests {
			if err := exportWriter.AddFile(path.Join("guests", guest.Name), guest); err != nil {
				return err
			}
		}
	}
	if slices.Contains(include, "broker") {
		if err := exportWriter.AddFile("broker", export.Broker); err != nil {
			return err
		}
	}
	return exportWriter.Close()
}

func CreateExport() (ExportData, error) {
	export := ExportData{}
	var err error
	if export.Clusters, err = restClient.ListClusters(""); err != nil {
		return export, err
	}

	if export.Hosts, err = restClient.ListHosts(""); err != nil {
		return export, err
	}

	if export.Realms, err = restClient.ListRealms(""); err != nil {
		return export, err
	}

	if export.Profiles, err = restClient.ListProfiles(""); err != nil {
		return export, err
	}

	if export.StoragePools, err = restClient.ListStoragePools(""); err != nil {
		return export, err
	}

	if export.Templates, err = restClient.ListTemplates(""); err != nil {
		return export, err
	}

	if export.Pools, err = restClient.ListGuestPools(""); err != nil {
		return export, err
	}
	if export.Guests, err = restClient.ListGuests(""); err != nil {
		return export, err
	}
	clusterID, err := restClient.ClusterID()
	if err != nil {
		return export, err
	}
	if export.Broker, err = restClient.GetBroker(clusterID); err != nil {
		return export, err
	}
	return export, nil
}

func ReadFile(file *os.File) (ExportData, error) {
	export := ExportData{}
	tr := tar.NewReader(file)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(hdr.Name, "/")
		if hdr.Name == "export/broker" {
			data, err := io.ReadAll(tr)
			if err != nil {
				return export, err
			}
			broker := rest.Broker{}
			err = json.Unmarshal(data, &broker)
			if err != nil {
				return export, err
			}
			export.Broker = broker
		} else if len(parts) < 3 {
			log.Printf("unexpected filename: %s", hdr.Name)
			continue
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return export, err
		}
		switch parts[1] {
		case "realms":
			realm := rest.Realm{}
			err = json.Unmarshal(data, &realm)
			if err != nil {
				return export, err
			}
			export.Realms = append(export.Realms, realm)
		case "profiles":
			profile := rest.Profile{}
			err = json.Unmarshal(data, &profile)
			if err != nil {
				return export, err
			}
			export.Profiles = append(export.Profiles, profile)
		case "clusters":
			cluster := rest.Cluster{}
			err = json.Unmarshal(data, &cluster)
			if err != nil {
				return export, err
			}
			export.Clusters = append(export.Clusters, cluster)
		case "hosts":
			host := rest.Host{}
			err = json.Unmarshal(data, &host)
			if err != nil {
				return export, err
			}
			export.Hosts = append(export.Hosts, host)
		case "templates":
			template := rest.Template{}
			err = json.Unmarshal(data, &template)
			if err != nil {
				return export, err
			}
			export.Templates = append(export.Templates, template)
		case "storagePools":
			sp := rest.StoragePool{}
			err = json.Unmarshal(data, &sp)
			if err != nil {
				return export, err
			}
			export.StoragePools = append(export.StoragePools, sp)
		case "pools":
			pool := rest.Pool{}
			err = json.Unmarshal(data, &pool)
			if err != nil {
				return export, err
			}
			export.Pools = append(export.Pools, pool)
		case "guests":
			guest := rest.Guest{}
			err = json.Unmarshal(data, &guest)
			if err != nil {
				return export, err
			}
			//only include external guests for now
			if guest.External {
				export.Guests = append(export.Guests, guest)
			}
		}
	}
	return export, nil
}

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "export cluster configuration",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("include", cmd.Flags().Lookup("include"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		var err error
		if args[0] == "-" {
			file = os.Stdout
		} else {
			file, err = os.Create(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
		}

		export, err := CreateExport()
		if err != nil {
			log.Fatalln(err)
		}
		export.WriteFile(file)
	},
}

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "import configuration from an export file",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("include", cmd.Flags().Lookup("include"))
		viper.BindPFlag("enable-shared-storage", cmd.Flags().Lookup("enable-shared-storage"))
		viper.BindPFlag("create-cluster", cmd.Flags().Lookup("create-cluster"))
		viper.BindPFlag("enable-pools", cmd.Flags().Lookup("enable-pools"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var file *os.File
		var err error
		if args[0] == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
		}
		data, err := ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		clusterID, err := restClient.ClusterID()
		if err != nil {
			log.Fatalln(err)
		}
		cluster, err := restClient.GetCluster(clusterID)
		if err != nil {
			log.Fatalln(err)
		}
		if cluster.License == nil {
			fmt.Println("License not found")
			license := ""
			prompt := &survey.Input{
				Message: fmt.Sprintf("Please enter a license key for cluster id %s:", clusterID),
			}
			survey.AskOne(prompt, &license)
			if err := cluster.SetLicense(restClient, license); err != nil {
				log.Fatalln(err)
			}
		}
		include := viper.GetStringSlice("include")
		hostidMap := map[string]string{}
		if viper.GetBool("create-cluster") {
			viper.Set("wait", true)
			viper.Set("progress-bar", true)

			for _, host := range data.Hosts {
				if hostRecord, err := restClient.GetHostByIP(host.IP); err == nil {
					hostidMap[host.Hostid] = hostRecord.Hostid
					continue //already exists
				}
				task, err := restClient.JoinHost(viper.GetString("user"), viper.GetString("password"), host.IP)
				if err != nil {
					log.Printf("Failed to add Host %s to the cluster\n", host.IP)
					continue
				}
				err = waitForTask(task, false, true)
				if err != nil {
					log.Printf("Failed to add Host %s to the cluster\n", host.IP)
					continue
				}
			}
		}
		for _, host := range data.Hosts {
			hostRecord, err := restClient.GetHostByIP(host.IP)
			if err != nil {
				continue
			}
			hostidMap[host.Hostid] = hostRecord.Hostid
			if slices.Contains(include, "hosts") {
				if host.State != hostRecord.State && host.Appliance.Role != "gateway" {
					task, err := hostRecord.SetState(restClient, host.State)
					if err != nil {
						log.Printf("Failed to set state for Host %s (%s)\n", host.Hostname, host.IP)
						continue
					}
					err = waitForTask(task, false, true)
					if err != nil {
						log.Printf("Failed to set state for Host %s (%s)\n", host.Hostname, host.IP)
					}
				}
				updateApplianceConf := false
				if host.Appliance.Loglevel != hostRecord.Appliance.Loglevel {
					updateApplianceConf = true
					hostRecord.Appliance.Loglevel = host.Appliance.Loglevel
				}
				if host.Appliance.Ntp != hostRecord.Appliance.Ntp {
					updateApplianceConf = true
					hostRecord.Appliance.Ntp = host.Appliance.Ntp
				}
				if host.Appliance.MaxCloneDensity != hostRecord.Appliance.MaxCloneDensity {
					updateApplianceConf = true
					hostRecord.Appliance.MaxCloneDensity = host.Appliance.MaxCloneDensity
				}
				if updateApplianceConf {
					log.Printf("Updating Host settings for %s (%s)\n", host.Hostname, host.IP)
					_, err = hostRecord.UpdateAppliance(restClient)
					if err != nil {
						log.Printf("Failed to update settings for Host %s (%s)\n", host.Hostname, host.IP)
					}
				}
			}
		}
		//shared storage
		oldSharedStorageId := ""
		newSharedStorageId := ""
		if len(data.Clusters) > 0 && data.Clusters[0].SharedStorage != nil && data.Clusters[0].SharedStorage.Enabled {
			oldSharedStorageId = data.Clusters[0].SharedStorage.ID
			if cluster.SharedStorage == nil || !cluster.SharedStorage.Enabled && viper.GetBool("enable-shared-storage") {
				task, err := cluster.EnableSharedStorage(restClient, data.Clusters[0].SharedStorage.StorageUtilization, data.Clusters[0].SharedStorage.MinSetSize)
				if err != nil {
					log.Printf("Failed to enable shared storage\n")
				}
				if err = waitForTask(task, false, true); err != nil {
					log.Printf("Enable shared storage task failed\n")
				}
			}
			cluster, err = restClient.GetCluster(clusterID)
			if err != nil {
				log.Fatalln(err)
			}
			newSharedStorageId = cluster.SharedStorage.ID
		}

		if slices.Contains(include, "realms") {
			for _, realm := range data.Realms {
				if _, err := restClient.GetRealm(realm.Name); err == nil {
					continue //already exists
				}
				fmt.Printf("Adding realm %s\n", realm.Name)
				if realm.ServiceAccount != nil {
					prompt := &survey.Password{
						Message: fmt.Sprintf("Password required for %s@%s:", realm.ServiceAccount.Username, realm.FQDN),
					}
					survey.AskOne(prompt, &realm.ServiceAccount.Password)
				}
				_, err := realm.Create(restClient)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
		if slices.Contains(include, "storagePools") {
			for _, storagePool := range data.StoragePools {
				if storagePool.Name == "HF_Shared" {
					continue //TODO: enable shared storage and update storageIds
				}
				if _, err := restClient.GetStoragePool(storagePool.ID); err == nil {
					continue //already exists
				}
				fmt.Printf("Adding storage pool %s\n", storagePool.Name)
				_, err := storagePool.Create(restClient)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}

		if slices.Contains(include, "profiles") {
			for _, profile := range data.Profiles {
				if _, err := restClient.GetProfile(profile.ID); err == nil {
					continue //already exists
				}
				if profile.UserVolumes != nil && profile.UserVolumes.Repository == oldSharedStorageId {
					if newSharedStorageId == "" {
						log.Printf("Shared Storage not found. Skipping profile %s\n", profile.Name)
						continue
					}
					profile.UserVolumes.Repository = newSharedStorageId
				}
				fmt.Printf("Adding profile %s\n", profile.Name)
				if profile.AdConfig != nil && profile.AdConfig.Username != "" {
					prompt := &survey.Password{
						Message: fmt.Sprintf("Password required for service account user %s:", profile.AdConfig.Username),
					}
					survey.AskOne(prompt, &profile.AdConfig.Password)
				}
				_, err := profile.Create(restClient)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
		if slices.Contains(include, "broker") {
			if err = restClient.SetBroker(clusterID, data.Broker); err != nil {
				log.Fatalln(err)
			}
		}

		if len(data.Clusters) > 0 && data.Clusters[0].Gateway != nil {
			gwHostInfo := map[string]rest.GatewayHost{}
			for id, host := range data.Clusters[0].Gateway.Hosts {
				if hostid, ok := hostidMap[id]; ok {
					gwHostInfo[hostid] = host
				}
			}
			data.Clusters[0].Gateway.Hosts = gwHostInfo
			if len(gwHostInfo) > 0 {
				if err = restClient.SetGateway(clusterID, *data.Clusters[0].Gateway); err != nil {
					log.Fatalln(err)
				}
			}
		}

		if slices.Contains(include, "templates") {
			for _, template := range data.Templates {
				if _, err := restClient.GetTemplate(template.Name); err == nil {
					continue //already exists
				}
				foundDisks := true
				for i, disk := range template.Disks {
					if disk.StorageID == oldSharedStorageId {
						if newSharedStorageId == "" {
							log.Printf("Shared Storage id not found. Skipping template %s\n", template.Name)
							foundDisks = false
							break
						}
						template.Disks[i].StorageID = newSharedStorageId
					}
					sp, err := restClient.GetStoragePool(template.Disks[i].StorageID)
					if err != nil {
						log.Printf("Storage Pool %s not found for template %s\n", disk.StorageID, template.Name)
						foundDisks = false
						break
					}
					if _, err := sp.DiskInfo(restClient, disk.Filename); err != nil {
						log.Printf("Disk %s not found for template: %s\n", disk.Filename, template.Name)
						foundDisks = false
						break
					}
				}
				if !foundDisks {
					fmt.Printf("Skipping template %s\n", template.Name)
					continue
				}
				fmt.Printf("Adding template %s\n", template.Name)
				_, err := template.Create(restClient)
				if err != nil {
					log.Printf("Error adding template: %v\n", err)
					continue
				}
				time.Sleep(time.Second * 10)
				template, err = restClient.GetTemplate(template.Name)
				if err != nil {
					log.Println(err)
				} else if template.State != "available" {
					log.Println(template.StateMessage)
				}
			}
		}
		if slices.Contains(include, "pools") {
			for _, pool := range data.Pools {
				if _, err := restClient.GetPool(pool.ID); err == nil {
					continue //already exists
				}
				if pool.GuestProfile == nil {
					continue //skip if guestProfile is missing
				}
				if pool.StorageID == oldSharedStorageId {
					if newSharedStorageId == "" {
						log.Printf("Shared Storage id not found. Skipping pool %s\n", pool.Name)
						continue
					}
					pool.StorageID = newSharedStorageId
				}
				foundDisks := true
				if pool.GuestProfile.TemplateName != "" {
					if _, err := restClient.GetTemplate(pool.GuestProfile.TemplateName); err != nil {
						log.Printf("Template %s not found.  Skipping pool %s\n", pool.GuestProfile.TemplateName, pool.Name)
						continue
					}
				}
				for i, disk := range pool.GuestProfile.Disks {
					if disk.StorageID == oldSharedStorageId {
						if newSharedStorageId == "" {
							log.Printf("Shared Storage id not found. Skipping pool %s\n", pool.Name)
							foundDisks = false
							break
						}
						pool.GuestProfile.Disks[i].StorageID = newSharedStorageId
					}
					sp, err := restClient.GetStoragePool(pool.GuestProfile.Disks[i].StorageID)
					if err != nil {
						log.Printf("Storage Pool %s not found for pool: %s\n", disk.StorageID, pool.Name)
						foundDisks = false
					}
					if _, err := sp.DiskInfo(restClient, disk.Filename); err != nil {
						log.Printf("Disk %s not found for pool: %s\n", disk.Filename, pool.Name)
						foundDisks = false
						break
					}
				}
				if !foundDisks {
					log.Printf("Skipping guest pool %s\n", pool.Name)
					continue
				}
				fmt.Printf("Adding guest pool %s\n", pool.Name)
				if !viper.GetBool("enable-pools") {
					pool.State = "disabled"
				}
				_, err := pool.Create(restClient)
				if err != nil {
					log.Printf("Error adding pool: %v\n", err)
				}
			}
		}
		if slices.Contains(include, "guests") {
			for _, guest := range data.Guests {
				if _, err := restClient.GetGuest(guest.Name); err == nil {
					continue //already exists
				}

				fmt.Printf("Adding external guest %s\n", guest.Name)
				externalGuest := rest.ExternalGuest{
					GuestName: guest.Name,
					Address:   guest.Address,
					Username:  guest.Username,
					ADGroup:   guest.ADGroup,
					Realm:     guest.Realm,
					OS:        guest.Os,
				}
				_, err := externalGuest.Create(restClient)
				if err != nil {
					log.Printf("Error adding external guest: %v\n", err)
				}
			}
		}

	},
}

func init() {
	importCmd.Flags().StringArray("include", []string{"broker", "clusters", "hosts", "guests", "pools", "profiles", "realms", "storagePools", "templates"}, "Data to import from the export file")
	importCmd.Flags().Bool("enable-shared-storage", false, "Automatically create shared storage")
	importCmd.Flags().Bool("create-cluster", false, "Automatically add hosts from the export file to the cluster")
	importCmd.Flags().Bool("enable-pools", false, "Enable guest pools automatically")
	RootCmd.AddCommand(importCmd)

	exportCmd.Flags().StringArray("include", []string{"broker", "clusters", "guests", "hosts", "pools", "profiles", "realms", "storagePools", "templates"}, "Data to include in the export file")
	RootCmd.AddCommand(exportCmd)

}
