package rest

import (
	"encoding/json"
)

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

func (client *Client) BrokerAssignments() ([]interface{}, error) {
	body, err := client.request("GET", "broker/assignments", nil)
	var assignments []interface{}
	if err != nil {
		return assignments, err
	}
	err = json.Unmarshal(body, &assignments)
	return assignments, err
}

func (client *Client) BrokerAssign(poolId string) (interface{}, error) {
	body, err := client.request("POST", "broker/assign/"+poolId, nil)
	var result interface{}
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

func (client *Client) BrokerConnectRDP(guest string, remote bool) ([]byte, error) {
	jsonData := map[string]interface{}{"guest": guest, "remote": remote, "outputType": "rdp"}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	body, err := client.request("POST", "broker/"+guest+"/connect", jsonValue)
	return body, err
}

func (client *Client) AssignGuest(poolId, username, realm, guest string) (interface{}, error) {
	var result interface{}
	jsonData := map[string]string{"realm": realm, "username": username, "guest": guest}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return result, err
	}
	body, err := client.request("POST", "broker/assign/"+poolId, jsonValue)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

func (client *Client) ReleaseGuest(poolId, username, guest string) error {
	jsonData := map[string]string{"poolId": poolId, "username": username, "guest": guest}
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
