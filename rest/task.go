package client

import (
	"encoding/json"
	"errors"
	"time"
)

type Task struct {
	Cancellable     bool      `json:"cancellable"`
	Description     string    `json:"description"`
	FinishedTime    time.Time `json:"finishedTime"`
	ID              string    `json:"id"`
	LastUpdatedTime time.Time `json:"lastUpdatedTime"`
	Message         string    `json:"message"`
	Name            string    `json:"name"`
	Progress        int       `json:"progress"`
	QueueTime       int       `json:"queueTime"`
	Ref             struct {
		Cluster string `json:"cluster"`
		Host    string `json:"host"`
	} `json:"ref"`
	StartTime time.Time   `json:"startTime"`
	State     string      `json:"state"`
	Tags      []string    `json:"tags"`
	Type      string      `json:"type"`
	Username  interface{} `json:"username"`
}

func (task Task) String() string {
	json, _ := json.MarshalIndent(task, "", "  ")
	return string(json)
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

func (task *Task) ForceCompleteTask(client *Client) error {
	if task.ID == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.Request("PUT", "task/"+task.ID+"/forcecomplete", nil)
	return err
}
