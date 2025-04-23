package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	Pending TaskStatus = "peding"
	Running TaskStatus = "running"
	Done    TaskStatus = "done"
	Failed  TaskStatus = "failed"
)

func (tm *TaskManager) process(id string) {
	tm.mu.Lock()
	task := tm.tasks[id]
	task.Status = Running
	tm.mu.Unlock()

	tm.mu.Lock()
	task.Result = map[string]string{"message": "done"}
	task.Status = Done
	tm.mu.Unlock()
}

func (tm *TaskManager) worker() {
	for id := range tm.queue {
		tm.process(id)
	}
}

func NewTaskManager(queueSize int) *TaskManager {
	tm := &TaskManager{
		tasks: make(map[string]*Task),
		queue: make(chan string, queueSize),
	}

	for i := 0; i < 5; i++ {
		go tm.worker()
	}
	return tm
}

func (tm *TaskManager) Submit(ctx context.Context, payload interface{}) (string, error) {
	id := uuid.NewString()
	task := &Task{ID: id, Status: Pending, CreatedAt: time.Now()}
	tm.mu.Lock()
	tm.tasks[id] = task
	tm.mu.Unlock()

	select {
	case tm.queue <- id:
		return id, nil
	default:
		return "", errors.New("Очередь заполнена")
	}
}

func (tm *TaskManager) GetID(id string) (*Task, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	t, ok := tm.tasks[id]
	return t, ok
}
