package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Host describes a host record from the rest api
type Host struct {
	Appliance struct {
		Broker           bool   `json:"broker"`
		ClusterID        string `json:"clusterId"`
		Cma              string `json:"cma"`
		CPUGovernor      string `json:"cpuGovernor"`
		CRS              bool   `json:"crs"`
		DbName           string `json:"dbName"`
		DoNotParticipate bool   `json:"doNotParticipate"`
		Firmware         struct {
			Active      string `json:"active"`
			Software    string `json:"software"`
			PendingSwap bool   `json:"pendingSwap"`
		} `json:"firmware"`
		Hostname        string `json:"hostname"`
		Loglevel        string `json:"loglevel"`
		MaxCloneDensity int    `json:"maxCloneDensity"`
		Ntp             string `json:"ntp"`
		Timezone        string `json:"timezone"`
		Role            string `json:"role"`
	} `json:"appliance"`
	Capabilities struct {
		StorageTypes []string `json:"storageTypes"`
	} `json:"capabilities"`
	Certificate struct {
		Expiration interface{} `json:"expiration"`
		Issuer     string      `json:"issuer"`
		Status     string      `json:"status"`
	} `json:"certificate"`
	Hardware struct {
		BIOS struct {
			ReleaseDate string `json:"releaseDate"`
			Vender      string `json:"vender"`
			Version     string `json:"version"`
		} `json:"BIOS"`
		Chasis struct {
			AssetTag     string `json:"assetTag"`
			Manufacturer string `json:"manufacturer"`
			Version      string `json:"version"`
		} `json:"Chasis"`
		HyperThreadingEnabled bool `json:"HyperThreadingEnabled"`
		PhysicalCPUs          int  `json:"PhysicalCPUs"`
		PhysicalCoresPerCPU   int  `json:"PhysicalCoresPerCPU"`
		Processor             []struct {
			Cores        int    `json:"cores"`
			Family       string `json:"family"`
			Manufacturer string `json:"manufacturer"`
			Threads      int    `json:"threads"`
			Version      string `json:"version"`
		} `json:"Processor"`
		System struct {
			Manufacturer string `json:"manufacturer"`
			ProductName  string `json:"productName"`
		} `json:"System"`
		TotalPhysicalMemory int `json:"TotalPhysicalMemory"`
		CPUFeatures         interface{}
		Memory              []struct {
			Size int    `json:"size"`
			Type string `json:"type"`
		} `json:"memory"`
		VideoCards []struct {
			Bus         int    `json:"bus"`
			DeviceClass int    `json:"deviceClass"`
			DeviceID    int    `json:"deviceId"`
			Domain      int    `json:"domain"`
			Func        int    `json:"func"`
			IommuGroup  int    `json:"iommu_group"`
			Mode        string `json:"mode"`
			Path        string `json:"path"`
			Slot        int    `json:"slot"`
			VendorID    int    `json:"vendorId"`
		} `json:"videoCards"`
	} `json:"hardware"`
	Hostid     string `json:"hostid"`
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	Networking struct {
		Interfaces struct {
		} `json:"interfaces"`
		Production struct {
			Dhcp      bool   `json:"dhcp"`
			DNS       string `json:"dns"`
			Interface string `json:"interface"`
			Search    string `json:"search"`
			Vlan      int    `json:"vlan"`
		} `json:"production"`
		Storage struct {
			Interface string `json:"interface"`
			IP        string `json:"ip"`
			Mask      string `json:"mask"`
			Vlan      int    `json:"vlan"`
		} `json:"storage"`
	} `json:"networking"`
	RdbID    string      `json:"rdbId"`
	Software interface{} `json:"software"`
	State    string      `json:"state"`
	Storage  struct {
		Blockdevices []struct {
			MajMin     string `json:"maj:min"`
			Mountpoint string `json:"mountpoint"`
			Name       string `json:"name"`
			Rm         string `json:"rm"`
			Ro         string `json:"ro"`
			Size       string `json:"size"`
			Type       string `json:"type"`
			Children   []struct {
				MajMin     string `json:"maj:min"`
				Mountpoint string `json:"mountpoint"`
				Name       string `json:"name"`
				Rm         string `json:"rm"`
				Ro         string `json:"ro"`
				Size       string `json:"size"`
				Type       string `json:"type"`
			} `json:"children,omitempty"`
		} `json:"blockdevices"`
		Disk struct {
			Zpool struct {
				CacheDevices []string `json:"cacheDevices"`
				Compression  string   `json:"compression"`
				Dedup        string   `json:"dedup"`
				Devices      []string `json:"devices"`
				Filesystems  struct {
					Hive struct {
						Filesystems struct {
							Conf struct {
								Mountpoint string `json:"mountpoint"`
							} `json:"conf"`
						} `json:"filesystems"`
						Mountpoint string `json:"mountpoint"`
					} `json:"hive"`
					Rethink struct {
						Mountpoint  string `json:"mountpoint"`
						Reservation string `json:"reservation"`
					} `json:"rethink"`
					Root struct {
						Mountpoint  string `json:"mountpoint"`
						Reservation string `json:"reservation"`
					} `json:"root"`
					Zdata struct {
					} `json:"zdata"`
				} `json:"filesystems"`
				Mountpoint string `json:"mountpoint"`
				Volumes    struct {
				} `json:"volumes"`
			} `json:"zpool"`
		} `json:"disk"`
		RAM struct {
			RamdiskPercent int `json:"ramdiskPercent"`
			SwapSize       int `json:"swapSize"`
			Swappiness     int `json:"swappiness"`
			Zram           struct {
				Compression string `json:"compression"`
				Dedup       string `json:"dedup"`
				Filesystems struct {
					Users struct {
						Mountpoint string `json:"mountpoint"`
					} `json:"users"`
				} `json:"filesystems"`
				Mountpoint string `json:"mountpoint"`
			} `json:"zram"`
		} `json:"ram"`
	} `json:"storage"`
	Tags []string `json:"tags"`
}

