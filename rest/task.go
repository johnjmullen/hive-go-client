package rest

import (
	"encoding/json"
	"errors"
	"fmt"
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
		return &task, errors.New("Id cannot be empty")
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

// ForceComplete marks a task as completed in the database
func (task *Task) ForceComplete(client *Client) error {
	if task.ID == "" {
		return errors.New("Id cannot be empty")
	}
	_, err := client.request("PUT", "task/"+task.ID+"/forcecomplete", nil)
	return err
}

// WatchTask monitors a task changefeed and sends updates to taskData
func (task Task) WatchTask(client *Client, taskData chan Task, done chan struct{}) {
	defer close(done)
	if task.State == "completed" || task.State == "failed" {
		return
	}
	newVal := Task{}
	feed, err := client.GetChangeFeed("task", map[string]string{"id": task.ID})
	if err != nil {
		return
	}
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-timer.C:
			//work around race condition
			t, _ := client.GetTask(task.ID)
			if t.State == "completed" || t.State == "failed" {
				taskData <- *t
				feed.Close()
				continue
			}
		case msg := <-feed.Data:
			if msg.Error != nil {
				fmt.Println(msg.Error)
				feed.Close()
				continue
			}
			err = json.Unmarshal(msg.NewValue, &newVal)
			if err != nil {
				fmt.Println("Error with json unmarshal", err)
				feed.Close()
				return
			}
			taskData <- newVal
			if newVal.State == "completed" || newVal.State == "failed" {
				feed.Close()
			}
		case <-feed.Done:
			return
		}
	}
}

//WaitForTask blocks until a task is complete and returns the task
func (task Task) WaitForTask(client *Client, printProgress bool) *Task {
	var progress int
	newVal := task
	done := make(chan struct{})
	taskData := make(chan Task)
	go task.WatchTask(client, taskData, done)
	for {
		select {
		case newVal = <-taskData:
			if printProgress && newVal.Progress != progress {
				progress = newVal.Progress
				fmt.Println(newVal.Progress)
			}
		case <-done:
			if printProgress && newVal.Progress != progress {
				fmt.Println(newVal.Progress)
			}
			return &newVal
		}
	}
}
