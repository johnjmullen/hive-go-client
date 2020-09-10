package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// PoolDisk disk structure for a pool record
type PoolDisk struct {
	BootOrder  int    `json:"bootOrder,omitempty"`
	DiskDriver string `json:"diskDriver,omitempty"`
	Filename   string `json:"filename,omitempty"`
	StorageID  string `json:"storageId,omitempty"`
	Type       string `json:"type,omitempty"`
}

// PoolCloudInit cloud-init settings for a pool record
type PoolCloudInit struct {
	Enabled       bool   `json:"enabled"`
	UserData      string `json:"userData,omitempty"`
	NetworkConfig string `json:"networkConfig,omitempty"`
}

// PoolGuestProfile guestProfile section of a pool record
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

// PoolHostDevice a hostdevice to share with a virtual machine
type PoolHostDevice struct {
	Bus    int    `json:"bus,omitempty"`
	Domain int    `json:"domain,omitempty"`
	Func   int    `json:"func,omitempty"`
	Slot   int    `json:"slot,omitempty"`
	Type   string `json:"type,omitempty"`
	UUID   string `json:"uuid,omitempty"`
}

// PoolInterface network interface settings for a pool
type PoolInterface struct {
	Emulation  string      `json:"emulation,omitempty"`
	Network    string      `json:"network,omitempty"`
	Vlan       interface{} `json:"vlan,omitempty"`
	IPAddress  string `json:"ipAddress"`
	MacAddress string      `json:"macAddress,omitempty"`
}

// PoolBackup data protection settings from a pool record
type PoolBackup struct {
	Enabled         bool   `json:"enabled"`
	Frequency       string `json:"frequency"`
	TargetStorageID string `json:"targetStorageId"`
}

//PoolAffinity host affinity settings for the pool
type PoolAffinity struct {
	CustomCPUFeatures  string   `json:"customCPUFeatures,omitempty"`
	UseHostPassthrough bool     `json:"useHostPassthrough"`
	AllowedHostIDs     []string `json:"allowedHostIds,omitempty"`
}

// Pool describes a guest pool record from the rest api
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
	PoolAffinity              *PoolAffinity     `json:"poolAffinity,omitempty"`
}

func (pool Pool) String() string {
	json, _ := json.MarshalIndent(pool, "", "  ")
	return string(json)
}

// ListGuestPools returns an array of all guest pools with an optional filter string
func (client *Client) ListGuestPools(query string) ([]Pool, error) {
	var pools []Pool
	path := "pools"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return pools, err
	}
	err = json.Unmarshal(body, &pools)
	return pools, err
}

// GetPool request a pool by id
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

// GetPoolByName request a task by name
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

//Create creates a new pool
func (pool *Pool) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("POST", "pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Update updates an existing pool record
func (pool *Pool) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("PUT", "pool/"+pool.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Delete removes a pool record
func (pool *Pool) Delete(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("Invalid pool")
	}
	_, err := client.request("DELETE", "pool/"+pool.ID, nil)
	return err
}

//Refresh refreshes a pool to ensure the definition is applied
func (pool *Pool) Refresh(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("Invalid pool")
	}
	_, err := client.request("POST", "pool/"+pool.ID+"/refresh", nil)
	return err
}

// WaitForPool waits for a pool to reach the desired state
func (pool Pool) WaitForPool(client *Client, targetState string, timeout time.Duration) error {
	if pool.State == targetState {
		return nil
	}
	newVal := Pool{}
	feed, err := client.GetChangeFeed("pool", map[string]string{"id": pool.ID})
	if err != nil {
		return err
	}
	timer := time.NewTimer(timeout)
	if timeout <= 0 && !timer.Stop() {
		<-timer.C
	}
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("Timed out")
		case msg := <-feed.Data:
			if msg.Error != nil {
				feed.Close()
				return msg.Error
			}
			err = json.Unmarshal(msg.NewValue, &newVal)
			if err != nil {
				err = fmt.Errorf("Error with json unmarshal: %v", err)
				feed.Close()
				return err
			}
			if newVal.State == targetState {
				feed.Close()
				return nil
			}
		}
	}
}
