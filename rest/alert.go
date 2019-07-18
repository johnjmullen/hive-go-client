package rest

import (
	"encoding/json"
	"errors"
	"time"
)

type Alert struct {
	Acknowledged bool   `json:"acknowledged"`
	Class        string `json:"class"`
	Count        int    `json:"count"`
	ID           string `json:"id"`
	Level        string `json:"level"`
	Message      string `json:"message"`
	Ref          struct {
		Cluster string `json:"cluster"`
		Host    string `json:"host"`
	} `json:"ref"`
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

func (alert Alert) String() string {
	json, _ := json.MarshalIndent(alert, "", "  ")
	return string(json)
}

func (alert *Alert) ToJson() ([]byte, error) {
	return json.Marshal(alert)
}

func (alert *Alert) FromJson(data []byte) error {
	return json.Unmarshal(data, alert)
}

func (client *Client) ListAlerts() ([]Alert, error) {
	var alerts []Alert
	body, err := client.Request("GET", "alerts", nil)
	if err != nil {
		return alerts, err
	}
	err = json.Unmarshal(body, &alerts)
	return alerts, err
}

func (client *Client) GetAlert(id string) (Alert, error) {
	var alert Alert
	if id == "" {
		return alert, errors.New("Name cannot be empty")
	}
	body, err := client.Request("GET", "alert/"+id, nil)
	if err != nil {
		return alert, err
	}
	err = json.Unmarshal(body, &alert)
	return alert, err
}

func (alert *Alert) Acknowledge(client *Client) error {
	if alert.ID == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.Request("POST", "alert/"+alert.ID+"/acknowledge", nil)
	return err
}
