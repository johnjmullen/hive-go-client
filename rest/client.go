package client

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

type Type interface {
	ToJson() ([]byte, error)
	FromJson([]byte) error
}

type authToken struct {
	Token string `json:"token"`
}

type Client struct {
	Host       string
	Port       uint
	httpClient *http.Client
	token      string
}

func NewClient(host string, port uint) *Client {
	return &Client{Host: host, Port: port}
}

func checkResponse(res *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = (fmt.Errorf("Error %d: %s", res.StatusCode, body))
	}
	return body, err
}

func (client *Client) Request(method, path string, data []byte) ([]byte, error) {
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
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	body, err := client.Request("POST", "auth", jsonValue)
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
