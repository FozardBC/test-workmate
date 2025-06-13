package taskprocessor

import (
	"test-app/internal/domain/models"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func SomeProcessingTask(task *models.Task) {
	task.Status.Status = models.StatusProcessing

	time.Sleep(time.Second * 10)

	task.Status = models.TaskStatus{Status: models.StatusCompleted}

	task.Finished_at = time.Now()

	task.Result = gofakeit.Street() // any random result

	task.Duration = time.Duration(task.Finished_at.Sub(task.Created_at).Seconds())

}
