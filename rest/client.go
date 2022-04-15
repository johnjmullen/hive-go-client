package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-version"
)

type authToken struct {
	Token string `json:"token"`
}

// Client is a wrapper around the hiveio rest api
type Client struct {
	Host          string
	Port          uint
	AllowInsecure bool
	httpClient    *http.Client
	token         string
}

//SetToken sets the token directly instead of calling auth
func (client *Client) SetToken(token string) {
	client.token = token
}

func (client *Client) getTaskFromResponse(body []byte, err error) (*Task, error) {
	if err != nil {
		return nil, err
	}
	var objMap map[string]string
	err = json.Unmarshal(body, &objMap)
	taskID, ok := objMap["taskId"]
	if err != nil || !ok || taskID == "" {
		return nil, fmt.Errorf("error parsing data. taskId not found")
	}
	return client.GetTask(taskID)
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
	headers := map[string]string{"Content-type": "application/json"}
	return checkResponse(client.requestWithHeaders(method, path, bytes.NewBuffer(data), headers, time.Second*120))
}

func (client *Client) postMultipart(path, filenameField, filepath string, params map[string]string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	//Add file and params to the form
	body := bytes.NewBufferString("")
	writer := multipart.NewWriter(body)
	_, err = writer.CreateFormFile(filenameField, info.Name())
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		err = writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}

	boundary := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", writer.Boundary()))
	mreader := io.MultiReader(body, f, boundary)
	headers := map[string]string{
		"Content-Type": fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary()),
	}
	return checkResponse(client.requestWithHeaders("POST", path, mreader, headers, time.Second*30))
	//req.ContentLength = fi.Size()+int64(body_buf.Len())+int64(close_buf.Len())
}

func (client *Client) requestWithHeaders(method, path string, body io.Reader, headers map[string]string, timeout time.Duration) (*http.Response, error) {
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
		client.httpClient = &http.Client{Transport: tr}
	}
	client.httpClient.Timeout = timeout

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	if client.token != "" {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}
	return client.httpClient.Do(req)
}

// Login attempts to connect to the server specified in Client with the provided username, password, and realm
func (client *Client) Login(username, password, realm string) error {
	if password == "" && (client.Host == "localhost" || client.Host == "::1" || client.Host == "127.0.0.1") {
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

// ChangeFeed wrapper around a websocket to monitor database changes
type ChangeFeed struct {
	Data chan ChangeFeedMessage
	Done chan struct{}
	conn *websocket.Conn
}

// ChangeFeedMessage contains a change from the database
type ChangeFeedMessage struct {
	OldValue json.RawMessage `json:"old_val"`
	NewValue json.RawMessage `json:"new_val"`
	Error    error
}

func (feed *ChangeFeed) monitorChangeFeed() {
	defer close(feed.Done)
	defer feed.conn.Close()
	for {
		var msg ChangeFeedMessage
		_, message, err := feed.conn.ReadMessage()
		if err != nil {
			msg.Error = err
			feed.Data <- msg
			return
		}
		if len(message) < 3 {
			continue
		}
		if string(message[:2]) == "44" {
			msg.Error = fmt.Errorf("%s", message[2:])
			feed.Data <- msg
		} else if string(message[:2]) != "42" {
			continue
		}

		var jsonMsg []json.RawMessage
		err = json.Unmarshal(message[2:], &jsonMsg)
		if err != nil {
			msg.Error = err
			feed.Data <- msg
			continue
		}
		if len(jsonMsg) < 3 {
			msg.Error = fmt.Errorf("invalid changefeed message")
			feed.Data <- msg
			continue
		}
		if strings.Contains(string(jsonMsg[0]), "initial") {
			var newValue []json.RawMessage
			err = json.Unmarshal(jsonMsg[2], &newValue)
			if err != nil {
				msg.Error = err
			} else if len(newValue) > 0 {
				msg.NewValue = newValue[0]
			}
		} else {
			err = json.Unmarshal(jsonMsg[2], &msg)
			if err != nil {
				msg.Error = err
			}
		}
		feed.Data <- msg
	}
}

func (feed *ChangeFeed) changeFeedKeepAlive(timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			feed.conn.WriteMessage(websocket.TextMessage, []byte("2"))
		case <-feed.Done:
			return
		}
	}
}

// Close disconnects the changefeed websocket
func (feed *ChangeFeed) Close() error {
	return feed.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

// GetChangeFeed returns a ChangeFeed for monitoring the specified table
// filter can be used to limit the changes monitored.
// example to monitor a single task:
// client.GetChangeFeed("task", map[string]string{"id": task.ID})
func (client *Client) GetChangeFeed(table string, filter map[string]string, includeInitial bool) (*ChangeFeed, error) {
	protocol := "wss"
	var token string
	if client.Port == 3000 {
		protocol = "ws"
	}
	if client.token != "" {
		token = "token=" + client.token + "&"
	}
	u := url.URL{Scheme: protocol, Host: fmt.Sprintf("%s:%d", client.Host, client.Port), Path: "/socket.io/", RawQuery: token + "transport=websocket"}
	dialer := websocket.Dialer{
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		HandshakeTimeout: 20 * time.Second,
	}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	options := map[string]interface{}{"table": table, "includeInitial": includeInitial, "filter": filter}
	jsonData := []interface{}{"query:change:register", options}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	err = c.WriteMessage(websocket.TextMessage, append([]byte("42"), jsonValue...))
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})
	incomingData := make(chan ChangeFeedMessage)
	feed := ChangeFeed{Data: incomingData, Done: done, conn: c}

	go feed.changeFeedKeepAlive(25 * time.Second)
	go feed.monitorChangeFeed()

	return &feed, nil
}

// HostVersion returns the software version of the host the client is connected to
func (client *Client) HostVersion() (Version, error) {
	var version Version
	body, err := client.request("GET", "host/version", nil)
	if err != nil {
		return version, err
	}
	err = json.Unmarshal(body, &version)
	return version, err
}

// HostID returns the hostid of the host the client is connected to
func (client *Client) HostID() (string, error) {
	body, err := client.request("GET", "host/hostid", nil)
	if err != nil {
		return "", err
	}
	var objMap map[string]string
	err = json.Unmarshal(body, &objMap)
	return objMap["id"], err
}

// ClusterID returns the cluster id of the host the client is connected to
func (client *Client) ClusterID() (string, error) {
	body, err := client.request("GET", "host/clusterid", nil)
	if err != nil {
		return "", err
	}
	var objMap map[string]string
	err = json.Unmarshal(body, &objMap)
	return objMap["id"], err
}

// CheckHostVersion returns an error if minimumVersion is greater than the host software version
func (client *Client) CheckHostVersion(minimumVersion string) error {
	hostVersion, err := client.HostVersion()
	if err != nil {
		return err
	}
	minVersion, err := version.NewVersion(minimumVersion)
	if err != nil {
		return err
	}
	v, err := version.NewVersion(hostVersion.Version)
	if err != nil {
		return fmt.Errorf("failed to parse host software version: %s", hostVersion.Version)
	}
	if v.LessThan(minVersion) {
		return fmt.Errorf("unsupported host version %s. Requires %s or greater", hostVersion.Version, minimumVersion)
	}
	return nil
}
