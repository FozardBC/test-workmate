package delete

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"test-app/internal/lib/api/response"

	"test-app/internal/api/middlewares/requestid"

	"github.com/gin-gonic/gin"
)

type Deleter interface {
	DeleteTask(ctx context.Context, id int64) error
}

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
			logHandler.Error("failed to delete task", "err", err.Error(), "id", id)

			c.JSON(http.StatusInternalServerError, response.Error("Internal error"))

			return
		}

		c.JSON(http.StatusOK, response.OK())

		logHandler.Info("Task deleted", "id", id)

	}
}
