package scheduler

import (
	"log"
	"rscheduler/processor"
)

var RScheduler rScheduler // global scheduler

func Init() {
	RScheduler = rScheduler{M: make(map[string]*processor.ProcList)}
	RScheduler.Start()
	enableGC()
	startHealthCheck()
	log.Println("调度器初始化成功")
}
