package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rscheduler/pkg/utils"
	"rscheduler/processor"
	"rscheduler/scheduler"
)

type processorAPI struct {
}

var ProcessorAPI = &processorAPI{}

func (p processorAPI) Info(c *gin.Context) {
	// TODO 后续重构，放在这边不太好
	scheduler.RScheduler.Lock.RLock()
	defer scheduler.RScheduler.Lock.RUnlock()

	processorInfoList := make([]processorInfo, 0)

	for _, procList := range scheduler.RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			procInfo := processorInfo{
				ID:           proc.ID,
				Name:         proc.Name,
				PID:          proc.PID,
				TotalTaskNum: proc.TakNum,
				Running:      proc.Running,
				TaskID:       utils.GetTaskID(proc.Task),
			}
			processorInfoList = append(processorInfoList, procInfo)
		}
	}

	c.JSON(http.StatusOK, ProcessorInfoResp{
		BaseResp:      NewBaseSuccessResp(),
		ProcessorInfo: processorInfoList,
	})
}

func (p processorAPI) Delete(c *gin.Context) {
	var req ProcessorDeleteReq
	if err := c.ShouldBindQuery(&req); err != nil {
		failResp := ProcessorDeleteResp{
			BaseResp: NewBaseFailResp(),
			Success:  false,
			Info:     err.Error(),
		}
		Response(c, failResp)
		return
	}

	err := scheduler.KillProcByProcID(req.ID, req.Force)
	if err != nil {
		failResp := ProcessorDeleteResp{
			BaseResp: NewBaseFailResp(),
			Success:  false,
			Info:     err.Error(),
		}
		Response(c, failResp)
		return
	}

	successResp := ProcessorDeleteResp{
		BaseResp: NewBaseSuccessResp(),
		Success:  true,
		Info:     fmt.Sprintf("Processor %s delete success", req.ID),
	}
	Response(c, successResp)
}
