package rest

import (
	"encoding/json"
	"errors"
)

type Pool struct {
	Density      []int `json:"density"`
	GuestProfile struct {
		CPU          []int  `json:"cpu"`
		Mem          []int  `json:"mem"`
		Os           string `json:"os"`
		Persistent   bool   `json:"persistent"`
		TemplateName string `json:"templateName"`
		Vga          string `json:"vga"`
	} `json:"guestProfile"`
	ID          string   `json:"id,omitempty"`
	InjectAgent bool     `json:"injectAgent"`
	Name        string   `json:"name"`
	ProfileID   string   `json:"profileId"`
	Seed        string   `json:"seed"`
	State       string   `json:"state,omitempty"`
	StorageID   string   `json:"storageId"`
	StorageType string   `json:"storageType"`
	Tags        []string `json:"tags,omitempty"`
	Type        string   `json:"type"`
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
