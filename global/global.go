package global

import (
	"log"
	"rscheduler/rslog"
)

var Logger *rslog.RsLogger // global logger
var GCChan chan struct{}   // 用于传递GC信号

func Init() {
	Logger = rslog.NewGlobalLogger()
	GCChan = make(chan struct{})
	log.Println("全局变量初始化成功")
}
