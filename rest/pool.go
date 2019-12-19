package rest

import (
	"encoding/json"
	"errors"
)

type PoolDisk struct {
	BootOrder  int    `json:"bootOrder,omitempty"`
	DiskDriver string `json:"diskDriver,omitempty"`
	Filename   string `json:"filename,omitempty"`
	StorageID  string `json:"storageId,omitempty"`
	Type       string `json:"type,omitempty"`
}

type PoolCloudInit struct {
	Enabled       bool   `json:"enabled"`
	UserData      string `json:"userData,omitempty"`
	NetworkConfig string `json:"networkConfig,omitempty"`
}

type PoolGuestProfile struct {
	AgentInstalled bool              `json:"agentInstalled"`
	CPU            []int             `json:"cpu,omitempty"`
	Debug          bool              `json:"debug,omitempty"`
	Disks          []*PoolDisk       `json:"disks,omitempty"`
	Firmware       string            `json:"firmware,omitempty"`
	Gpu            bool              `json:"gpu"`
	HostDevices    []*PoolHostDevice `json:"hostDevices,omitempty"`
	Interfaces     []*PoolInterface  `json:"interfaces,omitempty"`
	Mem            []int             `json:"mem,omitempty"`
	OS             string            `json:"os,omitempty"`
	Persistent     bool              `json:"persistent"`
	Protocol       string            `json:"protocol,omitempty"`
	TemplateName   string            `json:"templateName,omitempty"`
	Vga            string            `json:"vga,omitempty"`
	CloudInit      *PoolCloudInit    `json:"cloudInit,omitempty"`
}

type PoolHostDevice struct {
	Bus    int    `json:"bus,omitempty"`
	Domain int    `json:"domain,omitempty"`
	Func   int    `json:"func,omitempty"`
	Slot   int    `json:"slot,omitempty"`
	Type   string `json:"type,omitempty"`
	UUID   string `json:"uuid,omitempty"`
}

type PoolInterface struct {
	BootOrder  int    `json:"bootOrder,omitempty"`
	Emulation  string `json:"emulation,omitempty"`
	Network    string `json:"network,omitempty"`
	Vlan       int    `json:"vlan,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
}

type PoolBackup struct {
	Frequency       string `json:"frequency"`
	TargetStorageID string `json:"targetStorageId"`
}

type Pool struct {
	ID                        string            `json:"id,omitempty"`
	Density                   []int             `json:"density"`
	GuestProfile              *PoolGuestProfile `json:"guestProfile,omitempty"`
	InjectAgent               bool              `json:"injectAgent"`
	Name                      string            `json:"name"`
	PerformanceThreshold      int               `json:"performanceThreshold,omitempty"`
	ProfileID                 string            `json:"profileId,omitempty"`
	Seed                      string            `json:"seed,omitempty"`
	State                     string            `json:"state,omitempty"`
	StorageID                 string            `json:"storageId,omitempty"`
	StorageType               string            `json:"storageType,omitempty"`
	Tags                      []string          `json:"tags,omitempty"`
	TargetState               []string          `json:"targetState,omitempty"`
	Type                      string            `json:"type"`
	UserSessionLoginThreshold int               `json:"userSessionLoginThreshold,omitempty"`
	Backup                    *PoolBackup       `json:"backup,omitempty"`
}

func (pool Pool) String() string {
	json, _ := json.MarshalIndent(pool, "", "  ")
	return string(json)
}

func (client *Client) ListGuestPools(filter string) ([]Pool, error) {
	var pools []Pool
	path := "pools"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return pools, err
	}
	err = json.Unmarshal(body, &pools)
	return pools, err
}

func (client *Client) GetPool(id string) (*Pool, error) {
	var pool *Pool
	if id == "" {
		return pool, errors.New("Id cannot be empty")
	}
	body, err := client.request("GET", "pool/"+id, nil)
	if err != nil {
		return pool, err
	}
	err = json.Unmarshal(body, &pool)
	return pool, err
}

func (client *Client) GetPoolByName(name string) (*Pool, error) {
	var pools, err = client.ListGuestPools("name=" + name)
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		if pool.Name == name {
			return &pool, nil
		}
	}
	return nil, errors.New("Pool not found")
}

func (pool *Pool) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("POST", "pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (pool *Pool) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("PUT", "pool/"+pool.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (pool *Pool) Delete(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("Invalid pool")
	}
	_, err := client.request("DELETE", "pool/"+pool.ID, nil)
	return err
}

func (pool *Pool) Refresh(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("Invalid pool")
	}
	_, err := client.request("POST", "pool/"+pool.ID+"/refresh", nil)
	return err
}
