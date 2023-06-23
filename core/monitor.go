package core

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"rscheduler/config"
	"rscheduler/global"
	"time"
)

var cpuPercent int
var memPercent int

func InitMonitor() {
	go func() {
		for {
			getMemInfo() // no block
			getCPUInfo() // block 1s
			if !EnableNewTask() {
				// TODO Try to gc RScheduler
			}
		}
	}()
}

func getMemInfo() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		global.Logger.Error("get memory use failed, err: ", err)
		return
	}
	memPercent = int(vm.UsedPercent)
}

func getCPUInfo() {
	// TODO check it is adapt linux
	// This Method will spend 1s
	nowCPUPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		global.Logger.Error("get cpu use failed, err: ", err)
		return
	}
	cpuPercent = int(nowCPUPercent[0])
}

func EnableNewTask() bool {
	cfg := config.Config.TaskLimit
	if memPercent > cfg.MaxMem || cpuPercent > cfg.MaxCPU {
		return false
	}
	return true
}
