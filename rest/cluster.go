package rest

import (
	"encoding/json"
	"errors"
	"time"
)

// Broker settings from the cluster table
type Broker struct {
	AutoConnectUserDesktop    bool        `json:"autoConnectUserDesktop,omitempty"`
	BackgroundColor           string      `json:"backgroundColor,omitempty"`
	BgImage                   string      `json:"bgImage,omitempty"`
	BgImageFilename           string      `json:"bgImageFilename,omitempty"`
	ButtonTextColor           string      `json:"buttonTextColor,omitempty"`
	Disclaimer                string      `json:"disclaimer,omitempty"`
	Enabled                   bool        `json:"enabled"`
	External                  bool        `json:"external"`
	ExternalProfile           string      `json:"externalProfile,omitempty"`
	Favicon                   string      `json:"favicon,omitempty"`
	FaviconFilename           string      `json:"faviconFilename,omitempty"`
	HideRealms                bool        `json:"hideRealms,omitempty"`
	HideRelease               bool        `json:"hideRelease,omitempty"`
	Logo                      string      `json:"logo,omitempty"`
	LogoFilename              string      `json:"logoFilename,omitempty"`
	MainColor                 string      `json:"mainColor,omitempty"`
	PassthroughAuthentication bool        `json:"passthroughAuthentication,omitempty"`
	TextColor                 string      `json:"textColor,omitempty"`
	Title                     string      `json:"title,omitempty"`
	TwoFormAuth               interface{} `json:"twoFormAuth,omitempty"`
	AllowPhysical             bool        `json:"allowPhysical,omitempty"`
}

type GatewayHost struct {
	StartPort       uint   `json:"startPort"`
	EndPort         uint   `json:"endPort"`
	ExternalAddress string `json:"externalAddress"`
}

// Settings for email alerts
type EmailAlerts struct {
	Service   string                   `json:"service,omitempty"`
	Host      string                   `json:"host,omitempty"`
	SMTPPort  uint                     `json:"smtpPort,omitempty"`
	Secure    bool                     `json:"secure"`
	Username  string                   `json:"username"`
	Password  string                   `json:"password"`
	Level     string                   `json:"level,omitempty"`
	From      string                   `json:"from,omitempty"`
	To        string                   `json:"to"`
	BlackList []map[string]interface{} `json:"blackList,omitempty"`
}

// Gateway settings from the cluster table
type Gateway struct {
	Enabled               bool                   `json:"enabled"`
	ClientSourceIsolation bool                   `json:"clientSourceIsolation"`
	PersistentAssignment  bool                   `json:"persistentAssignment"`
	Hosts                 map[string]GatewayHost `json:"hosts"`
}

// ClusterBackup data protection settings from the cluster table
type ClusterBackup struct {
	Enabled     bool   `json:"enabled"`
	StartWindow string `json:"startWindow"`
	EndWindow   string `json:"endWindow"`
}

type ClusterSSO struct {
	Enabled      bool   `json:"enabled"`
	Provider     string `json:"provider"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret,omitempty"`
	DiscoveryURL string `json:"discoveryUrl,omitempty"`
	TenantID     string `json:"tenantId,omitempty"`
}

// Cluster record from the rest api
type Cluster struct {
	AdminPassword string   `json:"adminPassword"`
	Broker        *Broker  `json:"broker"`
	Gateway       *Gateway `json:"gateway"`
	HiveSense     struct {
		AwsAccessKeyID     string `json:"awsAccessKeyId,omitempty"`
		AwsSecretAccessKey string `json:"awsSecretAccessKey,omitempty"`
		CustomerName       string `json:"customerName,omitempty"`
		Enabled            bool   `json:"enabled"`
		LogStatus          string `json:"logStatus,omitempty"`
		UploadFrequency    int    `json:"uploadFrequency,omitempty"`
	} `json:"hiveSense"`
	ID      string `json:"id"`
	License *struct {
		Expiration time.Time `json:"expiration"`
		Type       string    `json:"type"`
		MaxGuests  int       `json:"maxGuests,omitempty"`
	} `json:"license"`
	Name          string `json:"name"`
	SharedStorage *struct {
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
	Backup      *ClusterBackup `json:"backup,omitempty"`
	Tags        []string       `json:"tags"`
	EmailAlerts EmailAlerts    `json:"emailAlerts,omitempty"`
	SSOProvider *ClusterSSO    `json:"ssoProvider,omitempty"`
}

func (cluster Cluster) String() string {
	json, _ := json.MarshalIndent(cluster, "", "  ")
	return string(json)
}

// ListClusters request a list of all clusters
func (client *Client) ListClusters(query string) ([]Cluster, error) {
	var clusters []Cluster
	path := "clusters"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return clusters, err
	}
	err = json.Unmarshal(body, &clusters)
	return clusters, err
}

// GetCluster request a cluster by id
func (client *Client) GetCluster(id string) (Cluster, error) {
	var cluster Cluster
	if id == "" {
		return cluster, errors.New("id cannot be empty")
	}
	body, err := client.request("GET", "cluster/"+id, nil)
	if err != nil {
		return cluster, err
	}
	err = json.Unmarshal(body, &cluster)
	return cluster, err
}

// JoinHost Add a new host to the cluster
func (client *Client) JoinHost(username, password, ipAddress string) (*Task, error) {
	jsonData := map[string]string{"remoteUsername": username, "remotePassword": password, "remoteIpAddress": ipAddress}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	//Why doesn't this use clusterId?
	return client.getTaskFromResponse(client.request("POST", "cluster/joinHost", jsonValue))
}