func (host Host) String() string {
	json, _ := json.MarshalIndent(host, "", "  ")
	return string(json)
}

// ListHosts returns an array of all host with an optional filter string
func (client *Client) ListHosts(query string) ([]Host, error) {
	var hosts []Host
	path := "hosts"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return hosts, err
	}
	err = json.Unmarshal(body, &hosts)
	return hosts, err
}

// GetHost requests a single guest by hostid
func (client *Client) GetHost(hostid string) (Host, error) {
	var host Host
	if hostid == "" {
		return host, errors.New("hostid cannot be empty")
	}
	body, err := client.request("GET", "host/"+hostid, nil)
	if err != nil {
		return host, err
	}
	err = json.Unmarshal(body, &host)
	return host, err
}

// GetHostByName requests a host by hostname
func (client *Client) GetHostByName(name string) (*Host, error) {
	var hosts, err = client.ListHosts("hostname=" + url.QueryEscape(name))
	if err != nil {
		return nil, err
	}
	for _, host := range hosts {
		if host.Hostname == name {
			return &host, nil
		}
	}
	return nil, errors.New("Host not found")
}

// GetHostByIP requests a host by hostname
func (client *Client) GetHostByIP(ip string) (*Host, error) {
	var hosts, err = client.ListHosts("ip=" + ip)
	if err != nil {
		return nil, err
	}
	for _, host := range hosts {
		if host.IP == ip {
			return &host, nil
		}
	}
	return nil, errors.New("Host not found")
}

//UpdateAppliance updates settings from Host.appliance
func (host *Host) UpdateAppliance(client *Client) (string, error) {
	var result string
	data := map[string]interface{}{"appliance": host.Appliance}
	jsonValue, _ := json.Marshal(data)
	body, err := client.request("PUT", "host/"+host.Hostid, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Delete removes a host from the database
func (host *Host) Delete(client *Client) error {
	if host.Hostid == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.request("DELETE", "host/"+host.Hostid, nil)
	if err != nil {
		return err
	}
	return err
}

// RestartServices calls restart hive services
func (host *Host) RestartServices(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/services/hive-services/restart", nil)
	return err
}

// Reboot calls reboot on a host
func (host *Host) Reboot(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/system/reboot", nil)
	return err
}

// Shutdown calls shutdown on a host
func (host *Host) Shutdown(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/system/shutdown", nil)
	return err
}

// Version is a structure containing version information returned by HostVersion
type Version struct {
	Major   uint   `json:"major"`
	Minor   uint   `json:"minor"`
	Patch   uint   `json:"patch"`
	Version string `json:"version"`
}

// SetState can be used to set a host's state to available or maintenance
func (host *Host) SetState(client *Client, state string) (*Task, error) {
	return client.getTaskFromResponse(client.request("POST", "host/"+host.Hostid+"/state?state="+state, nil))
}

// GetState gets the current state of the host
func (host *Host) GetState(client *Client) (string, error) {
	body, err := client.request("GET", "host/"+host.Hostid+"/state", nil)
	if err != nil {
		return "", err
	}
	var state string
	err = json.Unmarshal(body, &state)
	return state, err
}

// UnjoinCluster removes a host from the cluster
func (host *Host) UnjoinCluster(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/cluster/unjoin", nil)
	return err
}

//HostPackageInfo contains information about software and firmware packages
type HostPackageInfo struct {
	Packages []string `json:"packages"`
	Current  string   `json:"current"`
}

//ListSoftware returns the current software version and available packages
func (host *Host) ListSoftware(client *Client) (HostPackageInfo, error) {
	var info HostPackageInfo
	body, err := client.request("GET", "host/"+host.Hostid+"/firmware/software/packages", nil)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(body, &info)
	return info, err
}

//DeleteSoftware deletes a software package from a host
func (host *Host) DeleteSoftware(client *Client, pkg string) error {
	_, err := client.request("DELETE", "host/"+host.Hostid+"/firmware/software/"+pkg, nil)
	return err
}

//UploadSoftware uploads a firmware pkg file to the host
func (host *Host) UploadSoftware(client *Client, filename string) error {
	minVersion, _ := version.NewVersion("8.2.6-1110")
	v, err := version.NewVersion(strings.TrimPrefix(host.Appliance.Firmware.Software, "hiveio-fabric-v"))
	if err != nil || v.LessThan(minVersion) {
		_, err := client.postMultipart(fmt.Sprintf("host/%s/firmware/software/upload", host.Hostid), "data", filename, nil)
		return err
	}
	sp := StoragePool{ID: "softwarePackage"}
	return sp.Upload(client, filename, path.Base(filename))
}

// RestartNetworking calls restarts networking on the host
func (host *Host) RestartNetworking(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/networking/networking/restart", nil)
	return err
}

// EnableCRS enables crs on the host
func (host *Host) EnableCRS(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/enableCRS", nil)
	return err
}

// DisableCRS disables crs on the host
func (host *Host) DisableCRS(client *Client) error {
	_, err := client.request("POST", "host/"+host.Hostid+"/disableCRS", nil)
	return err
}
