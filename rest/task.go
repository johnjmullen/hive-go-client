package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Task describes a task record from the rest api
type Task struct {
	Cancellable     bool      `json:"cancellable"`
	Description     string    `json:"description"`
	FinishedTime    time.Time `json:"finishedTime"`
	ID              string    `json:"id"`
	LastUpdatedTime time.Time `json:"lastUpdatedTime"`
	Message         string    `json:"message"`
	Name            string    `json:"name"`
	Progress        float32   `json:"progress"`
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

// ListTasks returns an array of tasks with an optional filter string
func (client *Client) ListTasks(query string) ([]Task, error) {
	var tasks []Task
	path := "tasks"
	if query != "" {
		path += "?" + query
	}
	body, err := client.request("GET", path, nil)
	if err != nil {
		return tasks, err
	}
	err = json.Unmarshal(body, &tasks)
	return tasks, err
}

// GetTask request a task by id
func (client *Client) GetTask(id string) (*Task, error) {
	var task Task
	if id == "" {
		return &task, errors.New("id cannot be empty")
	}
	body, err := client.request("GET", "task/"+id, nil)
	if err != nil {
		return &task, err
	}
	err = json.Unmarshal(body, &task)
	return &task, err
}

// GetTaskByName request a task by name
func (client *Client) GetTaskByName(name string) (*Task, error) {
	var tasks, err = client.ListTasks("name=" + url.QueryEscape(name))
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

// ForceComplete marks a task as completed in the database
func (task *Task) ForceComplete(client *Client) error {
	if task.ID == "" {
		return errors.New("id cannot be empty")
	}
	_, err := client.request("PUT", "task/"+task.ID+"/forcecomplete", nil)
	return err
}

// WatchTask monitors a task changefeed and sends updates to taskData
func (task Task) WatchTask(client *Client, taskData chan Task, errorChannel chan error) {
	if task.State == "completed" || task.State == "failed" {
		taskData <- task
		return
	}
	newVal := Task{}
	feed, err := client.GetChangeFeed("task", map[string]string{"id": task.ID}, true)
	if err != nil {
		errorChannel <- err
		return
	}
	defer feed.Close()
	for {
		select {
		case msg := <-feed.Data:
			if msg.Error != nil {
				errorChannel <- msg.Error
				return
			}
			err = json.Unmarshal(msg.NewValue, &newVal)
			if err != nil {
				errorChannel <- err
				return
			}
			taskData <- newVal
			if newVal.State == "completed" || newVal.State == "failed" {
				return
			}
		case <-feed.Done:
			return
		}
	}
}

//WaitForTask blocks until a task is complete and returns the task
func (task Task) WaitForTask(client *Client, printProgress bool) (*Task, error) {
	var progress float32
	newVal := task
	errChannel := make(chan error)
	taskData := make(chan Task)
	go task.WatchTask(client, taskData, errChannel)
	for {
		select {
		case newVal = <-taskData:
			if printProgress && newVal.Progress != progress {
				progress = newVal.Progress
				fmt.Println(newVal.Progress)
			}
			if newVal.State == "completed" || newVal.State == "failed" {
				return &newVal, nil
			}
		case err := <-errChannel:
			return &newVal, err
		}
	}
}
