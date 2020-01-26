package rest

import (
	"encoding/json"
	"errors"
)

//RealmServiceAccount contains an active directory service account for the realm
type RealmServiceAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Realm describes a realm from the rest api
type Realm struct {
	Enabled        bool                 `json:"enabled"`
	FQDN           string               `json:"fqdn"`
	Name           string               `json:"name"`
	Tags           []string             `json:"tags,omitempty"`
	Verified       bool                 `json:"verified"`
	ServiceAccount *RealmServiceAccount `json:"serviceAccount,omitempty"`
}

func (realm Realm) String() string {
	json, _ := json.MarshalIndent(realm, "", "  ")
	return string(json)
}

// ListRealms returns an array of all realms with an optional filter string
func (client *Client) ListRealms(filter string) ([]Realm, error) {
	var realms []Realm
	path := "realms"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return realms, err
	}
	err = json.Unmarshal(body, &realms)
	return realms, err
}

// GetRealm requests a realm by name
func (client *Client) GetRealm(name string) (Realm, error) {
	var realm Realm
	if name == "" {
		return realm, errors.New("Name cannot be empty")
	}
	body, err := client.request("GET", "realm/"+name, nil)
	if err != nil {
		return realm, err
	}
	err = json.Unmarshal(body, &realm)
	return realm, err
}

// Create creates a new realm
func (realm *Realm) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(realm)
	body, err := client.request("POST", "realms", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Update updates an existing realm
func (realm *Realm) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(realm)
	body, err := client.request("PUT", "realm/"+realm.Name, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Delete deletes a realm
func (realm *Realm) Delete(client *Client) error {
	if realm.Name == "" {
		return errors.New("Name cannot be empty")
	}
	_, err := client.request("DELETE", "realm/"+realm.Name, nil)
	if err != nil {
		return err
	}
	return err
}
