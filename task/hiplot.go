package task

import (
	"encoding/json"
	"fmt"
	"time"

	"rscheduler/global"
	"rscheduler/rslog"
)

func NewHiPlotTask(b []byte) (t *Task) {
	ht := new(HiPlotTask)
	err := json.Unmarshal(b, ht)
	if err != nil {
		global.Logger.Error("Decode HiPlotTask failed, err: " + err.Error())
		return nil
	}
	t = &Task{
		Name:      ht.Name,
		ID:        ht.ID,
		CreatedAt: time.Now(),
		Runner:    ht,
		Logger:    rslog.NewTaskLogger(ht.Name, ht.ID),
		StopLog:   make(chan struct{}),
	}
	return t
}

func (t *HiPlotTask) CommendList() []string {
	commends := make([]string, 3)
	commends[0] = `ifelse(grepl("rscript" ,getwd()), "", setwd("./rscript/"))  `
	commends[1] = `taskID = "` + t.ID + `"`
	commends[2] = fmt.Sprintf(`hiFunc("%s", "%s", "%s", "%s", "%s")`, t.InputFile, t.ConfFile, t.OutputFilePrefix, t.Tool, t.Module)
	return commends
}
