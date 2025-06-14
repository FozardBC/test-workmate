package api

import (
	"log/slog"

	_ "test-app/docs"

	addTask "test-app/internal/api/handlers/task/add"
	"test-app/internal/api/handlers/task/delete"
	getTask "test-app/internal/api/handlers/task/get"
	"test-app/internal/lib/api/log"
	taskManager "test-app/internal/services/task-manager"

	"test-app/internal/api/middlewares/requestid"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	Router  *gin.Engine
	Log     *slog.Logger
	Service *taskManager.Manager
}

func New(log *slog.Logger, service *taskManager.Manager) *API {
	return &API{
		Router:  gin.New(),
		Log:     log,
		Service: service,
	}
}

func (api *API) Setup() {
	v1 := api.Router.Group("api/v1/")

	v1.Use(requestid.RequestIdMidlleware())
	v1.Use(gin.LoggerWithFormatter(log.Logging))

	v1.POST("/tasks", addTask.New(api.Log, api.Service))
	v1.GET("/tasks/:id", getTask.New(api.Log, api.Service))
	v1.DELETE("tasks/:id", delete.New(api.Log, api.Service))

	v1.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))

}
