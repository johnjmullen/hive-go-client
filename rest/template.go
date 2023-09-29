package rest

import (
	"encoding/json"
	"errors"
	"net/url"
)

// TemplateInterface a network interface from a template record
type TemplateInterface struct {
	Network   string `json:"network"`
	Vlan      int    `json:"vlan"`
	Emulation string `json:"emulation"`
}

// TemplateDisk a disk from a template record
type TemplateDisk struct {
	DiskDriver string `json:"diskDriver"`
	Filename   string `json:"filename"`
	Format     string `json:"format,omitempty"`
	Path       string `json:"path,omitempty"` //TODO: Is this actually used?
	Size       int    `json:"size"`
	StorageID  string `json:"storageId"`
	Type       string `json:"type"`
	OsVolume   int    `json:"osvolume"`
}

// Template a template record from the rest interface
type Template struct {
	Name               string                 `json:"name"`
	Vcpu               int                    `json:"vcpu"`
	Mem                int                    `json:"mem"`
	OS                 string                 `json:"os"`
	Firmware           string                 `json:"firmware,omitempty"`
	Gpu                bool                   `json:"gpu,omitempty"`
	DisplayDriver      string                 `json:"displayDriver,omitempty"`
	Interfaces         []*TemplateInterface   `json:"interfaces,omitempty"`
	Drivers            bool                   `json:"drivers"`
	Disks              []*TemplateDisk        `json:"disks,omitempty"`
	State              string                 `json:"state,omitempty"`
	StateMessage       string                 `json:"stateMessage,omitempty"`
	ManualAgentInstall bool                   `json:"manualAgentInstall"`
	TemplateMap        map[string]interface{} `json:"templateMap,omitempty"`
	BrokerOptions      GuestBrokerOptions     `json:"brokerOptions,omitempty"`
}

func (template Template) String() string {
	json, _ := json.MarshalIndent(template, "", "  ")
	return string(json)
}

// ListTemplates returns an array of templates with an optional filter string
func (client *Client) ListTemplates(query string) ([]Template, error) {
	path := "templates"
	if query != "" {
		path += "?" + query
	}
	var templates []Template
	body, err := client.request("GET", path, nil)
	if err != nil {
		return templates, err
	}
	err = json.Unmarshal(body, &templates)
	return templates, err
}

// GetTemplate request a template by name
func (client *Client) GetTemplate(name string) (Template, error) {
	var template Template
	if name == "" {
		return template, errors.New("name cannot be empty")
	}
	body, err := client.request("GET", "template/"+url.PathEscape(name), nil)
	if err != nil {
		return template, err
	}
	err = json.Unmarshal(body, &template)
	return template, err
}

// Create creates a new template
func (template *Template) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(template)
	body, err := client.request("POST", "templates", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Update updates an existing template
func (template *Template) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(template)
	body, err := client.request("PUT", "template/"+url.PathEscape(template.Name), jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

// Delete deletes a template record
func (template *Template) Delete(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("DELETE", "template/"+url.PathEscape(template.Name), nil)
	if err != nil {
		return err
	}
	return err
}

// Load stages a template to storage across the cluster
// stoarge can be "disk" or "ram"
func (template *Template) Load(client *Client, storage string) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}

	jsonData := map[string]string{"localStorage": storage}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := client.request("POST", "template/"+url.PathEscape(template.Name)+"/loadall", jsonValue)
	return err
}

// Unload removes a staged template from local storage across the cluster
func (template *Template) Unload(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "template/"+url.PathEscape(template.Name)+"/unloadall", nil)
	return err
}

// Analyze validates a template
func (template *Template) Analyze(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "template/"+url.PathEscape(template.Name)+"/analyze", nil)
	return err
}

// Author creates a virtual machine to author a template
func (template *Template) Author(client *Client) error {
	if template.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("PUT", "template/"+url.PathEscape(template.Name)+"/author", nil)
	return err
}

// Duplicate copies a template
func (template *Template) Duplicate(client *Client, dstName, dstStorage, dstFilename string) (*Task, error) {
	if template.Name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if len(template.Disks) < 1 {
		return nil, errors.New("no disks found")
	}
	jsonData := map[string]string{"srcStorage": template.Disks[0].StorageID, "srcFilename": template.Disks[0].Filename, "dstName": dstName, "dstStorage": dstStorage, "dstFilename": dstFilename}
	jsonValue, _ := json.Marshal(jsonData)
	return client.getTaskFromResponse(client.request("POST", "template/"+url.PathEscape(template.Name)+"/duplicate", jsonValue))
}
