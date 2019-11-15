package rest

import (
	"encoding/json"
	"errors"
	"fmt"
)

type StoragePool struct {
	ID                string   `json:"id,omitempty"`
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	Server            string   `json:"server,omitempty"`
	Path              string   `json:"path,omitempty"`
	URL               string   `json:"url,omitempty"`
	Username          string   `json:"username,omitempty"`
	Password          string   `json:"password,omitempty"`
	Key               string   `json:"key,omitempty"`
	MountOptions      []string `json:"mountOptions,omitempty"`
	Roles             []string `json:"roles,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	S3AccessKeyID     string   `json:"s3AccessKeyId,omitempty"`
	S3SecretAccessKey string   `json:"s3SecretAccessKey,omitempty"`
	S3Region          string   `json:"s3Region,omitempty"`
}

type DiskInfo struct {
	Filename            string   `json:"filename,omitempty"`
	VirtualSize         uint     `json:"virtual-size,omitempty"`
	ActualSize          uint     `json:"actual-size,omitempty"`
	DirtyFlag           bool     `json:"dirty-flag,omitempty"`
	ClusterSize         uint     `json:"cluster-size,omitempty"`
	Encrypted           bool     `json:"encrypted,omitempty"`
	BackingFilename     string   `json:"backing-filename,omitempty"`
	FullBackingFilename string   `json:"full-backing-filename,omitempty"`
	Snapshots           []string `json:"snapshots,omitempty"`
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
	body, err := client.request("GET", path, nil)
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
	pool := &StoragePool{}
	if id == "" {
		return pool, errors.New("id cannot be empty")
	}
	body, err := client.request("GET", "storage/pool/"+id, nil)
	if err != nil {
		return pool, err
	}
	err = json.Unmarshal(body, pool)
	return pool, err
}

func (pool *StoragePool) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("POST", "storage/pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (pool *StoragePool) Delete(client *Client) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	_, err := client.request("DELETE", "storage/pool/"+pool.ID, nil)
	return err
}

func (pool *StoragePool) CreateDisk(client *Client, filename, format string, size uint) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("Invalid Storage Pool")
	}
	jsonData := map[string]interface{}{"filename": filename, "size": size, "format": format}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "storage/pool/"+pool.ID+"/createDisk", jsonValue))
}

func (pool *StoragePool) ConvertDisk(client *Client, srcFilename, dstStorageId, dstFilename, dstFormat string) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("Invalid Storage Pool")
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
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "template/convert", jsonValue))
}

func (pool *StoragePool) CopyUrl(client *Client, url, filePath string) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("Invalid Storage Pool")
	}
	jsonData := map[string]interface{}{
		"url":      url,
		"filePath": filePath}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "storage/pool/"+pool.ID+"/copyUrl", jsonValue))
}

func (pool *StoragePool) DiskInfo(client *Client, filePath string) (DiskInfo, error) {
	var disk DiskInfo
	if pool.ID == "" {
		return disk, errors.New("Invalid Storage Pool")
	}
	jsonData := map[string]interface{}{"filePath": filePath}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return disk, err
	}
	body, err := client.request("POST", "storage/pool/"+pool.ID+"/diskInfo", jsonValue)
	if err != nil {
		return disk, err
	}
	err = json.Unmarshal(body, &disk)
	return disk, err
}

func (pool *StoragePool) GrowDisk(client *Client, filePath string, size uint) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	jsonData := map[string]interface{}{"filePath": filePath, "size": size}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.request("POST", "storage/pool/"+pool.ID+"/growDisk", jsonValue)
	return err
}

func (pool *StoragePool) DeleteFile(client *Client, filename string) error {
	if pool.ID == "" {
		return errors.New("Invalid Storage Pool")
	}
	body, err := client.request("DELETE", fmt.Sprintf("storage/pool/%s/%s", pool.ID, filename), nil)
	if err != nil {
		return err
	}
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

	body, err := client.request("GET", fmt.Sprintf("storage/pool/%s/browse", pool.ID), nil)
	if err != nil {
		return files, err
	}
	err = json.Unmarshal(body, &files)
	return files, err
}

func (client *Client) CopyFile(srcStorageID, srcFilePath, destStorageID, destFilePath string) (*Task, error) {
	data := map[string]string{
		"srcStorageId":  srcStorageID,
		"srcFilePath":   srcFilePath,
		"destStorageId": destStorageID,
		"destFilePath":  destFilePath,
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "storage/pool/copyFile", jsonValue))
}

func (client *Client) MoveFile(srcStorageID, srcFilePath, destStorageID, destFilePath string) (*Task, error) {
	data := map[string]string{
		"srcStorageId":  srcStorageID,
		"srcFilePath":   srcFilePath,
		"destStorageId": destStorageID,
		"destFilePath":  destFilePath,
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "storage/pool/moveFile", jsonValue))
}
