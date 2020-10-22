/**
 * @Time : 2019-07-12 10:31
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package logging

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/lestrrat-go/file-rotatelogs"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	LoggerRequestId = "trace-id"
	TraceId         = "trace-id"
)

func SetLogging(logger kitlog.Logger, logPath, levelOut string) kitlog.Logger {
	if logPath != "" {
		// default log
		logger = defaultLogger(logPath)
		logger = kitlog.WithPrefix(logger, "ts", kitlog.TimestampFormat(func() time.Time {
			return time.Now()
		}, "2006-01-02 15:04:05"))
	} else {
		//logger = kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})
		logger = term.NewLogger(os.Stdout, kitlog.NewLogfmtLogger, colorFunc())
		logger = kitlog.WithPrefix(logger, "ts", kitlog.TimestampFormat(func() time.Time {
			return time.Now()
		}, "2006-01-02 15:04:05"))
	}
	logger = level.NewFilter(logger, logLevel(levelOut))
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

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

func defaultLogger(filePath string) kitlog.Logger {
	linkFile, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	writer, err := rotatelogs.New(
		linkFile+"-%Y-%m-%d",
		rotatelogs.WithLinkName(linkFile),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*365),   // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return kitlog.NewLogfmtLogger(writer)
}

func colorFunc() func(keyvals ...interface{}) term.FgBgColor {
	return func(keyvals ...interface{}) term.FgBgColor {
		for i := 0; i < len(keyvals)-1; i += 2 {
			if keyvals[i] != "level" {
				continue
			}
			val := fmt.Sprintf("%v", keyvals[i+1])
			switch val {
			case "debug":
				return term.FgBgColor{Fg: term.DarkGray}
			case "info":
				return term.FgBgColor{Fg: term.Blue}
			case "warn":
				return term.FgBgColor{Fg: term.Yellow}
			case "error":
				return term.FgBgColor{Fg: term.Red}
			case "crit":
				return term.FgBgColor{Fg: term.Gray, Bg: term.DarkRed}
			default:
				return term.FgBgColor{}
			}
		}
		return term.FgBgColor{}
	}
}
