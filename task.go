package main

type Task struct {
	name string
	id   string
}

var DelayTask []*Task // TODO store task

func NewTask(name string) *Task {
	//id, _ := gonanoid.New()
	// TODO create id
	return &Task{
		name: name,
		id:   "123",
	}
}
