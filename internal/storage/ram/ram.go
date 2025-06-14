package ram

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"test-app/internal/domain/models"
	"test-app/internal/storage"
	"time"
)

var (
	taskIDKey = "taskID"

	ErrContextCancelledBef = errors.New("context cancelled before lock mutex")
	ErrContextCancelled    = errors.New("context cancelled")
)

type MemStorage struct {
	Tasks     map[int64]*models.Task
	Log       *slog.Logger
	CurrentId int64
	Mut       sync.RWMutex
}

func New(log *slog.Logger) *MemStorage {
	return &MemStorage{
		Tasks: map[int64]*models.Task{},
		Log:   log,
	}
}

func (s *MemStorage) Save(ctx context.Context, task *models.Task) (int64, error) {

	var id int64

	s.Mut.Lock()
	defer s.Mut.Unlock()

	select {
	case <-ctx.Done():
		s.Log.Debug(ErrContextCancelled.Error(), "NameTask", task.Name)
		return 0, ctx.Err()
	default:
	}

	id = s.CurrentId

	s.Tasks[id] = task

	s.CurrentId++

	return id, nil

}

func (s *MemStorage) Delete(ctx context.Context, id int64) error {

	if err := ctx.Err(); err != nil {
		s.Log.Debug(ErrContextCancelledBef.Error(), taskIDKey, id)

		return fmt.Errorf("%w:%w", ErrContextCancelledBef, err)
	}

	s.Mut.RLock()
	defer s.Mut.RUnlock()

	select {
	case <-ctx.Done():
		s.Log.Debug(ErrContextCancelledBef.Error(), taskIDKey, id)

		return fmt.Errorf("%w:%w", ErrContextCancelled, ctx.Err())
	default:
	}

	_, exists := s.Tasks[id]
	if !exists {
		s.Log.Debug(storage.ErrTaskNotFound.Error(), taskIDKey, id)

		return storage.ErrTaskNotFound
	}

	delete(s.Tasks, id)

	return nil
}

func (s *MemStorage) Task(ctx context.Context, id int64) (*models.Task, error) {
	if err := ctx.Err(); err != nil {
		s.Log.Debug(ErrContextCancelledBef.Error(), taskIDKey, id)

		return nil, fmt.Errorf("%w:%w", ErrContextCancelledBef, err)
	}

	s.Mut.RLock()
	defer s.Mut.RUnlock()

	select {
	case <-ctx.Done():
		s.Log.Debug(ErrContextCancelledBef.Error(), taskIDKey, id)

		return nil, fmt.Errorf("%w:%w", ErrContextCancelled, ctx.Err())
	default:
	}

	task, exists := s.Tasks[id]
	if !exists {
		s.Log.Error(storage.ErrTaskNotFound.Error(), taskIDKey, id)

		return nil, storage.ErrTaskNotFound
	}

	taskCopy := *task

	return &taskCopy, nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	s.Log.Debug("storage is active")

	time.Sleep(10 * time.Minute)
	return nil
}

func (s *MemStorage) Close() {
	s.Log.Info("Storage is closed")
}
