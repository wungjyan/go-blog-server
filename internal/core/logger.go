package core

import (
	"blog-server/global"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logConf := global.AppConfig.Log

	// 设置日志级别
	level, err := logrus.ParseLevel(logConf.Level)
	if err != nil {
		logger.Warnf("无效的日志级别 '%s'，使用默认级别 'info'", logConf.Level)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置日志格式
	switch strings.ToLower(logConf.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     logConf.Console,
		})
	}

	// 设置是否显示调用者信息
	logger.SetReportCaller(logConf.ShowCaller)

	// 设置输出
	setLogOutput(logger)

	return logger
}

func setLogOutput(logger *logrus.Logger) {
	logConf := global.AppConfig.Log

	if logConf.FilePath == "" {
		logger.SetOutput(os.Stdout)
		return
	}

	if err := ensureLogDir(logConf.FilePath); err != nil {
		logger.Errorf("创建日志目录失败: %v", err)
		logger.SetOutput(os.Stdout)
		return
	}

	fileWriter := &lumberjack.Logger{
		Filename:   logConf.FilePath,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackups,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
	}

	if logConf.Console {
		multiWriter := io.MultiWriter(os.Stdout, fileWriter)
		logger.SetOutput(multiWriter)
	} else {
		logger.SetOutput(fileWriter)
	}
}

func ensureLogDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}
