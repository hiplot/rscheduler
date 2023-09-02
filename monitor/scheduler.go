package monitor

import (
	"go.uber.org/zap"
	"log"
	"time"

	"rscheduler/config"
	"rscheduler/global"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type State int

const (
	None State = 1 << iota
	MemLimit
	CPULimit
)

var globalCpuPercent int
var globalMemPercent int

func InitMonitor() {
	go func() {
		for {
			getMemInfo() // no block
			getCPUInfo() // block 1s
			state := checkState()
			if state.ArriveMaxMem() {
				global.GCChan <- struct{}{}
			}
		}
	}()
	log.Println("监控初始化成功")
}

func getMemInfo() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		global.Logger.Error("get memory use failed, err: ", err)
		return
	}
	globalMemPercent = int(vm.UsedPercent)
}

func getCPUInfo() {
	// This method will spend 1s to collect cpu info
	nowCPUPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		global.Logger.Error("get cpu use failed, err: ", err)
		return
	}
	globalCpuPercent = int(nowCPUPercent[0])
}

func EnableNewTask() bool {
	return checkState() == None
}

func checkState() State {
	state := None
	cfg := config.Config.TaskLimit
	if globalMemPercent > cfg.MaxMem {
		global.Logger.Warnw("Memory use is too high", zap.Int("current", globalMemPercent), zap.Int("allowMax", cfg.MaxMem))
		state.WithMaxMem()
	}
	if globalCpuPercent > cfg.MaxCPU {
		global.Logger.Warnw("CPU use is too high", zap.Int("current", globalCpuPercent), zap.Int("allowMax", cfg.MaxCPU))
		state.WithMaxCPU()
	}
	return state
}

func (s *State) WithMaxMem() {
	*s |= MemLimit
}

func (s *State) WithMaxCPU() {
	*s |= CPULimit
}

func (s *State) ArriveMaxMem() bool {
	return *s&MemLimit != 0
}

func (s *State) ArriveMaxCPU() bool {
	return *s&CPULimit != 0
}
