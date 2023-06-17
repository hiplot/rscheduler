package main

import (
	"go.uber.org/zap"
	"sync"
)

type procMap struct {
	lock   sync.RWMutex
	m      map[string]*ProcList
	logger *rsLogger
}

var ProcMap procMap

func init() {
	ProcMap = procMap{
		m:      make(map[string]*ProcList),
		logger: newGlobalLogger(),
	}
}

func (pm *procMap) AddTask(t *Task) {
	proc := pm.getProc(t)
	if proc == nil {
		panic("get proc failed")
	}
	var err error
	_, err = proc.Exec(`taskID = "%s"`, t.id)
	_, err = proc.Exec(`source("./rscript/%s.R")`, t.name)
	if err != nil {
		pm.logger.Errorf("Exec failed: ")
		return
	}
}

func (pm *procMap) TaskComplete(taskName, taskID string) {
	// TODO collect result
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pList := pm.m[taskName]
	for i := pList.Back(); i != nil; i = i.Prev() {
		proc := i.Value.(*Proc)
		if proc.task != nil && proc.task.id == taskID {
			pm.logger.Infow("Task complete success", zap.String("taskName", taskName), zap.String("taskID", taskID))
			proc.task = nil
			return
		}
	}
}

func (pm *procMap) getProc(t *Task) *Proc {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pList := pm.m[t.name]
	if pList == nil || pList.Len() == 0 {
		proc := pm.makeNewProc(t.name)
		proc.task = t
		return proc
	}

	for procElement := pList.Front(); procElement != nil; procElement = procElement.Next() {
		proc := procElement.Value.(*Proc)
		if proc != nil && proc.task == nil {
			pm.logger.Infow("Find a idle processor", zap.String("taskName", t.name), zap.String("taskID", t.id))
			proc.task = t
			// put the procElement to the end of list
			pList.MoveToBack(procElement)
			return proc
		}
	}

	// TODO limit processor count
	proc := pm.makeNewProc(t.name)
	proc.task = t
	return proc
}

// This func is in order to reduce lock granularity
func (pm *procMap) makeNewProc(name string) *Proc {
	pm.lock.Unlock()
	defer pm.lock.Lock()
	return newProc(name)
}
