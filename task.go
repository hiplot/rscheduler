package main

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Task struct {
	name string
	id   string
	//logger *rsLogger
}

var DelayTask *ProcList // TODO store task

func NewTask(name string) *Task {
	id, _ := gonanoid.New()
	return &Task{
		name: name,
		id:   id,
		//logger: newTaskLogger(name, id),
	}
}
