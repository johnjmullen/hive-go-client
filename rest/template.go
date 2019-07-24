package rest

import (
	"encoding/json"
	"errors"
)

type TemplateInterface struct {
	Network   string `json:"network"`
	Vlan      int    `json:"vlan"`
	Emulation string `json:"emulation"`
}

type TemplateDisk struct {
	DiskDriver string `json:"diskDriver"`
	Filename   string `json:"filename"`
	Format     string `json:"format,omitempty"`
	Path       string `json:"path,omitempty"` //TODO: Is this actually used?
	Size       int    `json:"size"`
	StorageID  string `json:"storageId"`
	Type       string `json:"type"`
}

type Template struct {
	ID                 string               `json:"id,omitempty"`
	Name               string               `json:"name"`
	Vcpu               int                  `json:"vcpu"`
	Mem                int                  `json:"mem"`
	OS                 string               `json:"os"`
	Firmware           string               `json:"firmware,omitempty"`
	DisplayDriver      string               `json:"displayDriver,omitempty"`
	Interfaces         []*TemplateInterface `json:"interfaces,omitempty"`
	Drivers            bool                 `json:"drivers"`
	Disks              []*TemplateDisk      `json:"disks,omitempty"`
	State              string               `json:"state,omitempty"`
	StateMessage       string               `json:"stateMessage,omitempty"`
	ManualAgentInstall bool                 `json:"manualAgentInstall"`
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

func (template *Template) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(template)
	body, err := client.Request("POST", "templates", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (template *Template) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(template)
	body, err := client.Request("PUT", "template/"+template.Name, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (template *Template) Delete(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("DELETE", "template/"+template.Name, nil)
	if err != nil {
		return err
	}
	return err
}

func (template *Template) Load(client *Client, storage string) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}

	jsonData := map[string]string{"localStorage": storage}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := client.Request("POST", "template/"+template.Name+"/loadall", jsonValue)
	return err
}

func (template *Template) Unload(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("POST", "template/"+template.Name+"/unloadall", nil)
	return err
}

func (template *Template) Analyze(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("POST", "template/"+template.Name+"/analyze", nil)
	return err
}

func (template *Template) Author(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("PUT", "template/"+template.Name+"/author", nil)
	return err
}

/*func (template *Template) Duplicate(client *Client, dstName, dstStorage, dstFilename) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	jsonData := map[string]string{"dstName": dstName, "dstStorage": dstStorage, "dstFilename": dstFilename, "srcStorage":template. }
	jsonValue, _ := json.Marshal(jsonData)
	_, err := client.Request("POST", "template/"+template.Name+"/author", jsonValue)
	return err
}*/
