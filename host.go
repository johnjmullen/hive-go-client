package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Host struct {
	Appliance struct {
		Broker      bool   `json:"broker"`
		ClusterID   string `json:"clusterId"`
		Cma         string `json:"cma"`
		CPUGovernor string `json:"cpuGovernor"`
		DbName      string `json:"dbName"`
		Firmware    struct {
			Active      string `json:"active"`
			Software    string `json:"software"`
			PendingSwap bool   `json:"pendingSwap"`
		} `json:"firmware"`
		Hostname        string      `json:"hostname"`
		Loglevel        string      `json:"loglevel"`
		MaxCloneDensity int         `json:"maxCloneDensity"`
		Ntp             interface{} `json:"ntp"`
		Timezone        string      `json:"timezone"`
	} `json:"appliance"`
	Capabilities struct {
		StorageTypes []string `json:"storageTypes"`
	} `json:"capabilities"`
	Certificate struct {
		Expiration int64  `json:"expiration"`
		Issuer     string `json:"issuer"`
		Status     string `json:"status"`
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
		CPUFeatures         struct {
			Arch     string `json:"arch"`
			Features []struct {
				Name   string      `json:"name"`
				Policy interface{} `json:"policy"`
			} `json:"features"`
			ModelName string `json:"modelName"`
		} `json:"cpuFeatures"`
		Memory []struct {
			Size int    `json:"size"`
			Type string `json:"type"`
		} `json:"memory"`
		VideoCards []interface{} `json:"videoCards"`
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
	RdbID   string `json:"rdbId"`
	State   string `json:"state"`
	Storage struct {
		Blockdevices []struct {
			MajMin     string      `json:"maj:min"`
			Mountpoint interface{} `json:"mountpoint"`
			Name       string      `json:"name"`
			Rm         string      `json:"rm"`
			Ro         string      `json:"ro"`
			Size       string      `json:"size"`
			Type       string      `json:"type"`
			Children   []struct {
				MajMin     string      `json:"maj:min"`
				Mountpoint interface{} `json:"mountpoint"`
				Name       string      `json:"name"`
				Rm         string      `json:"rm"`
				Ro         string      `json:"ro"`
				Size       string      `json:"size"`
				Type       string      `json:"type"`
			} `json:"children,omitempty"`
		} `json:"blockdevices"`
		Disk struct {
			Zpool struct {
				CacheDevices []interface{} `json:"cacheDevices"`
				Compression  string        `json:"compression"`
				Dedup        string        `json:"dedup"`
				Devices      []string      `json:"devices"`
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
	CRS  bool     `json:"CRS"`
}

func (host *Host) String() string {
	json, _ := json.MarshalIndent(host, "", "  ")
	return string(json)
}

func (client *Client) ListHosts() ([]Host, error) {
	var hosts []Host
	res, err := client.Request("GET", "hosts", nil)
	if err != nil {
		return hosts, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return hosts, err
	}
	err = json.Unmarshal(body, &hosts)
	return hosts, err
}

func (client *Client) GetHost(hostid string) (Host, error) {
	var host Host
	if hostid == "" {
		return host, errors.New("hostid cannot be empty")
	}
	res, err := client.Request("GET", "host/"+hostid, nil)
	if err != nil {
		return host, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return host, err
	}
	err = json.Unmarshal(body, &host)
	return host, err
}

func (client *Client) DeleteHost(hostid string) error {
	if hostid == "" {
		return errors.New("name cannot be empty")
	}
	res, err := client.Request("DELETE", "host/"+hostid, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return err
}
