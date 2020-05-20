package rest

import (
	"encoding/json"
	"errors"
)

//ProfileADConfig contains the active directory settings for a profile
type ProfileADConfig struct {
	Domain    string `json:"domain,omitempty"`
	Ou        string `json:"ou,omitempty"`
	Password  string `json:"password,omitempty"`
	UserGroup string `json:"userGroup,omitempty"`
	Username  string `json:"username,omitempty"`
}

//ProfileBrokerOptions contains rdp settings for a profile
type ProfileBrokerOptions struct {
	AllowDesktopComposition   bool `json:"allowDesktopComposition,omitempty"`
	AudioCapture              bool `json:"audioCapture,omitempty"`
	DisableFullWindowDrag     bool `json:"disableFullWindowDrag,omitempty"`
	DisableMenuAnims          bool `json:"disableMenuAnims,omitempty"`
	DisablePrinter            bool `json:"disablePrinter,omitempty"`
	DisableThemes             bool `json:"disableThemes,omitempty"`
	DisableWallpaper          bool `json:"disableWallpaper,omitempty"`
	HideAuthenticationFailure bool `json:"hideAuthenticationFailure,omitempty"`
	InjectPassword            bool `json:"injectPassword,omitempty"`
	RedirectCSSP              bool `json:"redirectCSSP,omitempty"`
	RedirectClipboard         bool `json:"redirectClipboard,omitempty"`
	RedirectDisk              bool `json:"redirectDisk,omitempty"`
	RedirectPNP               bool `json:"redirectPNP,omitempty"`
	RedirectPrinter           bool `json:"redirectPrinter,omitempty"`
	RedirectSmartCard         bool `json:"redirectSmartCard,omitempty"`
	RedirectUSB               bool `json:"redirectUSB,omitempty"`
	SmartResize               bool `json:"smartResize,omitempty"`
	FailOnCertMismatch        bool `json:"failOnCertMismatch,omitempty"`
}

//ProfileUserVolumes contains user volume settings for a profile
type ProfileUserVolumes struct {
	BackupSchedule int    `json:"backupSchedule,omitempty"`
	Repository     string `json:"repository,omitempty"`
	Size           int    `json:"size,omitempty"`
	Target         string `json:"target,omitempty"`
}

//ProfileBackup contains data protection settings for a profile
type ProfileBackup struct {
	Enabled         bool        `json:"enabled"`
	Frequency       string      `json:"frequency"`
	TargetStorageID string      `json:"targetStorageId"`
	UserVolumeList  []string    `json:"userVolumeList,omitempty"`
	LastBackup      interface{} `json:"date,omitempty"`
}

// Profile is a profile record from the rest api
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
	EnableHTML5   bool                  `json:"enableHtml5,omitempty"`
}

func (profile Profile) String() string {
	json, _ := json.MarshalIndent(profile, "", "  ")
	return string(json)
}

//Create creates a new profile
func (profile *Profile) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(profile)
	body, err := client.request("POST", "profiles", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Update updates an existing profile
func (profile *Profile) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(profile)
	body, err := client.request("PUT", "profile/"+profile.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Delete deletes a profile
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

// ListProfiles returns an array of all profiles with an optional filter string
func (client *Client) ListProfiles(query string) ([]Profile, error) {
	var Profiles []Profile
	path := "profiles"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return Profiles, err
	}
	err = json.Unmarshal(body, &Profiles)
	return Profiles, err
}

// GetProfile requests a profile by id
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

// GetProfileByName requests a profile by name
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
