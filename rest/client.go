package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type authToken struct {
	Token string `json:"token"`
}

type Client struct {
	Host          string
	Port          uint
	AllowInsecure bool
	httpClient    *http.Client
	token         string
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
		err = (fmt.Errorf("{\"error\": %d, \"message\": %s}", res.StatusCode, body))
	}
	return body, err
}

func (client *Client) Request(method, path string, data []byte) ([]byte, error) {
	protocol := "https"
	if client.Port == 3000 {
		protocol = "http"
	}
	url := fmt.Sprintf("%s://%s:%d/api/%s", protocol, client.Host, client.Port, path)
	if client.httpClient == nil {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: client.AllowInsecure},
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
	if client.Host == "localhost" || client.Host == "::1" || client.Host == "127.0.0.1" {
		return nil
	}
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
