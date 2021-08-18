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
	"github.com/kplcloud/kplcloud/src/pkg/configmap"
	"github.com/kplcloud/kplcloud/src/pkg/deployment"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/config"
	kitcache "github.com/icowan/kit-cache"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"

	"github.com/kplcloud/kplcloud/src/api"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/cluster"
	pkgNs "github.com/kplcloud/kplcloud/src/pkg/namespace"
	"github.com/kplcloud/kplcloud/src/pkg/nodes"
	"github.com/kplcloud/kplcloud/src/pkg/syspermission"
	"github.com/kplcloud/kplcloud/src/pkg/sysrole"
	"github.com/kplcloud/kplcloud/src/pkg/sysuser"
	"github.com/kplcloud/kplcloud/src/redis"
	"github.com/kplcloud/kplcloud/src/repository"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Example: `## 启动命令
kplcloud start -p :8080 -g :8082
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

	//authSvc          auth.Service
	//accountSvc       account.Service

	sysUserSvc       sysuser.Service
	sysRoleSvc       sysrole.Service
	sysPermissionSvc syspermission.Service

	clusterSvc    cluster.Service
	nodeSvc       nodes.Service
	namespaceSvc  pkgNs.Service
	deploymentSvc deployment.Service
	configMapSvc  configmap.Service
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
	//authSvc = auth.New(logger, logging.TraceId, cf, store, cacheSvc, apiSvc, cf.GetBool("server", "debug"))
	//authSvc = auth.NewLogging(logger, logging.TraceId)(authSvc)
	//
	//// 用户信息模块
	//accountSvc = account.New(logger, logging.TraceId, store)
	//accountSvc = account.NewLogging(logger, logging.TraceId)(accountSvc)

	// 系统模块
	// 系统用户
	sysUserSvc = sysuser.New(logger, logging.TraceId, store)
	sysUserSvc = sysuser.NewLogging(logger, logging.TraceId)(sysUserSvc)
	// 系统角色
	sysRoleSvc = sysrole.New(logger, logging.TraceId, store)
	sysRoleSvc = sysrole.NewLogging(logger, logging.TraceId)(sysRoleSvc)
	// 用户信息模块
	sysPermissionSvc = syspermission.New(logger, logging.TraceId, store)
	sysPermissionSvc = syspermission.NewLogging(logger, logging.TraceId)(sysPermissionSvc)

	clusterSvc = cluster.New(logger, logging.TraceId, store, k8sClient)
	clusterSvc = cluster.NewLogging(logger, logging.TraceId)(clusterSvc)
	nodeSvc = nodes.New(logger, logging.TraceId, k8sClient, store)
	nodeSvc = nodes.NewLogging(logger, logging.TraceId)(nodeSvc)
	namespaceSvc = pkgNs.New(logger, logging.TraceId, k8sClient, store)
	namespaceSvc = pkgNs.NewLogging(logger, logging.TraceId)(namespaceSvc)
	deploymentSvc = deployment.New(logger, logging.TraceId, k8sClient, store)
	deploymentSvc = deployment.NewLogging(logger, logging.TraceId)(deploymentSvc)
	configMapSvc = configmap.New(logger, logging.TraceId, store, k8sClient)
	configMapSvc = configmap.NewLogging(logger, logging.TraceId)(configMapSvc)

	if tracer != nil {
		//authSvc = auth.NewTracing(tracer)(authSvc)
		sysUserSvc = sysuser.NewTracing(tracer)(sysUserSvc)
		sysRoleSvc = sysrole.NewTracing(tracer)(sysRoleSvc)
		clusterSvc = cluster.NewTracing(tracer)(clusterSvc)
		nodeSvc = nodes.NewTracing(tracer)(nodeSvc)
		namespaceSvc = pkgNs.NewTracing(tracer)(namespaceSvc)
		deploymentSvc = deployment.NewTracing(tracer)(deploymentSvc)
		configMapSvc = configmap.NewTracing(tracer)(configMapSvc)
	}

	g := &group.Group{}

	initHttpHandler(g)
	initGRPCHandler(g)
	initCancelInterrupt(g)

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
			vars := mux.Vars(request)
			clusterName, ok := vars["cluster"]
			if !ok {
				clusterName = request.Header.Get("Cluster")
			}
			ns, ok := vars["namespace"]
			if !ok {
				ns = request.Header.Get("Namespace")
			}
			ctx = context.WithValue(ctx, logging.TraceId, guid)
			ctx = context.WithValue(ctx, "token-context", token)
			ctx = context.WithValue(ctx, middleware.ContextKeyClusterName, clusterName)
			ctx = context.WithValue(ctx, middleware.ContextKeyNamespaceName, ns)
			return ctx
		}),
		kithttp.ServerBefore(middleware.TracingServerBefore(tracer)),
	}

	ems := []endpoint.Middleware{
		middleware.ClusterMiddleware(store),  //2
		middleware.TracingMiddleware(tracer), // 1
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 0
	}

	tokenEms := []endpoint.Middleware{
		//middleware.CheckAuthMiddleware(logger, cacheSvc, tracer), // 3
	}
	tokenEms = append(tokenEms, ems...)

	r := mux.NewRouter()

	// 授权登录模块
	//r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", auth.MakeHTTPHandler(authSvc, ems, opts)))
	//r.PathPrefix("/account").Handler(http.StripPrefix("/account", account.MakeHTTPHandler(accountSvc, tokenEms, opts)))

	r.PathPrefix("/cluster").Handler(http.StripPrefix("/cluster", cluster.MakeHTTPHandler(clusterSvc, tokenEms, opts)))
	r.PathPrefix("/node").Handler(http.StripPrefix("/node", nodes.MakeHTTPHandler(nodeSvc, tokenEms, opts)))
	r.PathPrefix("/namespace").Handler(http.StripPrefix("/namespace", pkgNs.MakeHTTPHandler(namespaceSvc, tokenEms, opts)))
	r.PathPrefix("/deployment").Handler(http.StripPrefix("/deployment", deployment.MakeHTTPHandler(deploymentSvc, tokenEms, opts)))
	r.PathPrefix("/configmap").Handler(http.StripPrefix("/configmap", configmap.MakeHTTPHandler(configMapSvc, tokenEms, opts)))

	// 以下为系统模块
	// 系统用户模块
	r.PathPrefix("/system/user").Handler(http.StripPrefix("/system/user", sysuser.MakeHTTPHandler(sysUserSvc, tokenEms, opts)))
	// 系统角色、权限
	r.PathPrefix("/system/role").Handler(http.StripPrefix("/system/role", sysrole.MakeHTTPHandler(sysRoleSvc, tokenEms, opts)))
	r.PathPrefix("/system/permission").Handler(http.StripPrefix("/system/permission", syspermission.MakeHTTPHandler(sysPermissionSvc, tokenEms, opts)))

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

func prepare() error {
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
		// 实例化仓库
		store = repository.New(db, logger, "traceId", tracer, rds, cacheSvc)
	}

	ctx := context.Background()

	// 读取所有配置
	settings, err := store.SysSetting().FindAll(ctx)
	if err != nil {
		_ = level.Error(logger).Log("store.SysSetting", "FindAll", "err", err.Error())
		return err
	}

	for _, v := range settings {
		cf.SetValue(v.Section, v.Key, v.Value)
	}

	db.LogMode(cf.GetBool("server", "debug"))
	logger = logging.SetLogging(logger, cf)

	//hashId = hashids.New("", cf.GetString(config.SectionServer, "app.key"), 12)

	// 链路追踪
	tracer, _, err = newJaegerTracer(cf)

	if err != nil {
		_ = level.Error(logger).Log("jaegerTracer", "connect", "err", err.Error())
	}
	// 实例化redis
	rds, err = redis.New(cf.GetString(config.SectionRedis, "hosts"),
		cf.GetString(config.SectionRedis, "password"),
		cf.GetString(config.SectionRedis, "prefix"),
		cf.GetInt(config.SectionRedis, "db"), tracer)
	if err != nil {
		_ = level.Error(logger).Log("redis", "connect", "err", err.Error())
	}
	_ = level.Info(logger).Log("rds", "connect", "success", true)

	// 实例化cache
	cacheSvc = kitcache.New(logger, logging.TraceId, rds)
	cacheSvc = kitcache.NewLoggingServer(logger, cacheSvc, logging.TraceId)

	// 实例化外部API
	apiSvc = api.NewApi(logger, logging.TraceId, tracer, cf, cacheSvc)

	store = repository.New(db, logger, "traceId", tracer, rds, cacheSvc)

	// 实例化k8s client
	k8sClient, err = kubernetes.NewClient(store)
	if err != nil {
		_ = level.Error(logger).Log("kubernetes", "NewClient", "err", err.Error())
	} else {
		k8sClient = kubernetes.NewLogging(logger, logging.TraceId)(k8sClient)
		if tracer != nil {
			k8sClient = kubernetes.NewTracing(tracer)(k8sClient)
		}
	}

	return nil
}
