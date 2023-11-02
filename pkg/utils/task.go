package utils

import "rscheduler/task"

func GetTaskID(t *task.Task) string {
	if t == nil {
		return ""
	}
	return t.ID
}
