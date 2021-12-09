package rest

import (
	"encoding/json"
	"errors"
	"net/url"
)

// UserPreferences
type UserPreferences struct {
	UiTimeout float64 `json:"uiTimeout,omitempty"`
}

// User describes a user record
type User struct {
	ID          string           `json:"id"`
	GroupName   string           `json:"groupName,omitempty"`
	Password    string           `json:"password,omitempty"`
	Preferences *UserPreferences `json:"preferences,omitempty"`
	Realm       string           `json:"realm,omitempty"`
	Role        string           `json:"role,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Username    string           `json:"username,omitempty"`
}

func (user User) String() string {
	json, _ := json.MarshalIndent(user, "", "  ")
	return string(json)
}

// ListUsers returns an array of all users with an optional filter string
func (client *Client) ListUsers(query string) ([]User, error) {
	var users []User
	path := "users"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(body, &users)
	return users, err
}

// GetUser request a user by id
func (client *Client) GetUser(id string) (*User, error) {
	var user *User
	if id == "" {
		return user, errors.New("id cannot be empty")
	}
	body, err := client.request("GET", "user/"+id, nil)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(body, &user)
	return user, err
}

// GetUserByName request a task by name
func (client *Client) GetUserByName(name string) (*User, error) {
	var users, err = client.ListUsers("name=" + url.QueryEscape(name))
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Username == name {
			return &user, nil
		}
	}
	return nil, errors.New("User not found")
}

// Create creates a new user
func (user *User) Create(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(user)
	body, err := client.request("POST", "users", jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Update updates an existing user
func (user *User) Update(client *Client) (string, error) {
	var result string
	jsonValue, _ := json.Marshal(user)
	body, err := client.request("PUT", "user/"+user.ID, jsonValue)
	if err == nil {
		result = string(body)
	}
	return result, err
}

//Delete deletes a user
func (user *User) Delete(client *Client) error {
	if user.ID == "" {
		return errors.New("ID cannot be empty")
	}
	_, err := client.request("DELETE", "user/"+user.ID, nil)
	if err != nil {
		return err
	}
	return err
}
