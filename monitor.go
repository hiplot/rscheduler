package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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
			if !enableNewTask() {
				// TODO GC procMap
			}
		}
	}()
}

func getMemInfo() {
	vm, err := mem.VirtualMemory()
	if err != nil {
		Logger.Error("get memory use failed, err: ", err)
		return
	}
	MemPrecent = int(vm.UsedPercent)
}

func getCPUInfo() {
	// TODO check it is adapt linux
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		Logger.Error("get cpu use failed, err: ", err)
		return
	}
	CPUPrecent = int(cpuPercent[0])
}

func enableNewTask() bool {
	if MemPrecent > AllowMaxMemoryPercent || CPUPrecent > AllowMaxCPUPercent {
		return false
	}
	return true
}
