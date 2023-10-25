package main

import (
	"rscheduler/api"
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
	api.Start()           // API服务 阻塞
}
