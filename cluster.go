package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Vlan   int `json:"vlan,omitempty"`
	client *Client
}

func (cluster *Cluster) String() string {
	json, _ := json.MarshalIndent(cluster, "", "  ")
	return string(json)
}

func (cluster *Cluster) ToJson() ([]byte, error) {
	return json.Marshal(cluster)
}

func (cluster *Cluster) FromJson(data []byte) error {
	return json.Unmarshal(data, cluster)
}

func (client *Client) ListClusters() ([]Cluster, error) {
	var clusters []Cluster
	body, err := client.Request("GET", "clusters", nil)
	if err != nil {
		return clusters, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &clusters)
	for _, cluster := range clusters {
		cluster.client = client
	}
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
	cluster.client = client
	return cluster, err
}

func (cluster *Cluster) JoinHost(username, password, ipAddress string) error {
	jsonData := map[string]string{"remoteUsername": username, "remotePassword": password, "remoteIpAddress": ipAddress}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = cluster.client.Request("POST", fmt.Sprintf("cluster/%s/joinHost", cluster.ID), jsonValue)
	//TODO: Need to watch task
	return err
}
