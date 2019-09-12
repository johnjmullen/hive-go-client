package rest

import (
	"encoding/json"
	"errors"
	"fmt"
)

type StoragePool struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Server       string   `json:"server"`
	Path         string   `json:"path"`
	Username     string   `json:"username,omitempty"`
	Password     string   `json:"password,omitempty"`
	Key          string   `json:"key,omitempty"`
	MountOptions []string `json:"mountOptions,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

func (sp StoragePool) String() string {
	json, _ := json.MarshalIndent(sp, "", "  ")
	return string(json)
}

func (client *Client) ListStoragePools(filter string) ([]StoragePool, error) {
	var pools []StoragePool
	path := "storage/pools"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.Request("GET", path, nil)
	if err != nil {
		return pools, err
	}
	err = json.Unmarshal(body, &pools)
	return pools, err
}

func (client *Client) GetStoragePoolByName(name string) (*StoragePool, error) {
	var pools, err = client.ListStoragePools("name=" + name)
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		if pool.Name == name {
			return &pool, nil
		}
	}
	return nil, errors.New("Storage Pool not found")
}

func (client *Client) GetStoragePool(id string) (*StoragePool, error) {
	var pool *StoragePool
	if id == "" {
		return pool, errors.New("id cannot be empty")
	}
	body, err := client.Request("GET", "storage/pool/"+id, nil)
	if err != nil {
		return pool, err
	}
	err = json.Unmarshal(body, &pool)
	return pool, err
}

func (pool *StoragePool) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.Request("POST", "storage/pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (pool *StoragePool) Delete(client *Client) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	_, err := client.Request("DELETE", "storage/pool/"+pool.ID, nil)
	return err
}

func (pool *StoragePool) CreateDisk(client *Client, filename, format string, size uint) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	jsonData := map[string]interface{}{"filename": filename, "size": size, "format": format}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.Request("POST", "storage/pool/"+pool.ID+"/createDisk", jsonValue)
	if err != nil {
		return err
	}
	return err
}

func (pool *StoragePool) ConvertDisk(client *Client, srcFilename, dstStorageId, dstFilename, dstFormat string) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	jsonData := map[string]interface{}{
		"srcStorage":  pool.ID,
		"srcFilename": srcFilename,
		"format":      "auto",
		"dstStorage":  dstStorageId,
		"dstFilename": dstFilename,
		"output":      dstFormat}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.Request("POST", "template/convert", jsonValue) //TODO: move to storage in rest api for 8.0
	if err != nil {
		return err
	}
	return err
}

func (pool *StoragePool) DeleteFile(client *Client, filename string) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	body, err := client.Request("DELETE", fmt.Sprintf("storage/pool/%s/%s", pool.ID, filename), nil)
	var res struct {
		Deleted bool `json:"deleted"`
	}
	err = json.Unmarshal(body, &res)
	if err == nil && !res.Deleted {
		err = (fmt.Errorf("Error: Unable to delete %s from %s", filename, pool.Name))
	}
	return err
}

func (pool *StoragePool) Browse(client *Client) ([]string, error) {
	var files []string
	if pool.ID == "" {
		return files, errors.New("Invalid Storage Pool")
	}

	body, err := client.Request("GET", fmt.Sprintf("storage/pool/%s/browse", pool.ID), nil)
	if err != nil {
		return files, err
	}
	err = json.Unmarshal(body, &files)
	return files, err
}
