package client

import (
	"encoding/json"
	"errors"
)

type Guest struct {
	Name string `json:"name"`
}

func (guest Guest) String() string {
	json, _ := json.MarshalIndent(guest, "", "  ")
	return string(json)
}

func (guest *Guest) ToJson() ([]byte, error) {
	return json.Marshal(guest)
}

func (guest *Guest) FromJson(data []byte) error {
	return json.Unmarshal(data, guest)
}

func (client *Client) ListGuests() ([]Guest, error) {
	var guests []Guest
	body, err := client.Request("GET", "guests", nil)
	if err != nil {
		return guests, err
	}
	err = json.Unmarshal(body, &guests)
	return guests, err
}

func (client *Client) GetGuest(guestid string) (Guest, error) {
	var guest Guest
	if guestid == "" {
		return guest, errors.New("guestid cannot be empty")
	}
	body, err := client.Request("GET", "guest/"+guestid, nil)
	if err != nil {
		return guest, err
	}
	err = json.Unmarshal(body, &guest)
	return guest, err
}

func (client *Client) DeleteGuest(guestid string) error {
	if guestid == "" {
		return errors.New("name cannot be empty")
	}
	_, err := client.Request("DELETE", "guest/"+guestid, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) GuestVersion() (Version, error) {
	var version Version
	body, err := client.Request("GET", "guest/version", nil)
	if err != nil {
		return version, err
	}
	err = json.Unmarshal(body, &version)
	return version, err
}
