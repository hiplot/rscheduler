package main

import (
	"log"
	"sync"
)

type procMap struct {
	lock sync.RWMutex
	m    map[string][]*Processor
}

var ProcMap procMap

func init() {
	ProcMap = procMap{
		m: make(map[string][]*Processor),
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
	for _, proc := range p.m[taskName] {
		if proc != nil && proc.task != nil && proc.task.id == taskID {
			log.Println("Task complete success")
			proc.task = nil
			return
		}
	}
}

func (p *procMap) getProc(t *Task) *Processor {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.m[t.name]) == 0 {
		proc := p.makeNewProc(t.name)
		proc.task = t
		return proc
	}

	for _, proc := range p.m[t.name] {
		if proc != nil && proc.task == nil {
			log.Println("Find a idle processor")
			proc.task = t
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
