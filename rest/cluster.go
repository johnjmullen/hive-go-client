package rest

import (
	"encoding/json"
	"errors"
	"time"
)

type Cluster struct {
	AdminPassword string `json:"adminPassword"`
	Gateway       struct {
		Enabled bool `json:"enabled"`
	} `json:"gateway"`
	HiveSense struct {
		AwsAccessKeyID     string `json:"awsAccessKeyId"`
		AwsSecretAccessKey string `json:"awsSecretAccessKey"`
		CustomerName       string `json:"customerName"`
		Enabled            bool   `json:"enabled"`
		UploadFrequency    int    `json:"uploadFrequency"`
	} `json:"hiveSense"`
	ID      string `json:"id"`
	License struct {
		Expiration        time.Time `json:"expiration"`
		Type              string    `json:"type"`
		LicenseExpiration time.Time `json:"license.expiration"`
	} `json:"license"`
	Name          string `json:"name"`
	SharedStorage struct {
		State string `json:"state"`
	} `json:"sharedStorage"`
	Tags []string `json:"tags"`
}

func (cluster Cluster) String() string {
	json, _ := json.MarshalIndent(cluster, "", "  ")
	return string(json)
}

func (client *Client) ListClusters() ([]Cluster, error) {
	var clusters []Cluster
	body, err := client.Request("GET", "clusters", nil)
	if err != nil {
		return clusters, err
	}
	err = json.Unmarshal(body, &clusters)
	return clusters, err
}

func (client *Client) GetCluster(id string) (Cluster, error) {
	var cluster Cluster
	if id == "" {
		return cluster, errors.New("Id cannot be empty")
	}
	body, err := client.Request("GET", "cluster/"+id, nil)
	if err != nil {
		return cluster, err
	}
	err = json.Unmarshal(body, &cluster)
	return cluster, err
}

func (client *Client) JoinHost(username, password, ipAddress string) error {
	jsonData := map[string]string{"remoteUsername": username, "remotePassword": password, "remoteIpAddress": ipAddress}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	//Why doesn't this use clusterId?
	_, err = client.Request("POST", "cluster/joinHost", jsonValue)
	//TODO: Need to watch task
	return err
}

func (cluster *Cluster) GetLicenseInfo(client *Client) (string, string, error) {
	body, err := client.Request("GET", "cluster/"+cluster.ID+"/license", nil)
	if err != nil {
		return "", "", err
	}
	var objMap map[string]string
	err = json.Unmarshal(body, &objMap)

	return objMap["expiration"], objMap["type"], err
}

func (cluster *Cluster) SetLicense(client *Client, key string) error {
	jsonData := map[string]string{"key": key}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.Request("PUT", "cluster/"+cluster.ID+"/license", jsonValue)
	return err
}
