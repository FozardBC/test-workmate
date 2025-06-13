package addTask

import (
	"context"
	"log/slog"
	"net/http"
	"test-app/internal/api/types"
	"test-app/internal/lib/api/response"

	"test-app/internal/api/middlewares/requestid"

	"github.com/gin-gonic/gin"
)

type Request struct {
	TaskName string `json:"task_name"`
}

type TaskSaver interface {
	SaveTask(ctx context.Context, name string) (int64, error)
}

func New(log *slog.Logger, saver TaskSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		logHandler := log.With(
			slog.String("requestID", requestid.Get(c)),
		)

		var req Request

		if err := c.BindJSON(&req); err != nil {
			logHandler.Error(types.ErrDecodeReqBody.Error(), "err", err.Error())

			c.JSON(http.StatusBadRequest, response.Error(types.ErrDecodeReqBody.Error()))
			return
		}

		id, err := saver.SaveTask(ctx, req.TaskName)

		if err != nil {
			logHandler.Error("failed to save task", "err", err.Error())

			c.JSON(http.StatusInternalServerError, response.Error("Internal error"))
		}

		c.JSON(http.StatusCreated, response.OKWithPayload(map[string]int64{"id": id}))

		logHandler.Info("Task saved", "name", req.TaskName, "id", id)

	}
}
