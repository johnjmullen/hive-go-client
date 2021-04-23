package rest

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"
)

/*{"backupSchedule":null,
"lastReplication":"2021-01-27T02:04:05.099Z",
"source":"/mnt/9c7c2d0d-3b14-44d1-9a78-54bc7c055770/user1.HOME.qcow2",
"state":"ready",
"stateMessage":"Agent added user volume",
"target":"nocache"}*/
//BrokerGuest describes a guest assignment
type BrokerGuest struct {
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	HostID     string `json:"hostid"`
	GuestState string `json:"guestState"`
	Username   string `json:"username"`
	PoolID     string `json:"poolID"`
	IP         string `json:"ip"`
	UserVolume *struct {
		BackupSchedule  interface{} `json:",omitempty"`
		LastReplication interface{} `json:",omitempty"`
		Source          string      `json:"source,omitempty"`
		State           string      `json:"state,omitempty"`
		StateMessage    interface{} `json:"stateMessage,omitempty"`
		Target          string      `json:"Target,omitempty"`
		RunningBackup   bool        `json:"runningBackup,omitempty"`
	} `json:"userVolume,omitempty"`
}

//BrokerPool describes a pool received from BrokerLogin
type BrokerPool struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Os        string        `json:"os"`
	UserGroup string        `json:"userGroup"`
	Guests    []BrokerGuest `json:"guests"`
}

type brokerLoginResponse struct {
	Token string       `json:"token"`
	Pools []BrokerPool `json:"pools"`
}

type BrokerConfig struct {
	AutoConnectUserDesktop    bool   `json:"autoConnectUserDesktop"`
	BackgroundColor           string `json:"backgroundColor"`
	BgImage                   string `json:"bgImage"`
	BgImageFilename           string `json:"bgImageFilename"`
	ButtonTextColor           string `json:"buttonTextColor"`
	ClientSourceIsolation     bool   `json:"clientSourceIsolation"`
	Disclaimer                string `json:"disclaimer"`
	Enabled                   bool   `json:"enabled"`
	External                  bool   `json:"external"`
	ExternalProfile           string `json:"externalProfile"`
	Favicon                   string `json:"favicon"`
	FaviconFilename           string `json:"faviconFilename"`
	HideRealms                bool   `json:"hideRealms"`
	HideRelease               bool   `json:"hideRelease"`
	Logo                      string `json:"logo"`
	LogoFilename              string `json:"logoFilename"`
	MainColor                 string `json:"mainColor"`
	PassthroughAuthentication bool   `json:"passthroughAuthentication"`
	TextColor                 string `json:"textColor"`
	Title                     string `json:"title"`
	TwoFormAuth               struct {
		Enabled      bool   `json:"enabled"`
		EnforceLocal bool   `json:"enforceLocal"`
		Type         string `json:"type"`
	} `json:"twoFormAuth"`
	Remote bool `json:"remote"`
}

// GetBrokerConfig returns the configuration options of the broker
func (client *Client) GetBrokerConfig() (BrokerConfig, error) {
	var config BrokerConfig
	body, err := client.request("GET", "broker/configuration", nil)
	if err != nil {
		return config, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &config)
	if err != nil {
		return config, err
	}
	return config, err
}

// BrokerLogin connects to the broker with the provided username, password, and realm
// returns a list of available pools for the user or an error
func (client *Client) BrokerLogin(username, password, realm, token string) ([]BrokerPool, error) {
	jsonData := map[string]string{"username": username, "password": password, "realm": realm}
	if token != "" {
		jsonData["token"] = token
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	body, err := client.request("POST", "authBrokerUser", jsonValue)
	if err != nil {
		return nil, err
	}

	var resp brokerLoginResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	client.token = resp.Token
	return resp.Pools, nil
}

type brokerAssignments struct {
	Pools []BrokerPool `json:"pools"`
}

// BrokerAssignments returns an array of assignments for the logged in user
func (client *Client) BrokerAssignments() ([]BrokerPool, error) {
	body, err := client.request("GET", "broker/assignments", nil)
	if err != nil {
		return nil, err
	}
	var assignments brokerAssignments
	fmt.Println(string(body))
	err = json.Unmarshal(body, &assignments)
	if err != nil {
		return nil, err
	}
	return assignments.Pools, err
}

// BrokerAssign requests a desktop from a pool
func (client *Client) BrokerAssign(poolID string) (BrokerGuest, error) {
	body, err := client.request("POST", "broker/assign/"+poolID, nil)
	var result BrokerGuest
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// BrokerConnect request the rdp file to connect to a guest.
// outputType can be rdp, json, or hio
func (client *Client) BrokerConnect(guest string, remote bool, outputType string) ([]byte, error) {
	jsonData := map[string]interface{}{"guest": guest, "remote": remote, "outputType": outputType}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	body, err := client.request("POST", "broker/"+guest+"/connect", jsonValue)
	return body, err
}

// AssignGuest assign a user to a specific guest
func (client *Client) AssignGuest(poolID, username, realm, guest string) (interface{}, error) {
	var result interface{}
	jsonData := map[string]string{"realm": realm, "username": username}
	hostVersion, err := client.HostVersion()
	if err != nil {
		return nil, err
	}
	minVersion, _ := version.NewVersion("8.3.0")
	v, err := version.NewVersion(hostVersion.Version)
	if err != nil || v.LessThan(minVersion) {
		jsonData["guest"] = guest
	} else {
		jsonData["guestName"] = guest
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return result, err
	}
	body, err := client.request("POST", "broker/assign/"+poolID, jsonValue)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// ReleaseGuest release a guest that is currently assigned
func (client *Client) ReleaseGuest(poolID, username, guest string) error {
	jsonData := map[string]string{"poolId": poolID, "username": username, "guest": guest}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	body, err := client.request("POST", "broker/release", jsonValue)
	var result interface{}
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	return err
}
