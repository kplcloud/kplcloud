/**
 * @Time : 2019-07-12 10:31
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"github.com/go-kit/kit/transport"
	"github.com/icowan/config"
	"github.com/jinzhu/gorm"
	"github.com/lestrrat-go/file-rotatelogs"

	"github.com/kplcloud/kplcloud/src/api"
	"github.com/kplcloud/kplcloud/src/encode"
)

const (
	TraceId = "traceId"
)

// LogErrorHandler is a transport error handler implementation which logs an error.
type LogErrorHandler struct {
	logger kitlog.Logger
	apiSvc api.Service
	appId  int
}

func (l *LogErrorHandler) Handle(ctx context.Context, err error) {
	var errDefined bool
	for k := range encode.ResponseMessage {
		if strings.Contains(err.Error(), k.Error().Error()) {
			errDefined = true
			break
		}
	}

	defer func() {
		_ = l.logger.Log("traceId", ctx.Value(TraceId), "err", err.Error())
	}()

	if !errDefined {
		go func(err error) {
			hostname, _ := os.Hostname()
			//res, err := l.apiSvc.Alarm().Warn(ctx, l.appId,
			//	fmt.Sprintf("\nMessage: 未定义错误! \nError: %s \nHostname: %s",
			//		err.Error(),
			//		hostname,
			//	))
			//if err != nil {
			//	log.Println(err)
			//	return
			//}
			//b, _ := json.Marshal(res)
			log.Println(fmt.Sprintf("host: %s, err: %s", hostname, err.Error()))
		}(err)
	}
}

func NewLogErrorHandler(logger kitlog.Logger, apiSvc api.Service, appId int) transport.ErrorHandler {
	return &LogErrorHandler{
		logger: logger,
		apiSvc: apiSvc,
		appId:  appId,
	}
}

func SetLogging(logger kitlog.Logger, cf *config.Config) kitlog.Logger {
	if cf.GetString(config.SectionServer, "log.path") != "" {
		// default log
		logger = defaultLogger(cf.GetString(config.SectionServer, "log.path"))
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
	logger = level.NewFilter(logger, logLevel(cf.GetString(config.SectionServer, "log.level")))
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

type gormLogger struct {
	gorm.LogWriter
	logger kitlog.Logger
}

func (g gormLogger) Println(v ...interface{}) {
	for _, dd := range v {
		_ = level.Debug(g.logger).Log("sql", dd)
	}
}

func NewGormLogger(logger kitlog.Logger) gorm.LogWriter {
	return &gormLogger{logger: logger}
}
