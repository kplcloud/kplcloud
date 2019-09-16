/**
 * @Time : 2019-07-12 10:31
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package logging

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitlogrus "github.com/go-kit/kit/log/logrus"
	"github.com/kplcloud/kplcloud/src/config"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

func SetLogging(logger log.Logger, cf *config.Config) log.Logger {
	if cf.GetString("server", "logs_path") != "" {
		logrusLogger, err := LogrusLogger(cf.GetString("server", "logs_path"))
		if err != nil {
			panic(err)
		}
		logLevel, _ := logrus.ParseLevel(cf.GetString("server", "log_level"))
		logrusLogger.SetLevel(logLevel)
		logger = kitlogrus.NewLogrusLogger(logrusLogger)
	} else {
		logger = log.NewLogfmtLogger(log.StdlibWriter{})
		logger = level.NewFilter(logger, logLevel(cf.GetString("server", "log_level")))
	}
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.WithPrefix(logger, "app", cf.GetString("server", "app_name"))

	return logger
}

func logLevel(logLevel string) (opt level.Option) {
	switch logLevel {
	case "warn":
		opt = level.AllowWarn()
	case "error":
		opt = level.AllowError()
	case "debug":
		opt = level.AllowDebug()
	case "info":
		opt = level.AllowInfo()
	case "all":
		opt = level.AllowAll()
	default:
		opt = level.AllowNone()
	}

	return
}

func LogrusLogger(filePath string) (*logrus.Logger, error) {
	//path, fileName := filepath.Split(filePath)
	linkFile, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	logrusLogger := logrus.New()
	writer, err := rotatelogs.New(
		linkFile+"-%Y-%m-%d",
		rotatelogs.WithLinkName(linkFile),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*365),   // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)
	if err != nil {
		logrusLogger.Error("Init log failed, err:", err)
		return nil, err
	}

	logrusLogger.SetOutput(writer)
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	return logrusLogger, nil
}
