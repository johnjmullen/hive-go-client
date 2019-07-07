package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Pool struct {
	Density   []int `json:"density"`
	GuestPool struct {
		CPU          []int  `json:"cpu"`
		Mem          []int  `json:"mem"`
		Os           string `json:"os"`
		Persistent   bool   `json:"persistent"`
		TemplateName string `json:"templateName"`
		Vga          string `json:"vga"`
	} `json:"guestPool"`
	ID          string   `json:"id"`
	InjectAgent bool     `json:"injectAgent"`
	Name        string   `json:"name"`
	PoolID      string   `json:"poolId"`
	Seed        string   `json:"seed"`
	State       string   `json:"state"`
	StorageID   string   `json:"storageId"`
	StorageType string   `json:"storageType"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
}

func (pool Pool) String() string {
	json, _ := json.MarshalIndent(pool, "", "  ")
	return string(json)
}

func (pool *Pool) ToJson() ([]byte, error) {
	return json.Marshal(pool)
}

func (pool *Pool) FromJson(data []byte) error {
	return json.Unmarshal(data, pool)
}

func (client *Client) ListPools() ([]Pool, error) {
	var Pools []Pool
	body, err := client.Request("GET", "pools", nil)
	if err != nil {
		return Pools, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &Pools)
	return Pools, err
}

func (client *Client) GetPool(id string) (Pool, error) {
	var Pool Pool
	if id == "" {
		return Pool, errors.New("Id cannot be empty")
	}
	body, err := client.Request("GET", "pool/"+id, nil)
	if err != nil {
		return Pool, err
	}
	err = json.Unmarshal(body, &Pool)
	return Pool, err
}

func (client *Client) CreatePool(Pool *Pool) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(Pool)
	body, err := client.Request("POST", "pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (client *Client) DeletePool(id string) error {
	if id == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.Request("DELETE", "pool/"+id, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) GetPoolByName(name string) (*Pool, error) {
	var pools, err = client.ListPools()
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
