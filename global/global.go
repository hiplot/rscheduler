package global

import (
	"rscheduler/rslog"
)

var Logger *rslog.RsLogger // global logger

func Init() {
	Logger = rslog.NewGlobalLogger()
}
