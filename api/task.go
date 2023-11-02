package api

import (
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
