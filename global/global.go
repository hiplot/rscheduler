package global

import (
	"log"
	"rscheduler/rslog"
)

const VERSION = "0.1.0"

var Logger *rslog.RsLogger // global logger
var GCChan chan struct{}   // 用于传递GC信号

func Init() {
	Logger = rslog.NewGlobalLogger()
	GCChan = make(chan struct{})
	log.Println("全局变量初始化成功")
}
