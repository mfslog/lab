package log

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(level string, path, fileName string) {
	hook := &lumberjack.Logger{
		Filename:   path + fileName, // 日志文件路径
		MaxSize:    128,             // megabytes
		MaxBackups: 30,              // 最多保留300个备份
		MaxAge:     7,               // days
		Compress:   true,            // 是否压缩 disabled by default
	}

}
