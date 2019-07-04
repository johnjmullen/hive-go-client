package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Realm struct {
	Enabled  bool     `json:"enabled"`
	FQDN     string   `json:"fqdn"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags,omitempty"`
	Verified bool     `json:"verified"`
}

func (realm Realm) String() string {
	return fmt.Sprintf("{\n Name: %v,\n Enabled: %v,\n fqdn: %v,\n Tags: %v,\n Verified: %v\n}\n", realm.Name, realm.Enabled, realm.FQDN, realm.Tags, realm.Verified)
}

func (realm *Realm) ToJson() ([]byte, error) {
	return json.Marshal(realm)
}

func (realm *Realm) FromJson(data []byte) error {
	return json.Unmarshal(data, realm)
}

func (client *Client) ListRealms() ([]Realm, error) {
	var realms []Realm
	body, err := client.Request("GET", "realms", nil)
	if err != nil {
		return realms, err
	}
	fmt.Println(string(body))
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

func (client *Client) CreateRealm(realm *Realm) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(realm)
	body, err := client.Request("POST", "realms", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (client *Client) DeleteRealm(name string) error {
	if name == "" {
		return errors.New("Name cannot be empty")
	}
	_, err := client.Request("DELETE", "realm/"+name, nil)
	if err != nil {
		return err
	}
	return err
}
