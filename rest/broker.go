package rest

import (
	"encoding/json"

	"github.com/hashicorp/go-version"
)

//BrokerPool describes a pool received from BrokerLogin
type BrokerPool struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Os        string        `json:"os"`
	UserGroup string        `json:"userGroup"`
	Guests    []interface{} `json:"guests"`
}

type brokerLoginResponse struct {
	Token string       `json:"token"`
	Pools []BrokerPool `json:"pools"`
}

// BrokerLogin connects to the broker with the provided username, password, and realm
// returns a list of available pools for the user or an error
func (client *Client) BrokerLogin(username, password, realm string) ([]BrokerPool, error) {
	jsonData := map[string]string{"username": username, "password": password, "realm": realm}
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
	if err == nil {
		client.token = resp.Token
	}
	return resp.Pools, err
}

// BrokerAssignments returns an array of assignments for the logged in user
func (client *Client) BrokerAssignments() ([]interface{}, error) {
	body, err := client.request("GET", "broker/assignments", nil)
	var assignments []interface{}
	if err != nil {
		return assignments, err
	}
	err = json.Unmarshal(body, &assignments)
	return assignments, err
}

// BrokerAssign requests a desktop from a pool
func (client *Client) BrokerAssign(poolID string) (interface{}, error) {
	body, err := client.request("POST", "broker/assign/"+poolID, nil)
	var result interface{}
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
