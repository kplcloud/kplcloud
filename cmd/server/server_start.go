/**
 * @Time: 2020/4/23 21:16
 * @Author: solacowa@gmail.com
 * @File: service_start
 * @Software: GoLand
 */

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/account"
	"github.com/kplcloud/kplcloud/src/pkg/auth"
	"github.com/kplcloud/kplcloud/src/pkg/syspermission"
	"github.com/kplcloud/kplcloud/src/pkg/sysrole"
	"github.com/kplcloud/kplcloud/src/pkg/sysuser"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Example: `## 启动命令
kit-admin start -p :8080 -g :8082
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return start()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := prepare(); err != nil {
				_ = level.Error(logger).Log("cmd", "start.PreRunE", "err", err.Error())
				return err
			}
			return nil
		},
	}

	tracer stdopentracing.Tracer

	authSvc          auth.Service
	sysUserSvc       sysuser.Service
	sysRoleSvc       sysrole.Service
	sysPermissionSvc syspermission.Service
	accountSvc       account.Service
)

func start() (err error) {
	// 关闭资源连接
	defer func() {
		_ = level.Debug(logger).Log("db", "close", "err", db.Close())
		if rds != nil {
			_ = level.Debug(logger).Log("redis", "close", "err", rds.Close(context.Background()))
		}
	}()

	// metrics 目前先设置两个指标
	//fieldKeys := []string{"method", "error", "service"}
	// 记数器
	//requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	//	Namespace: namespace,
	//	Subsystem: "activity_service",
	//	Name:      "request_count",
	//	Help:      "请求次数",
	//}, fieldKeys)
	//requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	//	Namespace: namespace,
	//	Subsystem: "activity_service",
	//	Name:      "request_latency_microseconds",
	//	Help:      "请求的总时间(微秒)",
	//}, fieldKeys)

	// 以下是各个服务的初始化
	// 授权登录
	authSvc = auth.New(logger, logging.TraceId, cf, store, cacheSvc, apiSvc, cf.GetBool("server", "debug"))
	authSvc = auth.NewLogging(logger, logging.TraceId)(authSvc)

	// 系统用户
	sysUserSvc = sysuser.New(logger, logging.TraceId, store)
	sysUserSvc = sysuser.NewLogging(logger, logging.TraceId)(sysUserSvc)
	// 系统角色
	sysRoleSvc = sysrole.New(logger, logging.TraceId, store)
	sysRoleSvc = sysrole.NewLogging(logger, logging.TraceId)(sysRoleSvc)
	// 用户信息模块
	accountSvc = account.New(logger, logging.TraceId, store)
	accountSvc = account.NewLogging(logger, logging.TraceId)(accountSvc)
	// 用户信息模块
	sysPermissionSvc = syspermission.New(logger, logging.TraceId, store)
	sysPermissionSvc = syspermission.NewLogging(logger, logging.TraceId)(sysPermissionSvc)

	if tracer != nil {
		authSvc = auth.NewTracing(tracer)(authSvc)
		sysUserSvc = sysuser.NewTracing(tracer)(sysUserSvc)
		sysRoleSvc = sysrole.NewTracing(tracer)(sysRoleSvc)
	}

	g := &group.Group{}

	initHttpHandler(g)
	initGRPCHandler(g)
	initCancelInterrupt(g)

	fmt.Println("asdfasdfasdfasdfasdfasd")

	_ = level.Error(logger).Log("server exit", g.Run())
	return nil
}

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	handlers := make(map[string]string, 3)
	if cf.GetBool("cors", "allow") {
		handlers["Access-Control-Allow-Origin"] = cf.GetString("cors", "origin")
		handlers["Access-Control-Allow-Methods"] = cf.GetString("cors", "methods")
		handlers["Access-Control-Allow-Headers"] = cf.GetString("cors", "headers")
		//reqFun = encode.BeforeRequestFunc(handlers)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range handlers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}
		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}

func initHttpHandler(g *group.Group) {
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

	// 以下为系统模块
	// 授权登录模块
	r.PathPrefix("/admin/auth").Handler(http.StripPrefix("/admin/auth", auth.MakeHTTPHandler(authSvc, ems, opts)))
	r.PathPrefix("/admin/account").Handler(http.StripPrefix("/admin/account", account.MakeHTTPHandler(accountSvc, tokenEms, opts)))
	// 系统用户模块
	r.PathPrefix("/admin/system/user").Handler(http.StripPrefix("/admin/system/user", sysuser.MakeHTTPHandler(sysUserSvc, tokenEms, opts)))
	// 系统角色、权限
	r.PathPrefix("/admin/system/role").Handler(http.StripPrefix("/admin/system/role", sysrole.MakeHTTPHandler(sysRoleSvc, tokenEms, opts)))
	r.PathPrefix("/admin/system/permission").Handler(http.StripPrefix("/admin/system/permission", syspermission.MakeHTTPHandler(sysPermissionSvc, tokenEms, opts)))

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
		//go func() {
		//_, _ = apiSvc.Alarm().Alert(context.Background(), cf.GetInt("service", "alarm.appId"), "服务它又起来了...")
		//apiSvc.Alert().Send(alert.AlertInfo, "服务又活了...")
		//}()
		// 注册到consul
		//go registerConsul()
		return http.ListenAndServe(httpAddr, nil)
	}, func(e error) {
		_ = level.Error(httpLogger).Log("transport", "HTTP", "httpListener.Close", "http", "err", e.Error())
		//apiSvc.Alert().Send(alert.AlertWarning, fmt.Sprintf("msg: %s, err: %s", "服务它停了,你猜它是不是挂了...", e.Error()))
		//_, _ = apiSvc.Alarm().Alert(context.Background(), cf.GetInt("service", "alarm.appId"), fmt.Sprintf("msg: %s, err: %s", "服务它停了,你猜它是不是挂了...", e.Error()))
		_ = level.Debug(logger).Log("db", "close", "err", db.Close())
		os.Exit(1)
	})
}

// gRPC server
func initGRPCHandler(g *group.Group) {
	//gRpcLogger := log.With(logger, "component", "gRPC")
	//
	//grpcOpts := []kitgrpc.ServerOption{
	//	kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(gRpcLogger)),
	//	kitgrpc.ServerBefore(func(ctx context.Context, mds metadata.MD) context.Context {
	//		ctx = context.WithValue(ctx, logging.TraceId, uuid.New().String())
	//		return ctx
	//	}),
	//}
	//
	//grpcListener, err := net.Listen("tcp", grpcAddr)
	//if err != nil {
	//	_ = logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	//}
	//
	//g.Add(func() error {
	//	_ = level.Debug(logger).Log("transport", "GRPC", "addr", grpcAddr)
	//	baseServer := googlegrpc.NewServer()
	//	accountPb.RegisterAccountServer(baseServer, account.MakeGRPCHandler(accountSvc, gRpcLogger, grpcOpts...))
	//	return baseServer.Serve(grpcListener)
	//}, func(error) {
	//	_ = level.Error(logger).Log("transport", "GRPC", "grpcListener.Close", grpcListener.Close())
	//})
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(err error) {
		close(cancelInterrupt)
	})
}

func registerConsul() {
	config := consulapi.DefaultConfig()
	config.Address = cf.GetString("service", "consul.host")
	config.Token = cf.GetString("service", "consul.token")
	client, err := consulapi.NewClient(config)
	if err != nil {
		_ = level.Error(logger).Log("consulapi", "NewClient", "err", err.Error())
		return
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = uuid.New().String()
	registration.Name = appName + "." + namespace
	registration.Port = 8080
	registration.Tags = []string{appName, namespace, "golang", "activity"}
	registration.Address = getLocalAddr()

	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/health", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	if err = client.Agent().ServiceRegister(registration); err != nil {
		_ = level.Error(logger).Log("client.Agent", "ServiceRegister", "err", err.Error())
	}
	//client.Agent().ServiceDeregister()
	var lastIndex uint64
	services, metainfo, err := client.Health().Service(appName, "v1000", true, &consulapi.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		logrus.Warn("error retrieving instances from Consul: %v", err)
	}
	lastIndex = metainfo.LastIndex

	addrs := map[string]struct{}{}
	for _, service := range services {
		fmt.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port)
		addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
	}
}

func getLocalAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}
