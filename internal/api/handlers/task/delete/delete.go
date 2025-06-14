package delete

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"test-app/internal/lib/api/response"
	"test-app/internal/storage"

	"test-app/internal/api/middlewares/requestid"

	"github.com/gin-gonic/gin"
)

type Deleter interface {
	DeleteTask(ctx context.Context, id int64) error
}

// New godoc
// @Summary      Delete a task by ID
// @Description  Delete a task with the specified ID. Returns no content if successful. If the task is not found, returns an error message.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "Task ID" example:"123"
// @Success      204
// @Failure      400  {object}  response.Response{error=string}  "Invalid task ID"
// @Failure      404  {object}  response.Response{error=string}  "Task not found"
// @Failure      500  {object}  response.Response{error=string}  "Internal server error"
// @Router       /tasks/{id} [delete]
func New(log *slog.Logger, deleter Deleter) gin.HandlerFunc {
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

		err = deleter.DeleteTask(ctx, int64(id))
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				logHandler.Error("task not found")

				c.JSON(http.StatusNotFound, response.Error("Task not exists"))
			}

			logHandler.Error("failed to delete task", "err", err.Error(), "id", id)

			c.JSON(http.StatusInternalServerError, response.Error("Internal error"))

			return
		}

		c.Status(http.StatusNoContent)

		logHandler.Info("Task deleted", "id", id)

	}
}
