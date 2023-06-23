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
	global.Init()
	config.Init()
	mq.RabbitMQInit()
	core.Init()

	g := gin.Default()
	g.GET("/completed", TaskCompleteHandler)
	g.GET("/newTask", NewTaskHandler)
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

func NewTaskHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	task := core.NewTask(taskName)

	if !core.EnableNewTask() {
		c.JSON(http.StatusOK, gin.H{
			"status": FailedResponseCode,
			"msg":    "limit",
		})
		return
	}

	err := core.RScheduler.AddTask(task)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": FailedResponseCode,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": SuccessResponseCode,
		"msg":    "",
	})
}
