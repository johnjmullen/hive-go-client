package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (client *Client) getTaskFromResponse(body []byte, err error) (*Task, error) {
	if err != nil {
		return nil, err
	}
	var objMap map[string]string
	err = json.Unmarshal(body, &objMap)
	taskId, ok := objMap["taskId"]
	if err != nil || !ok {
		return nil, fmt.Errorf("Error parsing data.  taskId not found")
	}
	return client.GetTask(taskId)
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

func (client *Client) request(method, path string, data []byte) ([]byte, error) {
	protocol := "https"
	if client.Port == 3000 {
		protocol = "http"
	}
	//TODO: separate queryString from path in function arguments
	urlString := fmt.Sprintf("%s://%s:%d/api/%s", protocol, client.Host, client.Port, path)
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	if client.httpClient == nil {
		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: client.AllowInsecure},
			DisableCompression: true,
		}
		client.httpClient = &http.Client{Transport: tr, Timeout: time.Second * 30}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(data))
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
	body, err := client.request("POST", "auth", jsonValue)
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
