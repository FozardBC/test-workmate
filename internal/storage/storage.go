package storage

import (
	"context"
	"errors"
	"test-app/internal/domain/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage interface {
	Save(ctx context.Context, task *models.Task) (int64, error)
	Delete(ctx context.Context, id int64) error
	Task(ctx context.Context, id int64) (*models.Task, error)
	Ping(ctx context.Context) error
	Close()
}
