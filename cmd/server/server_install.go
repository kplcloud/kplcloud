/**
 * @Time : 7/21/21 9:44 AM
 * @Author : solacowa@gmail.com
 * @File : server_install
 * @Software: GoLand
 */

package server

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/syspermission"
	"github.com/kplcloud/kplcloud/src/pkg/sysrole"
	"github.com/kplcloud/kplcloud/src/pkg/sysuser"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"time"

	installSrc "github.com/kplcloud/kplcloud/src/install"
)

var (
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "安装",
		Example: `## 安装kplcloud
kplcloud install -p :8080
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return install()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := installPre(); err != nil {
				_ = level.Error(logger).Log("cmd", "start.PreRunE", "err", err.Error())
				return err
			}
			return nil
		},
	}

	installSvc installSrc.Service
)

func install() (err error) {
	// 关闭资源连接
	defer func() {
		//_ = level.Debug(logger).Log("db", "close", "err", db.Close())
		if rds != nil {
			_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
		}
	}()
	logger = logging.SetLogging(logger, cf)

	installSvc = installSrc.New(logger, cf, configPath)

	//ctx := context.Background()

	//return installSvc.Database(ctx, "mysql","127.0.0.1", 32780, "root", "admin", "kplcloud")

	g := &group.Group{}

	httpLogger := log.With(logger, "component", "http")

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encode.JsonError),
		kithttp.ServerErrorHandler(logging.NewLogErrorHandler(level.Error(logger), apiSvc, cf.GetInt("service", "alarm.appId"))),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			guid := request.Header.Get("X-Request-Id")
			token := request.Header.Get("Authorization")

			ctx = context.WithValue(ctx, logging.TraceId, guid)
			ctx = context.WithValue(ctx, "token-context", token)
			return ctx
		}),
		kithttp.ServerBefore(middleware.TracingServerBefore(tracer)),
	}

	ems := []endpoint.Middleware{
		middleware.TracingMiddleware(tracer),                                                      // 1
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 0
	}

	tokenEms := []endpoint.Middleware{
		middleware.CheckAuthMiddleware(logger, cacheSvc, tracer),
	}
	tokenEms = append(tokenEms, ems...)

	r := mux.NewRouter()

	// 以下安装模块
	r.PathPrefix("/install/").Handler(http.StripPrefix("/install/", sysuser.MakeHTTPHandler(sysUserSvc, tokenEms, opts)))

	// 以下为业务模块

	// 对外metrics
	r.Handle("/metrics", promhttp.Handler())
	// 心跳检测
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
	// web页面
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(webPath)))

	http.Handle("/", accessControl(r, httpLogger))

	g.Add(func() error {
		_ = level.Debug(httpLogger).Log("transport", "HTTP", "addr", httpAddr)
		return http.ListenAndServe(httpAddr, nil)
	}, func(e error) {
		_ = level.Error(httpLogger).Log("transport", "HTTP", "httpListener.Close", "http", "err", e.Error())
		os.Exit(1)
	})

	initCancelInterrupt(g)

	_ = level.Error(logger).Log("server exit", g.Run())
	return nil
}

func installPre() error {
	cf, err = config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	if appName == "" {
		appName = cf.GetString(config.SectionServer, "app.name")
	}

	return err
}
