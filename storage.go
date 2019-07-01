package main

import (
    "encoding/json"
    "io/ioutil"
    "errors"
    "fmt"
)

type StoragePool struct {
    Id string             `json:"id,omitempty"`
    Name string           `json:"name"`
    Type string           `json:"type"`
    Server string         `json:"server"`
    Path string           `json:"path"`
    Username string       `json:"username,omitempty"`
    Password string       `json:"password,omitempty"`
    Key string            `json:"key,omitempty"`
    MountOptions []string `json:"mountOptions,omitempty"`
    Roles []string        `json:"roles,omitempty"`
    Tags []string         `json:"tags,omitempty"`
}

func (sp StoragePool) String() string {
     return fmt.Sprintf("{\n Name: %s,\n Id: %s,\n Type: %s,\n Server: %s,\n Path: %s\n}\n", sp.Name, sp.Id, sp.Type, sp.Server, sp.Path)
}

func (client *Client) ListStoragePools() ([]StoragePool, error) {
    var pools []StoragePool
    resp, err := client.Request("GET", "storage/pools", nil)
    if err != nil {
        return pools, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return pools, err
    }
    err = json.Unmarshal(body, &pools)
    return pools, err
}

func (client *Client) GetStoragePool(id string) (StoragePool, error) {
    var pool StoragePool
    if id == "" {
        return pool, errors.New("id cannot be empty")
    }
    resp, err := client.Request("GET", "storage/pool/"+id, nil)
    if err != nil {
        return pool, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return pool, err
    }
    err = json.Unmarshal(body, &pool)
    return pool, err
}

func (client *Client) CreateStoragePool(pool *StoragePool) (string, error) {
    var result string
    jsonValue, _ := json.Marshal(pool)
    resp, err := client.Request("POST", "storage/pools", jsonValue)
    if err != nil {
        return result, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err == nil {
        result = string(body)
    }
    return result, err
}

func (client *Client) DeleteStoragePool(id string) (error) {
    if id == "" {
        return errors.New("id cannot be empty")
    }
    resp, err := client.Request("DELETE", "storage/pool/"+id, nil)
    if err != nil {
        return err
    }
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    return err
}
