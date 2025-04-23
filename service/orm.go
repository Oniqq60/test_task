package service

import (
	"sync"
	"time"
)

type Task struct {
	ID        string
	Status    TaskStatus
	Result    interface{}
	Error     string
	CreatedAt time.Time
}

type TaskManager struct {
	mu    sync.Mutex
	tasks map[string]*Task
	queue chan string
}
