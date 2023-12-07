package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/eventials/go-tus"
)

// StoragePool describes a storage pool returned from the rest api
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
	S3Provider        string   `json:"s3Provider,omitempty"`
	S3Endpoint        string   `json:"s3Endpoint,omitempty"`
	Disabled          bool     `json:"disabled"`
	Replicated        bool     `json:"replicated,omitempty"`
}

// DiskInfo contains information about a disk from a storage pool
type DiskInfo struct {
	Filename            string   `json:"filename,omitempty"`
	Format              string   `json:"format,omitempty"`
	VirtualSize         uint     `json:"virtual-size,omitempty"`
	ActualSize          uint     `json:"actual-size,omitempty"`
	DirtyFlag           bool     `json:"dirty-flag,omitempty"`
	ClusterSize         uint     `json:"cluster-size,omitempty"`
	Encrypted           bool     `json:"encrypted,omitempty"`
	BackingFilename     string   `json:"backing-filename,omitempty"`
	FullBackingFilename string   `json:"full-backing-filename,omitempty"`
	Snapshots           []string `json:"snapshots,omitempty"`
}

func (pool StoragePool) String() string {
	json, _ := json.MarshalIndent(pool, "", "  ")
	return string(json)
}

// ListStoragePools returns an array of storage pools with an optional filter string
func (client *Client) ListStoragePools(query string) ([]StoragePool, error) {
	var pools []StoragePool
	path := "storage/pools"
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

// GetStoragePoolByName requests a storage pool by name
func (client *Client) GetStoragePoolByName(name string) (*StoragePool, error) {
	if name == "disk" {
		return &StoragePool{
			ID:   "disk",
			Name: "disk",
			Type: "disk",
			Path: "/zdata",
		}, nil
	}
	if name == "ram" {
		return &StoragePool{
			ID:   "ram",
			Name: "ram",
			Type: "ram",
			Path: "/zram",
		}, nil
	}
	var pools, err = client.ListStoragePools("name=" + url.PathEscape(name))
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		if pool.Name == name {
			return &pool, nil
		}
	}
	return nil, errors.New("storage Pool not found")
}

// GetStoragePool requests a storage pool by id
func (client *Client) GetStoragePool(id string) (*StoragePool, error) {
	if id == "disk" {
		return &StoragePool{
			ID:   "disk",
			Name: "disk",
			Type: "disk",
			Path: "/zdata",
		}, nil
	}
	if id == "ram" {
		return &StoragePool{
			ID:   "ram",
			Name: "ram",
			Type: "ram",
			Path: "/zram",
		}, nil
	}
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

// Create creates a new storage pool
func (pool *StoragePool) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(pool)
	body, err := client.request("POST", "storage/pools", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Delete deletes a storage pool
func (pool *StoragePool) Delete(client *Client) error {
	if pool.ID == "" {
		return errors.New("invalid Storage Pool")
	}
	_, err := client.request("DELETE", "storage/pool/"+pool.ID, nil)
	return err
}

// Stop disables a storage pool
func (pool *StoragePool) Stop(client *Client) error {
	if pool.ID == "" {
		return errors.New("invalid Storage Pool")
	}
	_, err := client.request("POST", "storage/pool/"+pool.ID+"/stop", nil)
	return err
}

// Start re-enables a stopped storage pool
func (pool *StoragePool) Start(client *Client) error {
	if pool.ID == "" {
		return errors.New("invalid Storage Pool")
	}
	_, err := client.request("POST", "storage/pool/"+pool.ID+"/start", nil)
	return err
}

// CreateDisk creates a new disk in the storage pool
func (pool *StoragePool) CreateDisk(client *Client, filename, format string, size uint) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("invalid Storage Pool")
	}
	jsonData := map[string]interface{}{"filename": filename, "size": size, "format": format}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "storage/pool/"+pool.ID+"/createDisk", jsonValue))
}

// ConvertDisk converts or copies a disk to a new file
func (pool *StoragePool) ConvertDisk(client *Client, srcFilename, dstStorageID, dstFilename, dstFormat string) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("invalid Storage Pool")
	}
	jsonData := map[string]interface{}{
		"srcStorage":  pool.ID,
		"srcFilename": srcFilename,
		"format":      "auto",
		"dstStorage":  dstStorageID,
		"dstFilename": dstFilename,
		"output":      dstFormat}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "template/convert", jsonValue))
}

// CopyURL downloads a file from a http url into a storage pool
func (pool *StoragePool) CopyURL(client *Client, url, filePath string) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("invalid Storage Pool")
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

