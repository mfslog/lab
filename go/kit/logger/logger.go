package logger

import (
	"github.com/go-kit/kit/log"
	kitLogrus "github.com/go-kit/kit/log/logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger log.Logger
)

func Init(level string, path, fileName string) {
	hook := &lumberjack.Logger{
		Filename:   path + fileName, // 日志文件路径
		MaxSize:    128,             // megabytes
		MaxBackups: 30,              // 最多保留300个备份
		MaxAge:     7,               // days
		Compress:   true,            // 是否压缩 disabled by default
	}
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(hook)
	logrusLogger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "02-01-2006 15:04:05", FullTimestamp: true})
	Logger = kitLogrus.NewLogrusLogger(logrusLogger)
}
