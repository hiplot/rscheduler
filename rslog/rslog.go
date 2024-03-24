package rslog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	SEP          = string(filepath.Separator)
	ProcessorLog = "processor"
	GlobalLog    = "global"
	TaskLog      = "task"
)

var once sync.Once

type RsLogger struct {
	*zap.SugaredLogger
	*os.File
}

func NewProcLogger(name, id string) *RsLogger {
	return newLogger(ProcessorLog, name, id)
}

func NewGlobalLogger() *RsLogger {
	return newLogger(GlobalLog, "", "global")
}

func NewTaskLogger(name, id string) *RsLogger {
	return newLogger(TaskLog, name, id)
}

func newLogger(kind, name, id string) *RsLogger {
	encoder := getEncoder()
	writeSyncer, file, err := getWriteSyncer(kind, name, id)
	var core zapcore.Core
	if err != nil {
		log.Println("get writeSyncer failed, err: ", err)
		// 直接将日志打到控制台
		writeSyncer = zapcore.AddSync(os.Stdout)
		file = os.Stdout
	}
	core = zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	zapLogger := zap.New(core)
	return &RsLogger{zapLogger.Sugar(), file}
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
}

func getWriteSyncer(kind, name, id string) (zapcore.WriteSyncer, *os.File, error) {
	now := time.Now().Format("2006-01-02")
	Path, _ := os.Getwd()
	once.Do(func() {
		_ = os.Mkdir("./log", 0777)
	})
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
