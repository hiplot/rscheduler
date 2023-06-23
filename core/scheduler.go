package core

import (
	"fmt"
	"go.uber.org/zap"
	"rscheduler/global"
	"rscheduler/mq"
	"sync"
	"time"
)

type rScheduler struct {
	lock sync.RWMutex
	M    map[string]*ProcList
}

func (rs *rScheduler) Start() {
	go func() {
		for {
			if !EnableNewTask() {
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
			t := NewTask(taskDetail)
			if t == nil {
				continue
			}
			// add task to scheduler
			err = rs.AddTask(t)
			if err != nil {
				global.Logger.Error("add task failed, err:", err)
			}
		}
	}()
}

func (rs *rScheduler) AddTask(t *Task) error {
	// bind task with processor
	proc := rs.getProc(t)
	if proc == nil {
		global.Logger.Error("get a nil processor, please check")
		return fmt.Errorf("get a nil processor, please check")
	}
	// do task
	var err error
	_, err = proc.Exec(`taskID = "%s"`, t.ID)
	_, err = proc.Exec(`source("./rscript/%s.R")`, t.Name)
	if err != nil {
		global.Logger.Error("Exec failed, err: ", err)
		return err
	}
	return nil
}

func (rs *rScheduler) TaskComplete(taskName, taskID string, kill bool) {
	// TODO collect result
	rs.lock.Lock()
	defer rs.lock.Unlock()
	pList := rs.M[taskName]
	for i := pList.Back(); i != nil; i = i.Prev() {
		proc := i.Value.(*Proc)
		if proc.task != nil && proc.task.ID == taskID {
			global.Logger.Infow("Task complete success", zap.String("taskName", taskName), zap.String("taskID", taskID))
			proc.task = nil
			if kill {
				_ = proc.Close()
				pList.Remove(i)
			}
			return
		}
	}
}

// bind task with processor
func (rs *rScheduler) getProc(t *Task) *Proc {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	pList := rs.M[t.Name]

	// create new pList and processor
	if pList == nil || pList.Len() == 0 {
		proc := rs.makeNewProc(t.Name)
		proc.task = t
		return proc
	}

	// try to get an idle processor
	for procElement := pList.Front(); procElement != nil; procElement = procElement.Next() {
		proc := procElement.Value.(*Proc)
		if proc != nil && proc.task == nil {
			global.Logger.Infow("Find an idle processor", zap.String("taskName", t.Name), zap.String("taskID", t.ID))
			proc.task = t
			// put the procElement to the end of list
			pList.MoveToBack(procElement)
			return proc
		}
	}

	// can not find an idle processor
	// create a new processor
	proc := rs.makeNewProc(t.Name)
	proc.task = t
	return proc
}

// This func is in order to reduce lock granularity
func (rs *rScheduler) makeNewProc(name string) *Proc {
	rs.lock.Unlock()
	defer rs.lock.Lock()
	return newProc(name)
}
