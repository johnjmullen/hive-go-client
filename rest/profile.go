package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Profile struct {
	AdConfig struct {
		Domain    string      `json:"domain,omitempty"`
		Ou        interface{} `json:"ou,omitempty,omitempty"`
		Password  string      `json:"password,omitempty"`
		UserGroup string      `json:"userGroup,omitempty"`
		Username  string      `json:"username,omitempty"`
	} `json:"adConfig,omitempty"`
	BrokerOptions struct {
		AllowDesktopComposition bool `json:"allowDesktopComposition,omitempty"`
		AudioCapture            bool `json:"audioCapture,omitempty"`
		RedirectCSSP            bool `json:"redirectCSSP,omitempty"`
		RedirectClipboard       bool `json:"redirectClipboard,omitempty"`
		RedirectDisk            bool `json:"redirectDisk,omitempty"`
		RedirectPNP             bool `json:"redirectPNP,omitempty"`
		RedirectPrinter         bool `json:"redirectPrinter,omitempty"`
		RedirectUSB             bool `json:"redirectUSB,omitempty"`
		SmartResize             bool `json:"smartResize,omitempty"`
	} `json:"brokerOptions,omitempty"`
	BypassBroker bool     `json:"bypassBroker,omitempty"`
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	Tags         []string `json:"tags,omitempty"`
	Timezone     string   `json:"timezone,omitempty"`
	UserVolumes  struct {
		BackupSchedule int    `json:"backupSchedule,omitempty"`
		Repository     string `json:"repository,omitempty"`
		Size           int    `json:"size,omitempty"`
		Target         string `json:"target,omitempty"`
	} `json:"userVolumes,omitempty"`
	Vlan int `json:"vlan,omitempty"`
}

func (profile Profile) String() string {
	json, _ := json.MarshalIndent(profile, "", "  ")
	return string(json)
}

func (profile *Profile) ToJson() ([]byte, error) {
	return json.Marshal(profile)
}

func (profile *Profile) FromJson(data []byte) error {
	return json.Unmarshal(data, profile)
}

func (client *Client) ListProfiles() ([]Profile, error) {
	var Profiles []Profile
	body, err := client.Request("GET", "profiles", nil)
	if err != nil {
		return Profiles, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &Profiles)
	return Profiles, err
}

func (client *Client) GetProfile(id string) (Profile, error) {
	var Profile Profile
	if id == "" {
		return Profile, errors.New("Id cannot be empty")
	}
	body, err := client.Request("GET", "profile/"+id, nil)
	if err != nil {
		return Profile, err
	}
	err = json.Unmarshal(body, &Profile)
	return Profile, err
}

func (client *Client) CreateProfile(Profile *Profile) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(Profile)
	body, err := client.Request("POST", "profiles", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (client *Client) DeleteProfile(id string) error {
	if id == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.Request("DELETE", "profile/"+id, nil)
	if err != nil {
		return err
	}
	return err
}
