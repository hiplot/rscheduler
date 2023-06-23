package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liang09255/lutils/conv"
	"net/http"
	"rscheduler/config"
	"rscheduler/core"
	"rscheduler/global"
	"rscheduler/mq"
)

func main() {
	config.Init()
	global.Init()
	mq.Init()
	core.Init()

	g := gin.Default()
	g.GET("/completed", TaskCompleteHandler)
	_ = g.Run(":8080")
}

const (
	FailedResponseCode  = 0
	SuccessResponseCode = 1
)

func TaskCompleteHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	taskID := c.Query("taskID")
	kill := conv.ToBool(c.Query("kill"))
	core.RScheduler.TaskComplete(taskName, taskID, kill)
	c.JSON(http.StatusOK, "Success")
}
