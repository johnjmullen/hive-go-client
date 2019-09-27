package rest

import (
	"encoding/json"
	"errors"
	"fmt"
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
	StartTime time.Time `json:"startTime"`
	State     string    `json:"state"`
	Tags      []string  `json:"tags"`
	Type      string    `json:"type"`
	Username  string    `json:"username"`
}

func (task Task) String() string {
	json, _ := json.MarshalIndent(task, "", "  ")
	return string(json)
}

func (client *Client) ListTasks(filter string) ([]Task, error) {
	var tasks []Task
	path := "tasks"
	if filter != "" {
		path += "?" + filter
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return tasks, err
	}
	err = json.Unmarshal(body, &tasks)
	return tasks, err
}

func (client *Client) GetTask(id string) (*Task, error) {
	var task *Task
	if id == "" {
		return task, errors.New("Id cannot be empty")
	}
	body, err := client.request("GET", "task/"+id, nil)
	if err != nil {
		return task, err
	}
	err = json.Unmarshal(body, &task)
	return task, err
}

func (client *Client) GetTaskByName(name string) (*Task, error) {
	var tasks, err = client.ListTasks("name=" + name)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if task.Name == name {
			return &task, nil
		}
	}
	return nil, errors.New("Task not found")
}

func (task *Task) ForceComplete(client *Client) error {
	if task.ID == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.request("PUT", "task/"+task.ID+"/forcecomplete", nil)
	return err
}

func (task *Task) WaitForTask(client *Client, printProgress bool) (*Task, error) {
	var progress int
	if task.State == "completed" || task.State == "failed" {
		return task, nil
	}
	feed, err := client.GetChangeFeed("task", map[string]string{"id": task.ID})
	var newVal Task
	if err != nil {
		return nil, err
	}
	for {
		select {
		case msg := <-feed.Data:
			if msg.Error != nil {
				fmt.Println(msg.Error)
				feed.Close()
				continue
			}
			err = json.Unmarshal(msg.NewValue, &newVal)
			if err != nil {
				feed.Close()
				return nil, err
			}
			if printProgress && newVal.Progress != progress {
				progress = newVal.Progress
				fmt.Println(newVal.Progress)
			}
			if newVal.State == "completed" {
				feed.Close()
				return &newVal, nil
			} else if newVal.State == "failed" {
				feed.Close()
				return &newVal, fmt.Errorf("Task Failed: %s", newVal.Message)
			}
		case <-feed.Done:
			return nil, fmt.Errorf("ChangeFeed was closed")
		}
	}
}