// DiskInfo retrieves information about a disk in the storage pool
func (pool *StoragePool) DiskInfo(client *Client, filePath string) (DiskInfo, error) {
	var disk DiskInfo
	if pool.ID == "" {
		return disk, errors.New("invalid Storage Pool")
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

// GrowDisk increases the size of a disk in a storage pool by size GB
func (pool *StoragePool) GrowDisk(client *Client, filePath string, size uint) (*Task, error) {
	if pool.ID == "" {
		return nil, errors.New("invalid Storage Pool")
	}
	jsonData := map[string]interface{}{"filePath": filePath, "size": size}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "storage/pool/"+pool.ID+"/growDisk", jsonValue))
}

// DeleteFile deletes a file from a storage pool
func (pool *StoragePool) DeleteFile(client *Client, filename string) error {
	if pool.ID == "" {
		return errors.New("invalid Storage Pool")
	}
	body, err := client.request("DELETE", fmt.Sprintf("storage/pool/%s/%s", pool.ID, url.PathEscape(filename)), nil)
	if err != nil {
		return err
	}
	var res struct {
		Deleted bool `json:"deleted"`
	}
	err = json.Unmarshal(body, &res)
	if err == nil && !res.Deleted {
		err = (fmt.Errorf("error: Unable to delete %s from %s", filename, pool.Name))
	}
	return err
}

// StoragePoolFileInfo contains information about the files returned by Browse
type StoragePoolFileInfo struct {
	Path     string `json:"Path"`
	Name     string `json:"Name"`
	Size     int    `json:"Size"`
	MimeType string `json:"MimeType"`
	ModTime  string `json:"ModTime"`
	IsDir    bool   `json:"IsDir"`
}

// Browse returns a list of files from a storage pool
func (pool *StoragePool) Browse(client *Client, filePath string, recursive bool) ([]StoragePoolFileInfo, error) {
	var files []StoragePoolFileInfo
	if pool.ID == "" {
		return files, errors.New("invalid Storage Pool")
	}
	options := "details=true"
	if len(filePath) > 0 {
		options += "&filePath=" + url.QueryEscape(filePath)
	}
	if recursive {
		options += "&recursive=true"
	}
	body, err := client.request("GET", fmt.Sprintf("storage/pool/%s/browse?%s", pool.ID, options), nil)
	if err != nil {
		return files, err
	}
	err = json.Unmarshal(body, &files)
	return files, err
}

// Download downloads a file from a storage pool
func (pool *StoragePool) Download(client *Client, filePath string) (*http.Response, error) {
	return pool.DownloadWithContext(context.Background(), client, filePath)
}

// DownloadWithContext downloads a file from a storage pool with a custom context
func (pool *StoragePool) DownloadWithContext(ctx context.Context, client *Client, filePath string) (*http.Response, error) {
	if pool.ID == "" {
		return nil, errors.New("invalid Storage Pool")
	}
	if err := client.CheckHostVersion("8.5.0"); err != nil {
		return nil, err
	}
	headers := map[string]string{"Content-type": "application/json"}
	resp, err := client.requestWithHeaders(ctx, "GET", fmt.Sprintf("storage/pool/%s/download?filePath=%s", pool.ID, url.QueryEscape(filePath)), bytes.NewBuffer(nil), headers, 0)
	return resp, err
}

// Upload uploads a local file into a storage pool
func (pool *StoragePool) Upload(client *Client, filename, targetFilename string) error {
	if pool.ID == "" {
		return errors.New("invalid Storage Pool")
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	header := make(http.Header)
	header.Add("Authorization", "Bearer "+client.token)
	conf := tus.Config{
		ChunkSize:           2 * 1024 * 1024,
		Resume:              false,
		OverridePatchMethod: false,
		Store:               nil,
		Header:              header,
		HttpClient:          client.httpClient,
	}

	uploadURL := fmt.Sprintf("https://%s:%d/upload/", client.Host, client.Port)
	tusClient, err := tus.NewClient(uploadURL, &conf)
	if err != nil {
		return err
	}

	upload, err := tus.NewUploadFromFile(f)
	upload.Metadata["storageId"] = pool.ID
	upload.Metadata["filename"] = targetFilename
	if err != nil {
		return err
	}

	uploader, err := tusClient.CreateUpload(upload)
	if err != nil {
		return err
	}

	err = uploader.Upload()
	return err
}

// CopyFile copies a file in a storage pool to a new file in another storage pool
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

// MoveFile moves or renames a file in a storage pool
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
