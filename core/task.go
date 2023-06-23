package core

import (
	"encoding/json"
	"rscheduler/global"
)

type Task struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
	//logger *rsLogger
}

func NewTask(b []byte) (t *Task) {
	return decode(b)
}

func decode(s []byte) *Task {
	t := new(Task)
	err := json.Unmarshal(s, t)
	if err != nil {
		global.Logger.Error("Decode task failed, err: " + err.Error())
		return nil
	}
	return t
}
