package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
初始化日志配置
logPath：日志存放路径
logName：日志命名
logDay：日志持久化天数
logLevel：日志等级
*/
func InitLog(logPath, logName string, logDay int, logLevel Level) {
	// 设置日志路径、名称以及存储天数
	infoWriter := getWriter(fmt.Sprintf("%s/%s.log", logPath, logName), logDay)       // debug、info存储路径
	warnWriter := getWriter(fmt.Sprintf("%s/%s_error.log", logPath, logName), logDay) // warning、error存储路径
	// 根据配置设置日志等级
	switch logLevel {
	case ErrorLevel:
		cores = zapcore.NewTee(
			cores,
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), errorLevel),
		)
	case WarningLevel:
		cores = zapcore.NewTee(
			cores,
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		)
	case InfoLevel:
		cores = zapcore.NewTee(
			cores,
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		)
	case DebugLevel:
		cores = zapcore.NewTee(
			cores,
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), debugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), debugLevel),
		)
	}
	log.loggerAction = zap.New(cores, zap.AddCaller())
	log.isInit = true
}
