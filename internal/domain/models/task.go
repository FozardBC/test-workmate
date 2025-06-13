package models

import "time"

type Task struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Status      TaskStatus    `json:"status"`
	Created_at  time.Time     `json:"created_at"`
	Finished_at time.Time     `json:"finished_at,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
	Result      interface{}   `json:"result,omitempty"`
}

type TaskStatus struct {
	Status string `json:"status"`
	Err    error  `json:"err,omitempty"`
}

var (
	StatusCreated    = "created"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)