// GetLicenseInfo Lookup license information for the current cluster
func (cluster *Cluster) GetLicenseInfo(client *Client) (string, string, error) {
	body, err := client.request("GET", "cluster/"+cluster.ID+"/license", nil)
	if err != nil {
		return "", "", err
	}
	var objMap map[string]string
	err = json.Unmarshal(body, &objMap)

	return objMap["expiration"], objMap["type"], err
}

// SetLicense replaces the license for the cluster
func (cluster *Cluster) SetLicense(client *Client, key string) error {
	jsonData := map[string]string{"key": key}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.request("PUT", "cluster/"+cluster.ID+"/license", jsonValue)
	return err
}

// EnableBackup enable automatic data protection
// startWindow and endWindow must be strings in the format "01:00:00"
func (cluster *Cluster) EnableBackup(client *Client, startWindow string, endWindow string) error {
	jsonData := map[string]string{"startWindow": startWindow, "endWindow": endWindow}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	_, err = client.request("POST", "cluster/"+cluster.ID+"/enableBackup", jsonValue)
	return err
}

// DisableBackup disable automatic data protection
func (cluster *Cluster) DisableBackup(client *Client) error {
	_, err := client.request("POST", "cluster/"+cluster.ID+"/disableBackup", nil)
	return err
}

// EnableSharedStorage enable shared storage on a cluster
// storageUtilization is a percentage for the amount of storage to be allocated for shared storage
// minSetSize is the number of host to grow the storage by (2 or 3)
func (cluster *Cluster) EnableSharedStorage(client *Client, storageUtilization int, minSetSize int) (*Task, error) {
	jsonData := map[string]int{"minSetSize": minSetSize, "storageUtilization": storageUtilization}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "cluster/"+cluster.ID+"/enableSharedStorage", jsonValue))
}

// DisableSharedStorage disables shared storage on the cluster
func (cluster *Cluster) DisableSharedStorage(client *Client) (*Task, error) {
	return client.getTaskFromResponse(client.request("POST", "cluster/"+cluster.ID+"/disableSharedStorage", nil))
}

// GetBroker returns the broker settings for the cluster
func (client *Client) GetBroker(clusterID string) (Broker, error) {
	var broker Broker
	body, err := client.request("GET", "cluster/"+clusterID+"/broker", nil)
	if err != nil {
		return broker, err
	}
	err = json.Unmarshal(body, &broker)

	return broker, err
}

// SetBroker updates broker settings for the cluster
func (client *Client) SetBroker(clusterID string, brokerSettings Broker) error {
	jsonValue, err := json.Marshal(brokerSettings)
	if err != nil {
		return err
	}
	_, err = client.request("PUT", "cluster/"+clusterID+"/broker", jsonValue)
	return err
}

// ResetBroker updates broker settings for the cluster
func (client *Client) ResetBroker(clusterID string) error {
	_, err := client.request("POST", "cluster/"+clusterID+"/broker/reset", nil)
	return err
}

// GetGateway returns the gateway settings for the cluster
func (client *Client) GetGateway(clusterID string) (Gateway, error) {
	var gateway Gateway
	body, err := client.request("GET", "cluster/"+clusterID+"/gateway", nil)
	if err != nil {
		return gateway, err
	}
	err = json.Unmarshal(body, &gateway)

	return gateway, err
}

// SetGateway updates gateway settings for the cluster
func (client *Client) SetGateway(clusterID string, gatewaySettings Gateway) error {
	jsonValue, err := json.Marshal(gatewaySettings)
	if err != nil {
		return err
	}
	_, err = client.request("PUT", "cluster/"+clusterID+"/gateway", jsonValue)
	return err
}

// UpdateSoftware applies a software package across the cluster
func (cluster Cluster) UpdateSoftware(client *Client, packageName string) (*Task, error) {
	jsonData := map[string]string{"packageName": packageName}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return client.getTaskFromResponse(client.request("POST", "cluster/"+cluster.ID+"/updatePackage", jsonValue))
}

// EmailAlerts updates email alert settings
func (cluster Cluster) SetEmailAlerts(client *Client, emailAlerts EmailAlerts) error {
	jsonValue, err := json.Marshal(emailAlerts)
	if err != nil {
		return err
	}
	_, err = client.request("POST", "cluster/"+cluster.ID+"/emailAlerts", jsonValue)
	return err
}

// ClearEmailAlerts updates email alert settings
func (cluster Cluster) ClearEmailAlerts(client *Client) error {
	_, err := client.request("POST", "cluster/"+cluster.ID+"/clearEmailAlerts", nil)
	return err
}

// SendTestEmail sends an email to test the email alert settings
func (cluster Cluster) SendTestEmail(client *Client) error {
	_, err := client.request("POST", "cluster/"+cluster.ID+"/sendTestEmail", nil)
	return err
}

func (cluster Cluster) EnableSSO(client *Client, settings ClusterSSO) error {
	jsonValue, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	_, err = client.request("POST", "cluster/"+cluster.ID+"/enableSSO", jsonValue)
	return err
}

func (cluster Cluster) DisableSSO(client *Client) error {
	_, err := client.request("POST", "cluster/"+cluster.ID+"/disableSSO", nil)
	return err
}

// SSOInfo returns the SSO information for the cluster
func (cluster Cluster) SSOInfo(client *Client) (ClusterSSO, error) {
	var sso ClusterSSO
	body, err := client.request("GET", "cluster/"+cluster.ID+"/ssoInfo", nil)
	if err != nil {
		return sso, err
	}
	err = json.Unmarshal(body, &sso)
	return sso, err
}
