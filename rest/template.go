package client

import (
	"encoding/json"
	"errors"
)

type Template struct {
	Vcpu          int    `json:"vcpu"`
	Mem           int    `json:"mem"`
	Os            string `json:"os"`
	Firmware      string `json:"firmware"`
	DisplayDriver string `json:"displayDriver"`
	Name          string `json:"name"`
	Interfaces    []struct {
		Network   string `json:"network"`
		Vlan      string `json:"vlan"`
		Emulation string `json:"emulation"`
	} `json:"interfaces"`
	Drivers bool `json:"drivers"`
	Disks   []struct {
		Type      string `json:"type"`
		StorageID string `json:"storageId"`
		Filename  string `json:"filename"`
		Emulation string `json:"emulation"`
		DiskSize  int    `json:"diskSize"`
	} `json:"disks"`
}

func (template Template) String() string {
	json, _ := json.MarshalIndent(template, "", "  ")
	return string(json)
}

func (template *Template) ToJson() ([]byte, error) {
	return json.Marshal(template)
}

func (template *Template) FromJson(data []byte) error {
	return json.Unmarshal(data, template)
}

func (client *Client) ListTemplates() ([]Template, error) {
	var templates []Template
	body, err := client.Request("GET", "templates", nil)
	if err != nil {
		return templates, err
	}
	err = json.Unmarshal(body, &templates)
	return templates, err
}

func (client *Client) GetTemplate(name string) (Template, error) {
	var template Template
	if name == "" {
		return template, errors.New("name cannot be empty")
	}
	body, err := client.Request("GET", "template/"+name, nil)
	if err != nil {
		return template, err
	}
	err = json.Unmarshal(body, &template)
	return template, err
}

func (client *Client) CreateTemplate(template *Template) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(template)
	body, err := client.Request("POST", "templates", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (client *Client) DeleteTemplate(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("DELETE", "template/"+name, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) LoadTemplate(name, storage string) ([]byte, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	jsonData := map[string]string{"localStorage": storage}
	jsonValue, _ := json.Marshal(jsonData)
	return client.Request("POST", "template/"+name+"/loadall", jsonValue)
}

func (client *Client) UnloadTemplate(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("POST", "template/"+name+"/unloadall", nil)
	return err
}
