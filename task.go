package coprhd

import (
	"errors"
	"fmt"
	"time"
)

const (
	GetTaskUriTpl = "vdc/tasks/%s.json"

	TaskPollDelay              = time.Millisecond * 250
	TaskStatePending TaskState = "pending"
	TaskStateError   TaskState = "error"
	TaskStateReady   TaskState = "ready"
)

var (
	ErrTaskWaitTimeout = errors.New("WaitDone timeout")
)

type (
	TaskService struct {
		*Client
	}

	Task struct {
		Name        string    `json:"name"`
		Id          string    `json:"id"`
		State       TaskState `json:"state"`
		Message     string    `json:"message"`
		Description string    `json:"description"`
		Progress    int       `json:"progress"`
		Resource    struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"resource"`
	}

	TaskState string
)

func (this *Client) Task() *TaskService {
	return &TaskService{this}
}

// Query returns the task object
func (this *TaskService) Query(id string) (Task, error) {
	path := fmt.Sprintf(GetTaskUriTpl, id)
	task := Task{}

	err := this.Get(path, nil, &task)

	return task, err
}

// WaitDone does a busy poll to wait for a task to reach the specified state with the timeout
func (this *TaskService) WaitDone(id string, state TaskState, to time.Duration) error {
	end := time.Now().Add(to)

	for {
		task, err := this.Query(id)
		if err != nil {
			return err
		}

		if task.State == TaskStateError {
			return errors.New(task.Message + ":" + task.Description)
		}

		if task.State == state {
			break
		}

		if time.Now().After(end) {
			return ErrTaskWaitTimeout
		}
		time.Sleep(TaskPollDelay)
	}

	return nil
}
