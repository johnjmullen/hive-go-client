package rest

import (
	"encoding/json"
	"errors"
)

type Guest struct {
	SessionInfo struct {
		SessionID     int    `json:"SessionID"`
		SourceIP      string `json:"SourceIP"`
		SourceName    string `json:"SourceName"`
		SessionState  int    `json:"sessionState"`
		SessionStatus string `json:"sessionStatus"`
	} `json:"SessionInfo"`
	AgentInstalled bool `json:"agentInstalled"`
	Cpus           int  `json:"cpus"`
	Currentmem     int  `json:"currentmem"`
	Disks          []struct {
		Backing    string `json:"backing,omitempty"`
		Dev        string `json:"dev"`
		DiskDriver string `json:"diskDriver"`
		Format     string `json:"format,omitempty"`
		Path       string `json:"path,omitempty"`
		Size       int64  `json:"size,omitempty"`
		StorageID  string `json:"storageId,omitempty"`
		Type       string `json:"type"`
	} `json:"disks"`
	GuestState string `json:"guestState"`
	Hostid     string `json:"hostid"`
	Hostname   string `json:"hostname"`
	Interfaces []struct {
		Emulation  string `json:"emulation"`
		IPAddress  string `json:"ipAddress"`
		MacAddress string `json:"macAddress"`
		Network    string `json:"network"`
		Vlan       int    `json:"vlan"`
	} `json:"interfaces"`
	Memory             int      `json:"memory"`
	Name               string   `json:"name"`
	Os                 string   `json:"os"`
	Persistent         bool     `json:"persistent"`
	PoolID             string   `json:"poolId"`
	PreviousGuestState string   `json:"previousGuestState"`
	ProfileID          string   `json:"profileId"`
	PublishedIP        string   `json:"publishedIp"`
	RdpUserInjected    bool     `json:"rdpUserInjected"`
	Realm              string   `json:"realm"`
	Stamp              float64  `json:"stamp"`
	Standalone         bool     `json:"standalone"`
	Tags               []string `json:"tags"`
	TargetState        []string `json:"targetState"`
	TemplateName       string   `json:"templateName"`
	UserVolume         struct {
		State         string `json:"state"`
		RunningBackup bool   `json:"runningBackup"`
	} `json:"userVolume"`
	Username string `json:"username"`
	UUID     string `json:"uuid"`
}

/*type Guest struct {
	AgentInstalled bool `json:"agentInstalled"`
	Cpus           int  `json:"cpus"`
	Currentmem     int  `json:"currentmem"`
	Disks          []struct {
		Backing    string `json:"backing,omitempty"`
		Dev        string `json:"dev"`
		DiskDriver string `json:"diskDriver"`
		Format     string `json:"format,omitempty"`
		Path       string `json:"path,omitempty"`
		Size       int    `json:"size,omitempty"`
		StorageID  string `json:"storageId,omitempty"`
		Type       string `json:"type"`
	} `json:"disks"`
	GuestState string `json:"guestState"`
	Hostid     string `json:"hostid"`
	Hostname   string `json:"hostname"`
	Interfaces []struct {
		Emulation  string `json:"emulation"`
		MacAddress string `json:"macAddress"`
		Network    string `json:"network"`
		Vlan       int    `json:"vlan"`
	} `json:"interfaces"`
	Memory             int      `json:"memory"`
	Name               string   `json:"name"`
	Os                 string   `json:"os"`
	Persistent         bool     `json:"persistent"`
	PoolID             string   `json:"poolId"`
	PreviousGuestState string   `json:"previousGuestState"`
	ProfileID          string   `json:"profileId"`
	PublishedIP        string   `json:"publishedIp"`
	Realm              string   `json:"realm"`
	Stamp              float64  `json:"stamp"`
	Standalone         bool     `json:"standalone"`
	Tags               []string `json:"tags"`
	TargetState        []string `json:"targetState"`
	TemplateName       string   `json:"templateName"`
	Username           string   `json:"username"`
	UUID               string   `json:"uuid"`
	RdpUserInjected    bool     `json:"rdpUserInjected"`
	HostDetails        struct {
		Hostname string `json:"hostname"`
		IP       string `json:"ip"`
	} `json:"hostDetails"`
}*/

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
