package monitor

import (
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"rscheduler/global"
	"sync"
	"time"
)

type ProcessInfo struct {
	Pid      int     // 进程id
	CpuUsage float64 // cpu使用率
	MemUsage float32 // 内存使用率
	RSS      uint64  // MB
}

var machineMemory *mem.VirtualMemoryStat
var processInfoPool sync.Pool

func init() {
	processInfoPool.New = NewProcessInfo
	var err error
	go func() {
		for {
			machineMemory, err = mem.VirtualMemory()
			if err != nil {
				global.Logger.Error("get machine memory failed, err: ", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func NewProcessInfo() interface{} {
	return &ProcessInfo{}
}

func GetProcessInfo(pid int) *ProcessInfo {
	if pid == 0 {
		global.Logger.Error("get process info failed, pid is 0")
		return &ProcessInfo{}
	}
	info, err := getProcessInfo(pid)
	if err != nil {
		global.Logger.Error("get process info failed, err: ", err)
		return &ProcessInfo{}
	}
	return info
}

func getProcessInfo(pid int) (*ProcessInfo, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}
	cpuPercent, err := p.CPUPercent()
	if err != nil {
		return nil, err
	}

	total := machineMemory.Total
	processMemory, err := p.MemoryInfoEx()
	if err != nil {
		return nil, err
	}
	used := processMemory.RSS
	memPercent := 100 * float32(used) / float32(total)
	rss := used / 1024 / 1024

	pi := processInfoPool.Get().(*ProcessInfo)
	pi.Pid = pid
	pi.CpuUsage = cpuPercent
	pi.MemUsage = memPercent
	pi.RSS = rss
	return pi, nil
}

func (p *ProcessInfo) zero() {
	p.Pid = 0
	p.CpuUsage = 0
	p.MemUsage = 0
	p.RSS = 0
}

func (p *ProcessInfo) Recycle() {
	p.zero()
	processInfoPool.Put(p)
}
