package api

import (
	"github.com/gin-gonic/gin"
	"github.com/liang09255/lutils/conv"
	"net/http"
	"rscheduler/scheduler"
)

type taskAPI struct{}

var TaskAPI = &taskAPI{}

func (t taskAPI) TaskCompleteHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	taskID := c.Query("taskID")
	kill := conv.ToBool(c.Query("kill"))
	scheduler.RScheduler.TaskComplete(taskName, taskID, kill)
	c.JSON(http.StatusOK, "Success")
}
