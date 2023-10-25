package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rscheduler/processor"
	"rscheduler/scheduler"
	"rscheduler/task"
)

type processorAPI struct {
}

var ProcessorAPI = &processorAPI{}

func (p processorAPI) Info(c *gin.Context) {
	scheduler.RScheduler.Lock.RLock()
	defer scheduler.RScheduler.Lock.RUnlock()
	processorInfoList := make([]processorInfo, 0)
	getTaskID := func(t *task.Task) string {
		if t == nil {
			return ""
		}
		return t.ID
	}
	for _, procList := range scheduler.RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			procInfo := processorInfo{
				ID:           proc.ID,
				Name:         proc.Name,
				PID:          proc.PID,
				TotalTaskNum: proc.TakNum,
				Running:      proc.Running,
				TaskID:       getTaskID(proc.Task),
			}
			processorInfoList = append(processorInfoList, procInfo)
		}
	}
	c.JSON(http.StatusOK, ProcessorInfoResponse{
		BaseResponse:  NewBaseSuccessResponse(),
		ProcessorInfo: processorInfoList,
	})
}
