package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
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
	BrokerOptions  struct {
		Port     uint   `json:"port,omitempty"`
		Protocol string `json:"protocol,omitempty"`
	} `json:"brokerOptions,omitempty"`
}

// PoolHostDevice a hostdevice to share with a virtual machine
type PoolHostDevice struct {
	Bus     int    `json:"bus"`
	Domain  int    `json:"domain"`
	Func    int    `json:"func"`
	Managed bool   `json:"managed"`
	Slot    int    `json:"slot"`
	Type    string `json:"type"`
	UUID    string `json:"uuid,omitempty"`
}

// PoolInterface network interface settings for a pool
type PoolInterface struct {
	Emulation  string      `json:"emulation,omitempty"`
	Network    string      `json:"network,omitempty"`
	Vlan       interface{} `json:"vlan,omitempty"`
	IPAddress  string      `json:"ipAddress"`
	MacAddress string      `json:"macAddress,omitempty"`
}

// PoolBackup data protection settings from a pool record
type PoolBackup struct {
	Enabled         bool   `json:"enabled"`
	Frequency       string `json:"frequency"`
	TargetStorageID string `json:"targetStorageId"`
}

// PoolAffinity host affinity settings for the pool
type PoolAffinity struct {
	CustomCPUFeatures  string   `json:"customCpuFeatures,omitempty"`
	UseHostPassthrough bool     `json:"useHostPassthrough"`
	AllowedHostIDs     []string `json:"allowedHostIds,omitempty"`
}

type PoolAssignmentAutoClear struct {
	Enabled bool `json:"enabled,omitempty"`
	MaxTime int  `json:"maxTime,omitempty"`
}
type PoolAssignment struct {
	Realm     string                   `json:"realm,omitempty"`
	Username  string                   `json:"username,omitempty"`
	ADGroup   string                   `json:"ADGroup,omitempty"`
	AutoClear *PoolAssignmentAutoClear `json:"autoClear,omitempty"`
}

// Pool describes a guest pool record from the rest api
type Pool struct {
	ID                        string            `json:"id,omitempty"`
	Density                   []int             `json:"density"`
	Description               string            `json:"description,omitempty"`
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
	PoolAffinity              *PoolAffinity     `json:"poolAffinity"`
	Assignment                *PoolAssignment   `json:"assignment,omitempty"`
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
		return pool, errors.New("id cannot be empty")
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
	var pools, err = client.ListGuestPools("name=" + url.QueryEscape(name))
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

// Create creates a new pool
func (pool *Pool) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("POST", "pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Update updates an existing pool record
func (pool *Pool) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("PUT", "pool/"+pool.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Delete removes a pool record
func (pool *Pool) Delete(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("invalid pool")
	}
	_, err := client.request("DELETE", "pool/"+pool.ID, nil)
	return err
}

// Refresh refreshes a pool to ensure the definition is applied
func (pool *Pool) Refresh(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("invalid pool")
	}
	_, err := client.request("POST", "pool/"+pool.ID+"/refresh", nil)
	return err
}

// WaitForPool waits for a pool to reach the desired state
func (pool Pool) WaitForPool(client *Client, targetState string, timeout time.Duration) error {
	return pool.WaitForPoolWithContext(context.Background(), client, targetState, timeout)
}

// WaitForPoolWithContext waits for a pool to reach the desired state with a context
func (pool Pool) WaitForPoolWithContext(ctx context.Context, client *Client, targetState string, timeout time.Duration) error {
	if pool.State == targetState {
		return nil
	}
	newVal := Pool{}
	feed, err := client.GetChangeFeedWithContext(ctx, "pool", map[string]string{"id": pool.ID}, false)
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
			return fmt.Errorf("timed out")
		case msg := <-feed.Data:
			if msg.Error != nil {
				feed.Close()
				return msg.Error
			}
			err = json.Unmarshal(msg.NewValue, &newVal)
			if err != nil {
				err = fmt.Errorf("error with json unmarshal: %v", err)
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

// Assign adds a user or group assignment for a standalone pool
func (pool *Pool) Assign(client *Client, realm, username, group string) error {
	if pool.ID == "" || client == nil {
		return errors.New("invalid pool")
	}
	assignment := PoolAssignment{
		Realm:    realm,
		Username: username,
		ADGroup:  group,
	}
	jsonValue, err := json.Marshal(assignment)
	if err != nil {
		return err
	}
	_, err = client.request("POST", "pool/"+pool.ID+"/assignment", jsonValue)
	return err
}

// DeleteAssignment removes the assignment for a standalone pool
func (pool *Pool) DeleteAssignment(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("invalid pool")
	}

	_, err := client.request("DELETE", "pool/"+pool.ID+"/assignment", nil)
	return err
}

// Snapshot stores pool state and creates disk snapshots for running guests
func (pool *Pool) Snapshot(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("invalid pool")
	}
	_, err := client.request("POST", "pool/"+pool.ID+"/snapshot", nil)
	return err
}

// Merge commits the guest snapshots back into their disk
func (pool *Pool) Merge(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("invalid pool")
	}
	_, err := client.request("POST", "pool/"+pool.ID+"/merge", nil)
	return err
}
