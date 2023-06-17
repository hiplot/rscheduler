package main

var Logger *rsLogger

var ProcMap procMap

func init() {
	Logger = newGlobalLogger()
	ProcMap = procMap{m: make(map[string]*ProcList)}
}
