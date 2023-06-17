package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liang09255/lutils/conv"
	"net/http"
)

func main() {
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
	ProcMap.TaskComplete(taskName, taskID, kill)
	c.JSON(http.StatusOK, "Success")
}

func NewTaskHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	task := NewTask(taskName)

	if !enableNewTask() {
		c.JSON(http.StatusOK, gin.H{
			"status": FailedResponseCode,
			"msg":    "limit",
		})
		return
	}

	err := ProcMap.AddTask(task)
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
