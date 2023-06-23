package core

var RScheduler rScheduler // global scheduler

var CPUPercent int // CPU占用率
var MemPercent int // 内存占用率

func Init() {
	RScheduler = rScheduler{M: make(map[string]*ProcList)}
}
