package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liang09255/lutils/conv"
)

func main() {
	g := gin.Default()
	g.GET("/completed", TaskCompleteHandler)
	g.GET("/newTask", NewTaskHandler)
	_ = g.Run(":8080")
}

func TaskCompleteHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	taskID := c.Query("taskID")
	kill := conv.ToBool(c.Query("kill"))
	ProcMap.TaskComplete(taskName, taskID, kill)
	c.JSON(200, "Success")
	return
}

func NewTaskHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	task := NewTask(taskName)
	ProcMap.AddTask(task)
}
