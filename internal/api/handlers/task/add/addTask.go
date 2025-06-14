package addTask

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"test-app/internal/api/types"
	"test-app/internal/lib/api/response"

	"test-app/internal/api/middlewares/requestid"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	TaskName string `json:"task_name" validate:"required,min=2,max=50"`
}

type Payload struct {
	ID int64 `json:"id"`
}

type TaskSaver interface {
	SaveTask(ctx context.Context, name string) (int64, error)
}

// addTask godoc
// @Summary		Add the new task
// @Description	add task with custom name
// @Tags		tasks
// @Accept		json
// @Produce		json
// @Param		task	body		Request				true	"Task create request"
// @Success		201		{object}	response.Response
// @Header 		201		{string}	Location			"URL of the created task (e.g., /tasks/123)"
// @Header 		201		{string}	x-request-id		"request-id"
// @Failure		400		{object}	response.Response
// @Failure		500		{object}	response.Response
// @Router		/tasks	[post]
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

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)

			logHandler.Error("invalid request", "err", err.Error())

			c.JSON(http.StatusBadRequest, response.ValidationError(validatorErr))

			return
		}

		id, err := saver.SaveTask(ctx, req.TaskName)

		if err != nil {
			logHandler.Error("failed to save task", "err", err.Error())

			c.JSON(http.StatusInternalServerError, response.Error("Internal error"))
		}

		c.Header("Location", fmt.Sprintf("/tasks/%d", id))

		c.JSON(http.StatusCreated, response.OKWithPayload(Payload{ID: id}))

		logHandler.Info("Task saved", "name", req.TaskName, "id", id)

	}
}
