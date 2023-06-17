package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

const (
	SEP          = string(filepath.Separator)
	ProcessorLog = "processor"
	GlobalLog    = "global"
)

type rsLogger struct {
	*zap.SugaredLogger
	*os.File
}

func newProcLogger(name, id string) *rsLogger {
	return newLogger(ProcessorLog, name, id)
}

func newGlobalLogger() *rsLogger {
	return newLogger(GlobalLog, "", "global")
}

func newLogger(kind, name, id string) *rsLogger {
	encoder := getEncoder()
	writeSyncer, file, err := getWriteSyncer(kind, name, id)
	var core zapcore.Core
	if err != nil {
		// TODO 发送日志初始化失败通知
		// 直接将日志打到控制台
		writeSyncer = zapcore.AddSync(os.Stdout)
		file = os.Stdout
	}
	core = zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	log := zap.New(core)
	return &rsLogger{log.Sugar(), file}
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getWriteSyncer(kind, name, id string) (zapcore.WriteSyncer, *os.File, error) {
	now := time.Now().Format("2006-01-02")
	Path, _ := os.Getwd()
	LogPath := Path + SEP + "log" + SEP + now + SEP + kind + SEP + name
	err := os.MkdirAll(LogPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	filePath := LogPath + SEP + id + ".txt"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return nil, nil, err
	}
	return zapcore.AddSync(file), file, nil
}
