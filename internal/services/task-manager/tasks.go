package taskManager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"test-app/internal/domain/models"
	taskprocessor "test-app/internal/services/task-processor"
	"test-app/internal/storage"
	"time"
)

type Manager struct {
	Storage storage.Storage
	Log     *slog.Logger
}

var (
	ErrDeleteTask = errors.New("failed to delete task")
)

func New(storage storage.Storage, log *slog.Logger) *Manager {
	return &Manager{
		Storage: storage,
		Log:     log,
	}
}

func (m *Manager) SaveTask(ctx context.Context, name string) (int64, error) {

	m.Log.Debug("started to save task", "name", name)

	Task := &models.Task{}

	Task.Name = name
	Task.Status = models.TaskStatus{Status: models.StatusCreated}
	Task.Created_at = time.Now()

	id, err := m.Storage.Save(ctx, Task)
	if err != nil {
		m.Log.Error("failed to save task", "err", err.Error())

		return 0, fmt.Errorf("failed to save task:%w", err)
	}

	go taskprocessor.SomeProcessingTask(Task)

	m.Log.Debug("task is saved. Starting to process", "taskID", Task.ID)

	return id, nil

}

func (m *Manager) Task(ctx context.Context, id int64) (*models.Task, error) {
	m.Log.Debug("getting task by id", "id", id)

	task, err := m.Storage.Task(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get task from storage:%w", err)
	}

	m.Log.Debug("task is recived")

	return task, nil
}

func (m *Manager) DeleteTask(ctx context.Context, id int64) error {
	m.Log.Debug("starting to delete task", "id", id)

	if err := m.Storage.Delete(ctx, id); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			m.Log.Debug(ErrDeleteTask.Error(), "err", err.Error())

			return fmt.Errorf("%w:%w", ErrDeleteTask, err)
		}

		m.Log.Error(ErrDeleteTask.Error(), "id", id, "error", err)

		return fmt.Errorf("%w:%w", ErrDeleteTask, err)
	}

	m.Log.Debug("task deleted")

	return nil
}
