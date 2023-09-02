package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liang09255/lutils/conv"
	"net/http"
	"rscheduler/config"
	"rscheduler/global"
	"rscheduler/monitor"
	"rscheduler/mq"
	"rscheduler/scheduler"
)

func main() {
	config.Init()         // 初始化配置
	global.Init()         // 初始化全局变量
	mq.Init()             // 初始化消息队列
	scheduler.Init()      // 初始化调度器
	monitor.InitMonitor() // 初始化监控

	g := gin.Default()
	g.GET("/completed", TaskCompleteHandler)
	_ = g.Run(":8080")
}

func TaskCompleteHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	taskID := c.Query("taskID")
	kill := conv.ToBool(c.Query("kill"))
	scheduler.RScheduler.TaskComplete(taskName, taskID, kill)
	c.JSON(http.StatusOK, "Success")
}
