package client

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Task struct {
	Enabled  bool     `json:"enabled"`
	FQDN     string   `json:"fqdn"`
	Name     string   `json:"id"`
	Tags     []string `json:"tags,omitempty"`
	Verified bool     `json:"verified"`
}

func (task Task) String() string {
	return fmt.Sprintf("{\n Name: %v,\n Enabled: %v,\n fqdn: %v,\n Tags: %v,\n Verified: %v\n}\n", task.Name, task.Enabled, task.FQDN, task.Tags, task.Verified)
}

func (task *Task) ToJson() ([]byte, error) {
	return json.Marshal(task)
}

func (task *Task) FromJson(data []byte) error {
	return json.Unmarshal(data, task)
}

func (client *Client) ListTasks() ([]Task, error) {
	var tasks []Task
	body, err := client.Request("GET", "tasks", nil)
	if err != nil {
		return tasks, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &tasks)
	return tasks, err
}

func (client *Client) GetTask(id string) (Task, error) {
	var task Task
	if id == "" {
		return task, errors.New("Name cannot be empty")
	}
	body, err := client.Request("GET", "task/"+id, nil)
	if err != nil {
		return task, err
	}
	err = json.Unmarshal(body, &task)
	return task, err
}

func (client *Client) ForceCompleteTask(id string) error {
	if id == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.Request("PUT", "task/"+id+"/forcecomplete", nil)
	return err
}
