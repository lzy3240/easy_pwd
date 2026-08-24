package log

import (
	"easy_pwd/global"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path"
	"sync"
	"time"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func Instance() *slog.Logger {
	if logger == nil {
		once.Do(func() {
			InitLog("")
		})
	}
	return logger
}

func InitLog(fPath string) {
	var logPath = ""
	if fPath != "" {
		logPath = fPath
	} else {
		logPath = global.LogFileDir
	}
	hook := &lumberjack.Logger{
		Filename:   path.Join(logPath, time.Now().Format("2006-01-02")+".log"), // 日志文件路径
		MaxSize:    100,                                                        // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 50,                                                         // 日志文件最多保存多少个备份
		MaxAge:     30,                                                         // 文件最多保存多少天
		Compress:   true,                                                       // 是否压缩
		LocalTime:  true,                                                       // 使用本地时间
	}

	writer := io.MultiWriter(hook, os.Stdout) //构造混合输出io, 控制台+文件

	l := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		AddSource:   true,            // 记录日志位置
		Level:       slog.LevelDebug, // 设置日志级别
		ReplaceAttr: nil,
	}))
	slog.SetDefault(l)

	logger = l
}

func SetValue(key string, value any) slog.Attr {
	return slog.Any(key, value)
}
