package main

import (
	"container/list"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"io"
	"log"
	"os/exec"
)

type Processor struct {
	id     string
	name   string
	cmd    *exec.Cmd
	inPipe *io.WriteCloser
	task   *Task
	logger *rsLogger
}

type ProcessorList struct {
	list.List
}

func newProcessorList() *ProcessorList {
	return &ProcessorList{list.List{}}
}

// newProcessor
// Build a new r session
// Auto load nameFunc.R
// In this file, we can write some init code
func newProcessor(name string) *Processor {
	log.Println("Create new process")
	cmd := exec.Command("R", "--vanilla")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Println("get stdinPipe failed, err:", err)
		return nil
	}

	id, _ := gonanoid.New()
	logger := newProcessorLogger(name, id)
	// redirect log to file
	cmd.Stdout = logger
	cmd.Stderr = logger

	err = cmd.Start()
	if err != nil {
		log.Println("start cmd failed, err:", err)
		return nil
	}

	proc := &Processor{
		id:     id,
		name:   name,
		cmd:    cmd,
		inPipe: &stdinPipe,
		logger: logger,
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

// ForceClose This method will directly kill the process
func (p *Processor) ForceClose() error {
	return p.cmd.Process.Kill()
}

// QuitR This method will wait for the current task to complete before exiting
func (p *Processor) QuitR() error {
	_, err := p.Exec("q()")
	return err
}

// Exec
// This func will automatically add \n at the end
// Please ensure exec one line
func (p *Processor) Exec(format string, a ...any) (i int, err error) {
	defer func() {
		// log failed code
		if err != nil {
			p.logger.Error("Exec code failed: ", fmt.Sprintf(format, a...))
		}
	}()
	if p.inPipe == nil {
		return 0, fmt.Errorf("inPipe is nil, exec failed")
	}
	return fmt.Fprintf(*p.inPipe, format+"\n", a...)
}
