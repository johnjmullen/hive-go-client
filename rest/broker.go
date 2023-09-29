package rest

import (
	"encoding/json"

	"github.com/hashicorp/go-version"
)

// BrokerGuest describes a guest assignment
type BrokerGuest struct {
	Name        string                  `json:"name"`
	UUID        string                  `json:"uuid"`
	HostID      string                  `json:"hostid"`
	GuestState  string                  `json:"guestState"`
	Username    string                  `json:"username"`
	PoolID      string                  `json:"poolID"`
	IP          string                  `json:"ip"`
	Connections []GuestBrokerConnection `json:"connections"`
	UserVolume  *UserVolume             `json:"userVolume,omitempty"`
}

// BrokerPool describes a pool received from BrokerLogin
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
	err = json.Unmarshal(body, &config)
	if err != nil {
		return config, err
	}
	return config, err
}

// BrokerLogin connects to the broker with the provided username, password, and realm
// returns nil or an error
func (client *Client) BrokerLogin(username, password, realm, token, mfaToken string) error {
	jsonData := map[string]string{"username": username, "password": password, "realm": realm}
	if token != "" {
		jsonData["token"] = token
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	body, err := client.request("POST", "authBrokerUser", jsonValue)
	if err != nil {
		return err
	}

	var resp brokerLoginResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	client.token = resp.Token
	return nil
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
func (client *Client) BrokerConnect(guest string, outputType string, connection string) ([]byte, error) {
	jsonData := map[string]interface{}{"guest": guest, "outputType": outputType, "connection": connection}
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
