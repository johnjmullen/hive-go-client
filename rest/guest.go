package rest

import (
	"encoding/json"
	"errors"
)

type Guest struct {
	SessionInfo *struct {
		SessionID     int    `json:"SessionID"`
		SourceIP      string `json:"SourceIP"`
		SourceName    string `json:"SourceName"`
		SessionState  int    `json:"sessionState"`
		SessionStatus string `json:"sessionStatus"`
	} `json:"SessionInfo,omitempty"`
	AgentInstalled bool `json:"agentInstalled"`
	Cpus           int  `json:"cpus,omitempty"`
	Currentmem     int  `json:"currentmem,omitempty"`
	Disks          []struct {
		Backing    string `json:"backing,omitempty"`
		Dev        string `json:"dev,omitempty"`
		DiskDriver string `json:"diskDriver,omitempty"`
		Format     string `json:"format,omitempty"`
		Path       string `json:"path,omitempty"`
		Size       int64  `json:"size,omitempty"`
		StorageID  string `json:"storageId,omitempty"`
		Type       string `json:"type,omitempty"`
	} `json:"disks,omitempty"`
	GuestState string `json:"guestState,omitempty"`
	Hostid     string `json:"hostid,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	Interfaces []struct {
		Emulation  string `json:"emulation,omitempty"`
		IPAddress  string `json:"ipAddress,omitempty"`
		MacAddress string `json:"macAddress,omitempty"`
		Network    string `json:"network,omitempty"`
		Vlan       int    `json:"vlan,omitempty"`
	} `json:"interfaces,omitempty"`
	Memory             int      `json:"memory,omitempty"`
	Name               string   `json:"name,omitempty"`
	Os                 string   `json:"os,omitempty"`
	Persistent         bool     `json:"persistent,omitempty"`
	PoolID             string   `json:"poolId,omitempty"`
	PreviousGuestState string   `json:"previousGuestState,omitempty"`
	ProfileID          string   `json:"profileId,omitempty"`
	PublishedIP        string   `json:"publishedIp,omitempty"`
	RdpUserInjected    bool     `json:"rdpUserInjected,omitempty"`
	Realm              string   `json:"realm,omitempty"`
	Stamp              float64  `json:"stamp,omitempty"`
	Standalone         bool     `json:"standalone,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	TargetState        []string `json:"targetState,omitempty"`
	TemplateName       string   `json:"templateName,omitempty"`
	UserVolume         struct {
		State         string `json:"state,omitempty"`
		RunningBackup bool   `json:"runningBackup,omitempty"`
	} `json:"userVolume,omitempty"`
	Username string `json:"username,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	Backup   *struct {
		State           string      `json:"state,omitempty"`
		Frequency       string      `json:"frequency"`
		TargetStorageID string      `json:"targetStorageId"`
		LastBackup      interface{} `json:"lastBackup,omitempty"`
		StateMessage    string      `json:"stateMessage,omitempty"`
	} `json:"backup,omitempty"`
}

func (guest Guest) String() string {
	json, _ := json.MarshalIndent(guest, "", "  ")
	return string(json)
}

func (client *Client) ListGuests(filter string) ([]Guest, error) {
	var guests []Guest
	path := "guests"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return guests, err
	}
	err = json.Unmarshal(body, &guests)
	return guests, err
}

func (client *Client) GetGuest(name string) (*Guest, error) {
	var guest Guest
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	body, err := client.request("GET", "guest/"+name, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &guest)
	return &guest, err
}

func (guest *Guest) Shutdown(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+guest.Name+"/shutdown", nil)
	if err != nil {
		return err
	}
	return err
}

func (guest *Guest) Reboot(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+guest.Name+"/reboot", nil)
	if err != nil {
		return err
	}
	return err
}

func (guest *Guest) Poweroff(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+guest.Name+"/poweroff", nil)
	if err != nil {
		return err
	}
	return err
}

func (guest *Guest) Reset(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+guest.Name+"/reset", nil)
	if err != nil {
		return err
	}
	return err
}

func (guest *Guest) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(guest)
	body, err := client.request("PUT", "guest/"+guest.Name, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Why is this POST instad of DELETE?
func (guest *Guest) Delete(client *Client) error {
	if guest.Name == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.request("POST", "guest/"+guest.Name+"/delete", nil)
	if err != nil {
		return err
	}
	return err
}
