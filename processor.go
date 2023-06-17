package main

import (
	"container/list"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"io"
	"os/exec"
)

type Proc struct {
	id     string
	name   string
	cmd    *exec.Cmd
	inPipe *io.WriteCloser
	task   *Task
	logger *rsLogger
}

type ProcList struct {
	list.List
}

func newProcList() *ProcList {
	return &ProcList{list.List{}}
}

// newProc
// Build a new r session
// Auto load nameFunc.R
// In this file, we can write some init code
func newProc(name string) *Proc {
	Logger.Info("Create new processor")
	cmd := exec.Command("R", "--vanilla")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		Logger.Error("get stdinPipe failed, err:", err)
		return nil
	}

	id, _ := gonanoid.New()
	logger := newProcLogger(name, id)
	// redirect log to file
	cmd.Stdout = logger
	cmd.Stderr = logger

	err = cmd.Start()
	if err != nil {
		Logger.Error("start cmd failed, err:", err)
		return nil
	}

	proc := &Proc{
		id:     id,
		name:   name,
		cmd:    cmd,
		inPipe: &stdinPipe,
		logger: logger,
	}

	_, err = proc.Exec(`source("./rscript/%sInit.R")`, name)
	if err != nil {
		Logger.Error("Exec failed, err: ", err)
		_ = proc.ForceClose()
		return nil
	}

	ProcMap.lock.Lock()
	defer ProcMap.lock.Unlock()

	if ProcMap.m[name] == nil {
		ProcMap.m[name] = newProcList()
	}

	ProcMap.m[name].PushBack(proc)
	return proc
}

// ForceClose This method will directly kill the process
func (p *Proc) ForceClose() error {
	return p.cmd.Process.Kill()
}

// Close This method will wait for the current task to complete before exiting
func (p *Proc) Close() error {
	_, err := p.Exec("q()")
	return err
}

// Exec
// This func will automatically add \n at the end
// Please ensure exec one line
func (p *Proc) Exec(format string, a ...any) (i int, err error) {
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
