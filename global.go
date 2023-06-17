package main

var Logger *rsLogger

var ProcMap procMap

var CPUPrecent int
var MemPrecent int

func init() {
	Logger = newGlobalLogger()
	ProcMap = procMap{m: make(map[string]*ProcList)}
}
