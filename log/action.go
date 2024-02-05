package log

import (
	"fmt"
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
)

func Debug(msg string, fields ...zap.Field) {
	checkInit()
	log.loggerAction.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	checkInit()
	log.loggerAction.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	checkInit()
	log.loggerAction.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	checkInit()
	log.loggerAction.Error(msg, fields...)
}

func checkInit() {
	if !log.isInit {
		panic(fmt.Errorf("log package not init"))
	}
}

func getWriter(filename string, times int) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(times)),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
