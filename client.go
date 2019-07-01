package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// type Type interface {
//     String() string
//     Marshal() string
//     UnMarshal(string)
// }

type authToken struct {
	Token string `json:"token"`
}

type Client struct {
	Host       string
	Port       uint
	httpClient *http.Client
	token      string
}

func checkResponse(resp *http.Response, err error) (*http.Response, error) {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return resp, (fmt.Errorf("Error %d: %s", resp.StatusCode, body))
	}
	return resp, err
}

func (client *Client) Request(method, path string, data []byte) (*http.Response, error) {
	url := fmt.Sprintf("https://%s:%d/api/%s", client.Host, client.Port, path)
	log.Print(method, " ", url)
	if client.httpClient == nil {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}
		client.httpClient = &http.Client{Transport: tr, Timeout: time.Second * 30}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	if client.token != "" {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}

	return checkResponse(client.httpClient.Do(req))

}

func (client *Client) Login(username, password, realm string) error {
	jsonData := map[string]string{"username": username, "password": password, "realm": realm}
	jsonValue, _ := json.Marshal(jsonData)
	var err error
	resp, err := client.Request("POST", "auth", jsonValue)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	auth := authToken{}
	err = json.Unmarshal(body, &auth)
	if err == nil {
		client.token = auth.Token
	}
	return err
}
