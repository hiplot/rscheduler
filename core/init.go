package core

var RScheduler rScheduler // global scheduler

func Init() {
	InitMonitor()
	RScheduler = rScheduler{M: make(map[string]*ProcList)}
	RScheduler.Start()
}
