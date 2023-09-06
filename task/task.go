package task

import (
	"rscheduler/config"
	"time"

	"rscheduler/monitor"
	"rscheduler/rslog"

	"go.uber.org/zap"
)

type TaskRunner interface {
	CommendList() []string
}

type Task struct {
	Name      string `json:"Name"`
	ID        string `json:"ID"`
	CreatedAt time.Time
	StartAt   time.Time
	Runner    TaskRunner
	Logger    *rslog.RsLogger
	PID       int
	StopLog   chan struct{}
}

type HiPlotTask struct {
	InputFile        string `json:"inputFile"`
	ConfFile         string `json:"confFile"`
	OutputFilePrefix string `json:"outputFilePrefix"`
	Tool             string `json:"tool"`
	Module           string `json:"module"`
	ID               string `json:"ID"`
	Name             string `json:"Name"` // HiPlot
}

func (t *Task) CommendList() []string {
	commends := make([]string, 2)
	commends[0] = `taskID = "` + t.ID + `"`
	commends[1] = `source("./rscript/` + t.Name + `.R")`
	return commends
}

// StartLogger Auto record process runtime status
func (t *Task) StartLogger() {
	go func() {
		t.Logger.Infow("Start task logger", zap.Int("PID", t.PID))
		ticker := time.NewTicker(100 * time.Millisecond)
	LogLoop:
		for {
			select {
			case <-ticker.C:
				info := monitor.GetProcessInfo(t.PID)
				if info == nil {
					continue
				}
				t.Logger.Infow("Task runtime status",
					zap.Float64("cpuUsage", info.CpuUsage),
					zap.Float32("memUsage", info.MemUsage),
					zap.Uint64("RSS", info.RSS))
				info.Recycle()
			case <-t.StopLog:
				break LogLoop
			}
		}
	}()
}

func (t *Task) StopLogger() {
	t.StopLog <- struct{}{}
}

func (t *Task) SetStartTime() {
	t.StartAt = time.Now()
}

func (t *Task) IsTimeout() bool {
	return time.Now().Sub(t.StartAt) > time.Duration(config.Config.TaskLimit.TaskTimeout)*time.Second
}
