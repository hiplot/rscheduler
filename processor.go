package main

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Processor struct {
	name   string
	cmd    *exec.Cmd
	inPipe *io.WriteCloser
	task   *Task
}

type ProcessorList struct {
	list.List
}

func newProcessorList() *ProcessorList {
	return &ProcessorList{list.List{}}
}

// newProcess
// Build a new r session
// Auto load nameFunc.R
// In this file, we can write some init code
func newProcess(name string) *Processor {
	log.Println("Create new process")
	cmd := exec.Command("R", "--vanilla")
	// TODO redirect log to file
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Println("get stdinPipe failed, err:", err)
		return nil
	}
	err = cmd.Start()
	if err != nil {
		log.Println("start cmd failed, err:", err)
		return nil
	}

	proc := &Processor{
		name:   name,
		cmd:    cmd,
		inPipe: &stdinPipe,
	}

	_, err = proc.Exec(`source("./rscript/%sInit.R")`, name)
	if err != nil {
		log.Println(err)
	}

	ProcMap.lock.Lock()
	defer ProcMap.lock.Unlock()

	if ProcMap.m[name] == nil {
		ProcMap.m[name] = newProcessorList()
	}

	ProcMap.m[name].PushBack(proc)
	return proc
}

// Exec
// This func will automatically add \n at the end
// Please ensure exec one line
func (p *Processor) Exec(format string, a ...any) (int, error) {
	if p.inPipe == nil {
		return 0, fmt.Errorf("inPipe is nil, exec failed")
	}
	return fmt.Fprintf(*p.inPipe, format+"\n", a...)
}
