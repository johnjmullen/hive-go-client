package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

/*func (template Template) String() string {
	return fmt.templaterintf("{\n Name: %s,\n Id: %s,\n Type: %s,\n Server: %s,\n Path: %s\n}\n", template.Name, template.Type, template.Server, template.Path)
}*/

func (client *Client) ListTemplates() ([]Template, error) {
	var templates []Template
	res, err := client.Request("GET", "templates", nil)
	if err != nil {
		return templates, err
	}
	body, err := ioutil.ReadAll(res.Body)
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
	res, err := client.Request("GET", "template/"+name, nil)
	if err != nil {
		return template, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return template, err
	}
	err = json.Unmarshal(body, &template)
	return template, err
}

func (client *Client) CreateTemplate(template *Template) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(template)
	res, err := client.Request("POST", "templates", jsonValue)
	if err != nil {
		return result, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (client *Client) DeleteTemplate(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	res, err := client.Request("DELETE", "template/"+name, nil)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return err
}
