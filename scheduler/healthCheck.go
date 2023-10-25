package scheduler

import (
	"container/list"
	"rscheduler/processor"
	"time"
)

func startHealthCheck() {
	go func() {
		for {
			healthCheck()
			time.Sleep(time.Second * 5)
		}
	}()
}

func healthCheck() {
	RScheduler.Lock.Lock()
	defer RScheduler.Lock.Unlock()
	for _, procList := range RScheduler.M {
		if procList == nil {
			continue
		}
		removeList := make([]*list.Element, 0)
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if passHealthCheck(proc) {
				continue
			}
			_ = proc.ForceClose()
			removeList = append(removeList, i)
		}
		for _, i := range removeList {
			procList.Remove(i)
		}
	}
}

func passHealthCheck(proc *processor.Proc) bool {
	// 检查进程是否异常退出，内存是否超限
	if !proc.HealthCheck() || !proc.MemCheck() {
		return false
	}
	// 检查繁忙进程是否超时
	if !proc.IsIdle() && proc.Task.IsTimeout() {
		return false
	}
	return true
}
