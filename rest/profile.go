package rest

import (
	"encoding/json"
	"errors"

	"github.com/ghodss/yaml"
)

type Profile struct {
	AdConfig *struct {
		Domain    string      `json:"domain,omitempty"`
		Ou        interface{} `json:"ou,omitempty"`
		Password  string      `json:"password,omitempty"`
		UserGroup string      `json:"userGroup,omitempty"`
		Username  string      `json:"username,omitempty"`
	} `json:"adConfig,omitempty"`
	BrokerOptions *struct {
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
	UserVolumes  *struct {
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

func (profile *Profile) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, profile)
}

func (profile *Profile) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(profile)
	body, err := client.Request("POST", "profiles", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (profile *Profile) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(profile)
	body, err := client.Request("PUT", "profile/"+profile.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (profile *Profile) Delete(client *Client) error {
	if profile.ID == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.Request("DELETE", "profile/"+profile.ID, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) ListProfiles() ([]Profile, error) {
	var Profiles []Profile
	body, err := client.Request("GET", "profiles", nil)
	if err != nil {
		return Profiles, err
	}
	err = json.Unmarshal(body, &Profiles)
	return Profiles, err
}

func (client *Client) GetProfile(id string) (*Profile, error) {
	var profile *Profile
	if id == "" {
		return profile, errors.New("Id cannot be empty")
	}
	body, err := client.Request("GET", "profile/"+id, nil)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(body, &profile)
	return profile, err
}

func (client *Client) GetProfileByName(name string) (*Profile, error) {
	var profiles, err = client.ListProfiles()
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if profile.Name == name {
			return &profile, nil
		}
	}
	return nil, errors.New("Profile not found")
}
