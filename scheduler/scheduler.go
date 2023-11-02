package scheduler

import (
	"fmt"
	"rscheduler/processor"
	"sync"
	"time"

	"rscheduler/config"
	"rscheduler/global"
	"rscheduler/monitor"
	"rscheduler/mq"
	"rscheduler/task"

	"go.uber.org/zap"
)

type rScheduler struct {
	Lock        sync.RWMutex
	M           map[string]*processor.ProcList
	BusyProcNum int
	IdleProcNum int
}

func (rs *rScheduler) Start() {
	go func() {
		for {
			// Check System Status and Custom Rules
			if !monitor.EnableNewTask() || !rs.enableBusyProc() {
				time.Sleep(time.Second) // prevent dead loop
				continue
			}
			// get a new task from rabbitmq
			taskDetail, err := mq.RabbitMQ.Get()
			if err != nil {
				global.Logger.Error("get delivery failed, err:", err)
				time.Sleep(time.Second) // prevent dead loop
				continue
			}
			// decode and pack task
			t := task.NewHiPlotTask(taskDetail)
			if t == nil {
				continue
			}
			// add task to scheduler
			err = rs.RunTask(t)
			if err != nil {
				global.Logger.Error("run task failed, err:", err)
			}
		}
	}()
}

func (rs *rScheduler) RunTask(t *task.Task) error {
	// get a processor
	proc := rs.getProc(t.Name)
	if proc == nil {
		return fmt.Errorf("get a nil processor, please check")
	}
	// bind task with processor
	proc.BindTask(t)
	// run task
	return proc.Start()
}

// TaskComplete This method will be called when the task is completed
func (rs *rScheduler) TaskComplete(taskName, taskID string, kill bool) {
	rs.Lock.Lock()
	defer func() {
		// update processor num
		if !kill {
			rs.IdleProcNum++
		}
		rs.BusyProcNum--
		// release Lock
		rs.Lock.Unlock()
	}()

	pList := rs.M[taskName]
	for i := pList.Back(); i != nil; i = i.Prev() {
		proc := i.Value.(*processor.Proc)
		if proc.Task != nil && proc.Task.ID == taskID {
			proc.Complete()
			// judge whether to kill the processor
			if !rs.enableIdleProc() || !proc.MemCheck() || proc.PreDelete {
				kill = true
			}
			if kill {
				_ = proc.ForceClose()
				pList.Remove(i)
			}
			return
		}
	}
}

// bind task with processor
func (rs *rScheduler) getProc(name string) (proc *processor.Proc) {
	rs.Lock.Lock()
	defer func() {
		if proc == nil {
			global.Logger.Errorw("Get processor failed, processor is nil 1111", zap.String("taskName", name))
			return
		}
		proc.SetRun()
		rs.BusyProcNum++
		rs.Lock.Unlock()
	}()

	pList := rs.M[name]

	// create new pList and processor
	if pList == nil || pList.Len() == 0 {
		proc = rs.makeNewProc(name)
		return proc
	}

	// try to get an idle processor
	for procElement := pList.Front(); procElement != nil; procElement = procElement.Next() {
		proc = procElement.Value.(*processor.Proc)
		if proc != nil && proc.IsIdle() && proc.HealthCheck() {
			global.Logger.Infow("Find an idle processor", zap.String("taskName", name))
			rs.IdleProcNum--
			// put the procElement to the end of list
			pList.MoveToBack(procElement)
			return proc
		}
	}

	// can not find an idle processor
	// create a new processor
	proc = rs.makeNewProc(name)
	if proc == nil {
		global.Logger.Errorw("Create a new processor failed, proc is nil", zap.String("taskName", name))
		return nil
	}
	return proc
}

// create a new processor
func (rs *rScheduler) makeNewProc(name string) *processor.Proc {
	rs.Lock.Unlock()
	proc := processor.NewProc(name)
	// push processor into scheduler
	rs.addNewProc(proc)
	return proc
}

// push processor into scheduler
func (rs *rScheduler) addNewProc(proc *processor.Proc) {
	rs.Lock.Lock()
	if rs.M[proc.Name] == nil {
		rs.M[proc.Name] = processor.NewProcList()
	}
	rs.M[proc.Name].PushBack(proc)
}

// 是否允许创建新的空闲处理器
func (rs *rScheduler) enableIdleProc() bool {
	if rs.IdleProcNum >= config.Config.TaskLimit.MaxIdleProcessor {
		global.Logger.Infow("Arrive max idle processor num")
		return false
	}
	return true
}

// 是否允许创建新的繁忙处理器
func (rs *rScheduler) enableBusyProc() bool {
	if rs.BusyProcNum >= config.Config.TaskLimit.MaxBusyProcessor {
		global.Logger.Infow("Arrive max busy processor num")
		return false
	}
	return true
}
