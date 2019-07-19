package rest

import (
	"encoding/json"
	"errors"
)

type Realm struct {
	Enabled  bool     `json:"enabled"`
	FQDN     string   `json:"fqdn"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags,omitempty"`
	Verified bool     `json:"verified"`
}

func (realm Realm) String() string {
	json, _ := json.MarshalIndent(realm, "", "  ")
	return string(json)
}

func (client *Client) ListRealms() ([]Realm, error) {
	var realms []Realm
	body, err := client.Request("GET", "realms", nil)
	if err != nil {
		return realms, err
	}
	err = json.Unmarshal(body, &realms)
	return realms, err
}

func (client *Client) GetRealm(name string) (Realm, error) {
	var realm Realm
	if name == "" {
		return realm, errors.New("Name cannot be empty")
	}
	body, err := client.Request("GET", "realm/"+name, nil)
	if err != nil {
		return realm, err
	}
	err = json.Unmarshal(body, &realm)
	return realm, err
}

func (realm *Realm) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(realm)
	body, err := client.Request("POST", "realms", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (realm *Realm) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(realm)
	body, err := client.Request("PUT", "realm/"+realm.Name, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (realm *Realm) Delete(client *Client) error {
	if realm.Name == "" {
		return errors.New("Name cannot be empty")
	}
	_, err := client.Request("DELETE", "realm/"+realm.Name, nil)
	if err != nil {
		return err
	}
	return err
}
