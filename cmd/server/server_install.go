/**
 * @Time : 7/21/21 9:44 AM
 * @Author : solacowa@gmail.com
 * @File : server_install
 * @Software: GoLand
 */

package server

import (
	"context"
	"fmt"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/kplcloud/kplcloud/src/repository"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/config"
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	"github.com/kplcloud/kplcloud/src/encode"
	installSrc "github.com/kplcloud/kplcloud/src/install"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/middleware"
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
		if db != nil {
			_ = level.Debug(logger).Log("db", "close", "err", db.Close())
		}
		if rds != nil {
			_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
		}
	}()

	installSvc = installSrc.New(logger, cf, configPath, store)

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
		//middleware.CheckAuthMiddleware(logger, cacheSvc, tracer),
	}
	tokenEms = append(tokenEms, ems...)
	installSvc.StoreToConfig(context.Background())

	r := mux.NewRouter()

	// 以下安装模块
	r.PathPrefix("/install").Handler(http.StripPrefix("/install", installSrc.MakeHTTPHandler(installSvc, tokenEms, opts)))

	// 心跳检测
	r.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
	// web页面
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(webPath + "/install/")))

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

	logger = logging.SetLogging(logger, cf)
	if strings.EqualFold(cf.GetString("database", "drive"), "mysql") {
		dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
			cf.GetString("database", "user"), cf.GetString("database", "password"),
			cf.GetString("database", "host"), cf.GetInt("database", "port"),
			cf.GetString("database", "database"))

		// 连接数据库
		db, err = mysqlclient.NewMysql(dbUrl, true)
		if err != nil {
			_ = level.Error(logger).Log("db", "connect", "err", err)
			err = encode.ErrInstallDbConnect.Wrap(err)
			return err
		}
		store = repository.New(db, logger, "traceId", nil, nil)
	}

	return err
}
