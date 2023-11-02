package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liang09255/lutils/conv"
	"net/http"
	"rscheduler/global"
	"rscheduler/processor"
	"rscheduler/scheduler"
	"time"
)

type taskAPI struct{}

var TaskAPI = &taskAPI{}

func (t taskAPI) TaskCompleteHandler(c *gin.Context) {
	taskName := c.Query("taskName")
	taskID := c.Query("taskID")
	kill := conv.ToBool(c.Query("kill"))
	scheduler.RScheduler.TaskComplete(taskName, taskID, kill)
	c.JSON(http.StatusOK, NewBaseSuccessResp())
}

func (t taskAPI) Info(c *gin.Context) {
	// TODO 后续重构，放在这边不太好
	scheduler.RScheduler.Lock.RLock()
	defer scheduler.RScheduler.Lock.RUnlock()

	taskInfoList := make([]taskInfo, 0)

	for _, procList := range scheduler.RScheduler.M {
		for i := procList.Front(); i != nil; i = i.Next() {
			proc := i.Value.(*processor.Proc)
			if proc.IsIdle() {
				continue
			}
			t := proc.Task
			if t == nil {
				global.Logger.Errorln("Proc have not task but status is running")
				continue
			}
			tInfo := taskInfo{
				ID:          t.ID,
				Name:        t.Name,
				ProcessorID: proc.ID,
				Memory:      proc.NowMem,
				CPU:         conv.ToUint64(proc.NowCPU),
				CreateAt:    t.CreatedAt.Unix(),
				StartAt:     t.StartAt.Unix(),
				RunTime:     int64(time.Now().Sub(t.StartAt).Seconds()),
			}
			taskInfoList = append(taskInfoList, tInfo)
		}
	}

	c.JSON(http.StatusOK, TaskInfoResp{
		BaseResp: NewBaseSuccessResp(),
		TaskInfo: taskInfoList,
	})
}

func (t taskAPI) Delete(c *gin.Context) {
	var req TaskDeleteReq
	if err := c.ShouldBindQuery(&req); err != nil {
		failResp := TaskDeleteResp{
			BaseResp: NewBaseFailResp(),
			Success:  false,
			Info:     err.Error(),
		}
		Response(c, failResp)
		return
	}

	err := scheduler.KillProcByTaskID(req.ID)
	if err != nil {
		failResp := TaskDeleteResp{
			BaseResp: NewBaseFailResp(),
			Success:  false,
			Info:     err.Error(),
		}
		Response(c, failResp)
		return
	}

	successResp := TaskDeleteResp{
		BaseResp: NewBaseSuccessResp(),
		Success:  true,
		Info:     fmt.Sprintf("Task %s delete success", req.ID),
	}
	Response(c, successResp)
}
