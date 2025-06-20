package getTask

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"test-app/internal/domain/models"
	"test-app/internal/lib/api/response"
	"test-app/internal/storage"

	"test-app/internal/api/middlewares/requestid"

	"github.com/gin-gonic/gin"
)

type TaskGetter interface {
	Task(ctx context.Context, id int64) (*models.Task, error)
}

// New godoc
// @Summary      Get a task by ID
// @Description  Retrieve a task with the specified ID. Returns the task details if found.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "Task ID"
// @Success      200  {object}  response.Response{payload=models.Task}  "Task details"
// @Failure      400  {object}  response.Response{status=string,error=string}  "Invalid task ID"
// @Failure      404  {object}  response.Response{error=string}  "Task not found"
// @Failure      500  {object}  response.Response{error=string}  "Internal server error"
// @Router       /tasks/{id} [get]
func New(log *slog.Logger, taskGetter TaskGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		logHandler := log.With(
			slog.String("requestID", requestid.Get(c)),
		)

		idParam := c.Param("id")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logHandler.Error("failed to convertation url parameter", "paramName", "id", "param", idParam)

			c.JSON(http.StatusBadRequest, response.Error(fmt.Sprintf("Invalid parameter:%s", idParam)))

			return
		}

		task, err := taskGetter.Task(ctx, int64(id))
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				logHandler.Debug("Task is not exists", "id", id)

				c.JSON(http.StatusNotFound, response.Error("Task is not exists"))

				return
			}
			logHandler.Error("failed to get task", "err", err.Error(), "id", id)

			c.JSON(http.StatusInternalServerError, response.Error("Internal error"))

			return
		}

		c.JSON(http.StatusOK, response.OKWithPayload(task))

		logHandler.Info("Task is recived", "name", task.Name, "id", id)
	}
}
