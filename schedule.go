package main

import (
	"log"
	"sync"
)

type procMap struct {
	lock sync.RWMutex
	m    map[string]*ProcessorList
}

var ProcMap procMap

func init() {
	ProcMap = procMap{
		m: make(map[string]*ProcessorList),
	}
}

func (p *procMap) AddTask(t *Task) {
	proc := p.getProc(t)
	if proc == nil {
		panic("get proc failed")
	}
	var err error
	_, err = proc.Exec(`taskID = "%s"`, t.id)
	_, err = proc.Exec(`source("./rscript/%s.R")`, t.name)
	if err != nil {
		log.Println(err)
		return
	}
}

func (p *procMap) TaskComplete(taskName, taskID string) {
	// TODO collect result
	p.lock.Lock()
	defer p.lock.Unlock()
	pList := p.m[taskName]
	for i := pList.Back(); i != nil; i = i.Prev() {
		proc := i.Value.(*Processor)
		if proc.task != nil && proc.task.id == taskID {
			log.Println("Task complete success")
			proc.task = nil
			return
		}
	}
}

func (p *procMap) getProc(t *Task) *Processor {
	p.lock.Lock()
	defer p.lock.Unlock()
	pList := p.m[t.name]
	if pList == nil || pList.Len() == 0 {
		proc := p.makeNewProc(t.name)
		proc.task = t
		return proc
	}

	for procElement := pList.Front(); procElement != nil; procElement = procElement.Next() {
		proc := procElement.Value.(*Processor)
		if proc != nil && proc.task == nil {
			log.Println("Find a idle processor")
			proc.task = t
			// put the procElement to the end of list
			pList.MoveToBack(procElement)
			return proc
		}
	}

	// TODO limit processor count
	proc := p.makeNewProc(t.name)
	proc.task = t
	return proc
}

// This func is in order to reduce lock granularity
func (p *procMap) makeNewProc(name string) *Processor {
	p.lock.Unlock()
	defer p.lock.Lock()
	return newProcess(name)
}
