package core

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"rscheduler/global"
	"time"
)

const (
	AllowMaxCPUPercent    = 90
	AllowMaxMemoryPercent = 80
)

func init() {
	go func() {
		for {
			getMemInfo()
			getCPUInfo()
			if !EnableNewTask() {
				// TODO GC procMap
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
	MemPercent = int(vm.UsedPercent)
}

func getCPUInfo() {
	// TODO check it is adapt linux
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		global.Logger.Error("get cpu use failed, err: ", err)
		return
	}
	CPUPercent = int(cpuPercent[0])
}

func EnableNewTask() bool {
	if MemPercent > AllowMaxMemoryPercent || CPUPercent > AllowMaxCPUPercent {
		return false
	}
	return true
}
