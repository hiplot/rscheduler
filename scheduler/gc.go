package scheduler

import (
	"container/list"
	"go.uber.org/zap"
	"rscheduler/global"
	"rscheduler/processor"
	"time"
)

var gcTime time.Time

func enableGC() {
	go func() {
		for {
			select {
			case <-global.GCChan:
				if !allowGC() {
					continue
				}
				gc()
			}
		}
	}()
}

func gc() {
	global.Logger.Infow("Start Scheduler GC")
	RScheduler.Lock.Lock()
	defer RScheduler.Lock.Unlock()

	startTime := time.Now()
	for _, procList := range RScheduler.M {
		if procList == nil {
			continue
		}
		removeList := make([]*list.Element, 0)
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if passGCCheck(proc) {
				continue
			}
			_ = proc.ForceClose()
			removeList = append(removeList, i)
		}
		for _, i := range removeList {
			procList.Remove(i)
		}
	}
	global.Logger.Infow("Scheduler GC complete", zap.Duration("useTime", time.Now().Sub(startTime)))
}

func passGCCheck(proc *processor.Proc) bool {
	// 检查进程是否正在运行
	if proc == nil || !proc.IsIdle() {
		return true
	}
	return false
}

// 避免频繁GC
func allowGC() bool {
	if time.Now().Sub(gcTime) > 1*time.Minute {
		gcTime = time.Now()
		return true
	}
	return false
}
