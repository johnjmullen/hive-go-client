package rest

import (
	"encoding/json"
	"errors"
)

type ProfileADConfig struct {
	Domain    string `json:"domain,omitempty"`
	Ou        string `json:"ou,omitempty"`
	Password  string `json:"password,omitempty"`
	UserGroup string `json:"userGroup,omitempty"`
	Username  string `json:"username,omitempty"`
}

type ProfileBrokerOptions struct {
	AllowDesktopComposition bool `json:"allowDesktopComposition,omitempty"`
	AudioCapture            bool `json:"audioCapture,omitempty"`
	RedirectCSSP            bool `json:"redirectCSSP,omitempty"`
	RedirectClipboard       bool `json:"redirectClipboard,omitempty"`
	RedirectDisk            bool `json:"redirectDisk,omitempty"`
	RedirectPNP             bool `json:"redirectPNP,omitempty"`
	RedirectPrinter         bool `json:"redirectPrinter,omitempty"`
	RedirectUSB             bool `json:"redirectUSB,omitempty"`
	SmartResize             bool `json:"smartResize,omitempty"`
}

type ProfileUserVolumes struct {
	BackupSchedule int    `json:"backupSchedule,omitempty"`
	Repository     string `json:"repository,omitempty"`
	Size           int    `json:"size,omitempty"`
	Target         string `json:"target,omitempty"`
}

type ProfileBackup struct {
	Frequency       string      `json:"frequency"`
	TargetStorageID string      `json:"targetStorageId"`
	UserVolumeList  []string    `json:"userVolumeList,omitempty"`
	LastBackup      interface{} `json:"date,omitempty"`
}

type Profile struct {
	AdConfig      *ProfileADConfig      `json:"adConfig,omitempty"`
	BrokerOptions *ProfileBrokerOptions `json:"brokerOptions,omitempty"`
	BypassBroker  bool                  `json:"bypassBroker"`
	ID            string                `json:"id,omitempty"`
	Name          string                `json:"name"`
	Tags          []string              `json:"tags,omitempty"`
	Timezone      string                `json:"timezone,omitempty"`
	UserVolumes   *ProfileUserVolumes   `json:"userVolumes,omitempty"`
	Vlan          int                   `json:"vlan,omitempty"`
	Backup        *ProfileBackup        `json:"backup,omitempty"`
}

func (profile Profile) String() string {
	json, _ := json.MarshalIndent(profile, "", "  ")
	return string(json)
}

func (profile *Profile) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(profile)
	body, err := client.request("POST", "profiles", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (profile *Profile) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(profile)
	body, err := client.request("PUT", "profile/"+profile.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (profile *Profile) Delete(client *Client) error {
	if profile.ID == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.request("DELETE", "profile/"+profile.ID, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) ListProfiles(filter string) ([]Profile, error) {
	var Profiles []Profile
	path := "profiles"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.request("GET", path, nil)
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
	body, err := client.request("GET", "profile/"+id, nil)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(body, &profile)
	return profile, err
}

func (client *Client) GetProfileByName(name string) (*Profile, error) {
	var profiles, err = client.ListProfiles("name=" + name)
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
