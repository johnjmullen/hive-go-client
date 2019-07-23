package rest

import (
	"encoding/json"
	"errors"
)

type PoolDisk struct {
	BootOrder   float64     `json:"bootOrder,omitempty"`
	DiskDriver  interface{} `json:"diskDriver,omitempty"`
	Filename    string      `json:"filename,omitempty"`
	StoragePool string      `json:"storagePool,omitempty"`
	Type        interface{} `json:"type,omitempty"`
}

type PoolGuestProfile struct {
	AgentInstalled bool              `json:"agentInstalled,omitempty"`
	CPU            []int             `json:"cpu,omitempty"`
	Debug          bool              `json:"debug,omitempty"`
	Disks          []*PoolDisk       `json:"disks"`
	Firmware       string            `json:"firmware,omitempty"`
	Gpu            bool              `json:"gpu,omitempty"`
	HostDevices    []*PoolHostDevice `json:"hostDevices,omitempty"`
	Interfaces     []*PoolInterface  `json:"interfaces,omitempty"`
	Mem            []float64         `json:"mem,omitempty"`
	OS             string            `json:"os,omitempty"`
	Persistent     bool              `json:"persistent,omitempty"`
	Protocol       string            `json:"protocol,omitempty"`
	TemplateName   string            `json:"templateName,omitempty"`
	Vga            interface{}       `json:"vga,omitempty"`
}

type PoolHostDevice struct {
	Bus    float64 `json:"bus,omitempty"`
	Domain float64 `json:"domain,omitempty"`
	Func   float64 `json:"func,omitempty"`
	Slot   float64 `json:"slot,omitempty"`
	Type   string  `json:"type,omitempty"`
	UUID   string  `json:"uuid,omitempty"`
}

type PoolInterface struct {
	BootOrder float64     `json:"bootOrder,omitempty"`
	Emulation interface{} `json:"emulation,omitempty"`
	Network   string      `json:"network,omitempty"`
	Vlan      interface{} `json:"vlan,omitempty"`
}

type Pool struct {
	ID                        string            `json:"id,omitempty"`
	Density                   []float64         `json:"density"`
	GuestProfile              *PoolGuestProfile `json:"guestProfile,omitempty"`
	InjectAgent               bool              `json:"injectAgent,omitempty"`
	Name                      string            `json:"name"`
	PerformanceThreshold      float64           `json:"performanceThreshold,omitempty"`
	ProfileID                 interface{}       `json:"profileId,omitempty"`
	Seed                      string            `json:"seed,omitempty"`
	State                     interface{}       `json:"state,omitempty"`
	StorageID                 string            `json:"storageId,omitempty"`
	StorageType               interface{}       `json:"storageType,omitempty"`
	Tags                      []string          `json:"tags,omitempty"`
	TargetState               []interface{}     `json:"targetState,omitempty"`
	Type                      interface{}       `json:"type"`
	UserSessionLoginThreshold float64           `json:"userSessionLoginThreshold,omitempty"`
}

func (pool Pool) String() string {
	json, _ := json.MarshalIndent(pool, "", "  ")
	return string(json)
}

func (client *Client) ListGuestPools() ([]Pool, error) {
	var pools []Pool
	body, err := client.Request("GET", "pools", nil)
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
	body, err := client.Request("GET", "pool/"+id, nil)
	if err != nil {
		return pool, err
	}
	err = json.Unmarshal(body, &pool)
	return pool, err
}

func (client *Client) GetPoolByName(name string) (*Pool, error) {
	var pools, err = client.ListGuestPools()
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
	body, err := client.Request("POST", "pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (pool *Pool) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.Request("PUT", "pool", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (pool *Pool) Delete(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("Invalid pool")
	}
	_, err := client.Request("DELETE", "pool/"+pool.ID, nil)
	return err
}

func (pool *Pool) Refresh(client *Client) error {
	if pool.ID == "" || client == nil {
		return errors.New("Invalid pool")
	}
	_, err := client.Request("POST", "pool/"+pool.ID+"/refresh", nil)
	return err
}
