package processor

import (
	"container/list"
	"fmt"
	"io"
	"os/exec"
	"rscheduler/config"
	"time"

	"rscheduler/global"
	"rscheduler/monitor"
	"rscheduler/rslog"
	"rscheduler/task"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

type Proc struct {
	ID         string
	Name       string
	PID        int    // 进程id
	TakNum     uint32 // 总执行任务数
	InitialMem uint64 // 初始内存占用
	NowMem     uint64 // 当前内存占用
	Running    bool   // 是否正在运行
	CMD        *exec.Cmd
	InPipe     *io.WriteCloser
	Task       *task.Task
	Logger     *rslog.RsLogger
}

type ProcList struct {
	list.List
}

func NewProcList() *ProcList {
	return &ProcList{list.List{}}
}

// NewProc
// Build a new r session
// Auto load nameInit.R
// In this file, we can write some init code
func NewProc(name string) *Proc {
	global.Logger.Info("Create new processor")
	cmd := exec.Command("R", "--vanilla")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		global.Logger.Error("get stdinPipe failed, err:", err)
		return nil
	}

	id, _ := gonanoid.New()
	logger := rslog.NewProcLogger(name, id)
	// redirect log to file
	cmd.Stdout = logger
	cmd.Stderr = logger

	err = cmd.Start()
	if err != nil {
		global.Logger.Error("start CMD failed, err:", err)
		return nil
	}

	proc := &Proc{
		ID:     id,
		Name:   name,
		CMD:    cmd,
		PID:    cmd.Process.Pid,
		InPipe: &stdinPipe,
		Logger: logger,
	}

	// load init file
	_, err = proc.Exec(`source("./rscript/%sInit.R")`, name)
	if err != nil {
		global.Logger.Error("Exec failed, err: ", err)
		_ = proc.ForceClose()
		return nil
	}
	return proc
}

// ForceClose This method will directly kill the process
func (p *Proc) ForceClose() error {
	return p.CMD.Process.Kill()
}

// Close This method will wait for the current Task to complete before exiting
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
			p.Logger.Error("Exec code failed: ", fmt.Sprintf(format, a...))
		}
	}()
	if p.InPipe == nil {
		return 0, fmt.Errorf("InPipe is nil, exec failed")
	}
	return fmt.Fprintf(*p.InPipe, format+"\n", a...)
}

func (p *Proc) HealthCheck() bool {
	return p.CMD != nil && p.CMD.ProcessState == nil
}

func (p *Proc) IsIdle() bool {
	return !p.Running
}

func (p *Proc) RefreshProcState() {
	pi := monitor.GetProcessInfo(p.PID)
	if p.InitialMem == 0 {
		p.InitialMem = pi.RSS
	}
	p.NowMem = pi.RSS
	pi.Recycle()
}

func (p *Proc) SetRun() {
	p.Running = true
}

func (p *Proc) CancelRun() {
	p.Running = false
}

func (p *Proc) Complete() {
	global.Logger.Infow("Task complete success",
		zap.String("taskName", p.Task.Name),
		zap.String("taskID", p.Task.ID),
		zap.Duration("useTime", time.Now().Sub(p.Task.StartAt)),
	)
	p.Task.StopLogger()
	p.UnbindTask()
	p.CancelRun()
}

func (p *Proc) BindTask(t *task.Task) {
	p.Task = t
	t.PID = p.PID
	p.TakNum++
}

func (p *Proc) UnbindTask() {
	p.Task = nil
}

func (p *Proc) Start() error {
	if p.Task == nil {
		return fmt.Errorf("task is nil, please bind Task first")
	}
	p.Task.StartLogger()
	p.Task.SetStartTime()
	// do Task
	commends := p.Task.Runner.CommendList()
	for _, commend := range commends {
		_, err := p.Exec(commend)
		if err != nil {
			return fmt.Errorf("exec failed, err: %v", err)
		}
	}
	return nil
}

func (p *Proc) MemCheck() bool {
	p.RefreshProcState()
	if p.IsIdle() {
		return p.idleMemCheck()
	}
	return p.busyMemCheck()
}

func (p *Proc) idleMemCheck() bool {
	return p.NowMem < uint64(config.Config.TaskLimit.MaxIdleProcessorMem)
}

func (p *Proc) busyMemCheck() bool {
	return p.NowMem < uint64(config.Config.TaskLimit.MaxBusyProcessorMem)
}
