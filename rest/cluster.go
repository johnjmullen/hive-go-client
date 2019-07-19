package rest

import (
	"encoding/json"
	"errors"
)

type Cluster struct {
	AdConfig struct {
		Domain    string      `json:"domain,omitempty"`
		Ou        interface{} `json:"ou,omitempty,omitempty"`
		Password  string      `json:"password,omitempty"`
		UserGroup string      `json:"userGroup,omitempty"`
		Username  string      `json:"username,omitempty"`
	} `json:"adConfig,omitempty"`
	BrokerOptions struct {
		AllowDesktopComposition bool `json:"allowDesktopComposition,omitempty"`
		AudioCapture            bool `json:"audioCapture,omitempty"`
		RedirectCSSP            bool `json:"redirectCSSP,omitempty"`
		RedirectClipboard       bool `json:"redirectClipboard,omitempty"`
		RedirectDisk            bool `json:"redirectDisk,omitempty"`
		RedirectPNP             bool `json:"redirectPNP,omitempty"`
		RedirectPrinter         bool `json:"redirectPrinter,omitempty"`
		RedirectUSB             bool `json:"redirectUSB,omitempty"`
		SmartResize             bool `json:"smartResize,omitempty"`
	} `json:"brokerOptions,omitempty"`
	BypassBroker bool     `json:"bypassBroker,omitempty"`
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	Tags         []string `json:"tags,omitempty"`
	Timezone     string   `json:"timezone,omitempty"`
	UserVolumes  struct {
		BackupSchedule int    `json:"backupSchedule,omitempty"`
		Repository     string `json:"repository,omitempty"`
		Size           int    `json:"size,omitempty"`
		Target         string `json:"target,omitempty"`
	} `json:"userVolumes,omitempty"`
	Vlan int `json:"vlan,omitempty"`
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
