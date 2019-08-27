package rest

import (
	"encoding/json"
	"errors"
	"time"
)

//Broker information from the cluster table
type Broker struct {
	AutoConnectUserDesktop    bool   `json:"autoConnectUserDesktop,omitempty"`
	BackgroundColor           string `json:"backgroundColor,omitempty"`
	BgImage                   string `json:"bgImage,omitempty"`
	BgImageFilename           string `json:"bgImageFilename,omitempty"`
	ButtonTextColor           string `json:"buttonTextColor,omitempty"`
	Disclaimer                string `json:"disclaimer,omitempty"`
	Enabled                   bool   `json:"enabled"`
	Favicon                   string `json:"favicon,omitempty"`
	FaviconFilename           string `json:"faviconFilename,omitempty"`
	HideRealms                bool   `json:"hideRealms,omitempty"`
	HideRelease               bool   `json:"hideRelease,omitempty"`
	Logo                      string `json:"logo,omitempty"`
	LogoFilename              string `json:"logoFilename,omitempty"`
	MainColor                 string `json:"mainColor,omitempty"`
	PassthroughAuthentication bool   `json:"passthroughAuthentication,omitempty"`
	TextColor                 string `json:"textColor,omitempty"`
	Title                     string `json:"title,omitempty"`
	TwoFormAuth               struct {
	} `json:"twoFormAuth"`
}

//Gateway settings from the cluster table
type Gateway struct {
	Enabled bool `json:"enabled"`
	PortMap struct {
		F91577A2F6F8 struct {
			EndPort   int `json:"endPort"`
			StartPort int `json:"startPort"`
		} `json:"f91577a2f6f8"`
	} `json:"portMap"`
	URI string `json:"uri"`
}

//Cluster record from hive-rest
type Cluster struct {
	AdminPassword string   `json:"adminPassword"`
	Broker        *Broker  `json:"broker"`
	Gateway       *Gateway `json:"gateway"`
	HiveSense     struct {
		AwsAccessKeyID     string `json:"awsAccessKeyId"`
		AwsSecretAccessKey string `json:"awsSecretAccessKey"`
		CustomerName       string `json:"customerName"`
		Enabled            bool   `json:"enabled"`
		LogStatus          string `json:"logStatus"`
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
		Enabled bool `json:"enabled"`
		Hosts   []struct {
			Hostid string `json:"hostid"`
			State  string `json:"state"`
		} `json:"hosts"`
		ID                 string `json:"id"`
		MinSetSize         int    `json:"minSetSize"`
		State              string `json:"state"`
		StorageUtilization int    `json:"storageUtilization"`
	} `json:"sharedStorage"`
	Tags []string `json:"tags"`
}

func (cluster Cluster) String() string {
	json, _ := json.MarshalIndent(cluster, "", "  ")
	return string(json)
}

func (client *Client) ListClusters(filter string) ([]Cluster, error) {
	var clusters []Cluster
	path := "clusters"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.Request("GET", path, nil)
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
