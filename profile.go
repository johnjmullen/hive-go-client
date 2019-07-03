package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Profile struct {
	AdConfig struct {
		Domain    string      `json:"domain"`
		Ou        interface{} `json:"ou"`
		Password  string      `json:"password"`
		UserGroup string      `json:"userGroup"`
		Username  string      `json:"username"`
	} `json:"adConfig"`
	BrokerOptions struct {
		AllowDesktopComposition bool `json:"allowDesktopComposition"`
		AudioCapture            bool `json:"audioCapture"`
		RedirectCSSP            bool `json:"redirectCSSP"`
		RedirectClipboard       bool `json:"redirectClipboard"`
		RedirectDisk            bool `json:"redirectDisk"`
		RedirectPNP             bool `json:"redirectPNP"`
		RedirectPrinter         bool `json:"redirectPrinter"`
		RedirectUSB             bool `json:"redirectUSB"`
		SmartResize             bool `json:"smartResize"`
	} `json:"brokerOptions"`
	BypassBroker bool     `json:"bypassBroker"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Tags         []string `json:"tags"`
	Timezone     string   `json:"timezone"`
	UserVolumes  struct {
		BackupSchedule int    `json:"backupSchedule"`
		Repository     string `json:"repository"`
		Size           int    `json:"size"`
		Target         string `json:"target"`
	} `json:"userVolumes"`
	Vlan int `json:"vlan"`
}

func (client *Client) ListProfiles() ([]Profile, error) {
	var Profiles []Profile
	res, err := client.Request("GET", "profiles", nil)
	if err != nil {
		return Profiles, err
	}
	body, err := ioutil.ReadAll(res.Body)
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
	res, err := client.Request("GET", "Profile/"+id, nil)
	if err != nil {
		return Profile, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Profile, err
	}
	err = json.Unmarshal(body, &Profile)
	return Profile, err
}

func (client *Client) CreateProfile(Profile *Profile) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(Profile)
	res, err := client.Request("POST", "storage/Profiles", jsonValue)
	if err != nil {
		return result, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err == nil {
		result = string(body)
	}
	return result, err
}

func (client *Client) DeleteProfile(id string) error {
	if id == "" {
		return errors.New("Id cannot be empty")
	}
	res, err := client.Request("DELETE", "storage/Profile/"+id, nil)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return err
}
