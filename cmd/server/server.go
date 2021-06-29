package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/transport/amqp"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/config"
	mysqlclient "github.com/icowan/mysql-client"
	redisclient "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	kplamqp "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/casbin"
	"github.com/kplcloud/kplcloud/src/email"
	"github.com/kplcloud/kplcloud/src/git-repo"
	"github.com/kplcloud/kplcloud/src/istio"
	"github.com/kplcloud/kplcloud/src/jenkins"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/account"
	"github.com/kplcloud/kplcloud/src/pkg/audit"
	"github.com/kplcloud/kplcloud/src/pkg/auth"
	"github.com/kplcloud/kplcloud/src/pkg/build"
	"github.com/kplcloud/kplcloud/src/pkg/configmap"
	"github.com/kplcloud/kplcloud/src/pkg/consul"
	"github.com/kplcloud/kplcloud/src/pkg/cronjob"
	"github.com/kplcloud/kplcloud/src/pkg/deployment"
	"github.com/kplcloud/kplcloud/src/pkg/discovery"
	"github.com/kplcloud/kplcloud/src/pkg/event"
	"github.com/kplcloud/kplcloud/src/pkg/git"
	"github.com/kplcloud/kplcloud/src/pkg/group"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/pkg/ingress"
	"github.com/kplcloud/kplcloud/src/pkg/market"
	"github.com/kplcloud/kplcloud/src/pkg/member"
	"github.com/kplcloud/kplcloud/src/pkg/monitor"
	"github.com/kplcloud/kplcloud/src/pkg/msgs"
	"github.com/kplcloud/kplcloud/src/pkg/namespace"
	"github.com/kplcloud/kplcloud/src/pkg/nodes"
	"github.com/kplcloud/kplcloud/src/pkg/notice"
	"github.com/kplcloud/kplcloud/src/pkg/permission"
	"github.com/kplcloud/kplcloud/src/pkg/persistentvolume"
	"github.com/kplcloud/kplcloud/src/pkg/persistentvolumeclaim"
	"github.com/kplcloud/kplcloud/src/pkg/pod"
	"github.com/kplcloud/kplcloud/src/pkg/proclaim"
	"github.com/kplcloud/kplcloud/src/pkg/project"
	"github.com/kplcloud/kplcloud/src/pkg/public"
	"github.com/kplcloud/kplcloud/src/pkg/role"
	"github.com/kplcloud/kplcloud/src/pkg/statistics"
	"github.com/kplcloud/kplcloud/src/pkg/storage"
	"github.com/kplcloud/kplcloud/src/pkg/template"
	"github.com/kplcloud/kplcloud/src/pkg/terminal"
	"github.com/kplcloud/kplcloud/src/pkg/tools"
	"github.com/kplcloud/kplcloud/src/pkg/virtualservice"
	"github.com/kplcloud/kplcloud/src/pkg/wechat"
	"github.com/kplcloud/kplcloud/src/pkg/workspace"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	wechatsdk "github.com/yijizhichang/wechat-sdk"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultHttpPort      = ":8080"
	DefaultConfigPath    = "app.cfg"
	DefaultKubeConfig    = "config.yaml"
	DefaultAdminEmail    = "kplcloud@nsini.com"
	DefaultAdminPassword = "admin"
	DefaultInitDBSQL     = "./database/kplcloud.sql"
)

var (
	httpAddr      = envString("HTTP_ADDR", DefaultHttpPort)
	configPath    = envString("CONFIG_PATH", DefaultConfigPath)
	kubeConfig    = envString("KUBE_CONFIG", DefaultKubeConfig)
	adminEmail    = envString("ADMIN_EMAIL", DefaultAdminEmail)
	adminPassword = envString("ADMIN_PASSWORD", DefaultAdminPassword)
	sqlPath       = envString("INIT_SQL", DefaultInitDBSQL)

	logger log.Logger
	store  repository.Repository
	db     *gorm.DB
	cf     *config.Config
	err    error
	rds    redisclient.RedisClient
)

var (
	rootCmd = &cobra.Command{
		Use:               "server",
		Short:             "开普勒平台服务端",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# 开普勒平台服务端

您可以通过改命令来启动您的服务

可用的配置类型：

[start]

有关开普勒平台的相关概述，请参阅 https://github.com/kplcloud/kplcloud
`,
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务",
		Example: `## 启动命令
server start -p :8080 -c ./app.cfg -k ./k8s-config.yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			run()
			return nil
		},
	}
)

func init() {

	rootCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的端口: :8080")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", DefaultConfigPath, "配置文件路径: ./app.cfg")
	startCmd.PersistentFlags().StringVarP(&kubeConfig, "kube.config", "k", DefaultKubeConfig, "kubernetes config文件: ./config.yaml")
	startCmd.PersistentFlags().StringVar(&adminEmail, "admin.email", DefaultAdminEmail, "初始化管理员邮箱")
	startCmd.PersistentFlags().StringVar(&adminPassword, "admin.password", DefaultAdminPassword, "初始化管理员密码")
	startCmd.PersistentFlags().StringVar(&sqlPath, "init.sqlpath", DefaultInitDBSQL, "初始化sql文件")

	addFlags(rootCmd)
	rootCmd.AddCommand(startCmd)
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func run() {

	cf, err = config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	_ = cf.SetValue("server", "kube_config", kubeConfig)

	logger = logging.SetLogging(logger, cf.GetString("server", "log_path"), cf.GetString("server", "log_level"))

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		cf.GetString(config.SectionMysql, "user"),
		cf.GetString(config.SectionMysql, "password"),
		cf.GetString(config.SectionMysql, "host"),
		cf.GetString(config.SectionMysql, "port"),
		cf.GetString(config.SectionMysql, "database"))

	// 连接数据库
	db, err = mysqlclient.NewMysql(dbUrl, cf.GetBool(config.SectionServer, "debug"))
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		panic(err)
	}

	store = repository.NewRepository(db, logger, logging.TraceId)
	if err = importToDb(cf.GetString("server", "app_key")); err != nil {
		_ = level.Error(logger).Log("import", "database", "err", err)
		panic(err)
	}

	// 连接redis
	rds, err = redisclient.NewRedisClient(cf.GetString(config.SectionRedis, "hosts"),
		cf.GetString(config.SectionRedis, "password"),
		cf.GetString(config.SectionRedis, "prefix"),
		cf.GetInt(config.SectionRedis, "db"))

	if err != nil {
		_ = level.Error(logger).Log("redis", "connect", "err", err)
		panic(err)
	}

	_ = level.Info(logger).Log("rds", "connect", "success", true)
	// amqp client
	amqpClient, err := kplamqp.NewAmqp(cf)
	if err != nil {
		_ = level.Error(logger).Log("kplamqp", "NewAmqp", "err", err)
		panic(err)
	}

	defer func() {
		_ = db.Close()         // close db
		_ = rds.Close()        // close redis
		_ = amqpClient.Close() // close amqp
	}()

	// k8s client
	k8sClient, err := kubernetes.NewClient(cf)
	if err != nil {
		_ = level.Error(logger).Log("kubernetes", "NewClient", "err", err)
		panic(err)
	}

	// istio client
	istioClient, err := istio.NewClient(k8sClient.Config())
	if err != nil {
		_ = logger.Log("istio", "NewClient", "err", err)
		panic(err)
	}

	// jenkins client
	var jenkinsClient jenkins.Jenkins
	if cf.GetString("jenkins", "host") != "" {
		jenkinsClient, err = jenkins.NewJenkins(cf)
		if err != nil {
			_ = level.Error(logger).Log("jenkins", "NewJenkins", "err", err)
			panic(err)
		}
	}

	// casbin
	casbinClient, err := casbin.NewCasbin(cf, db, rds)
	if err != nil {
		_ = level.Error(logger).Log("casbin", "NewCasbin", "err", err.Error())
		panic(err)
	}

	// git client
	gitClient := git_repo.NewClient(cf)

	//wechat client
	wxClient := wechatsdk.NewWechat(&wechatsdk.Config{
		AppID:            cf.GetString("wechat", "app_id"),
		AppSecret:        cf.GetString("wechat", "app_secret"),
		Token:            cf.GetString("wechat", "token"),
		EncodingAESKey:   cf.GetString("wechat", "encoding_aes_key"),
		Cache:            rds,
		ThirdAccessToken: false,
		ProxyUrl:         "", //代理地址
	})

	// email client
	mailClient := email.NewEmail(cf)

	fieldKeys := []string{"method"}

	var (
		authSvc      = auth.NewService(logger, cf, casbinClient, store)
		noticeSvc    = notice.NewService(logger, cf, amqpClient, store)
		hookQueueSvc = hooks.NewServiceHookQueue(logger, amqpClient, cf, store, noticeSvc)

		// start k8s rds
		namespaceSvc  = namespace.NewService(logger, cf, k8sClient, store)
		storageSvc    = storage.NewService(logger, k8sClient, store)
		pvcSvc        = persistentvolumeclaim.NewService(logger, k8sClient, store)
		pvSvc         = persistentvolume.NewService(logger, k8sClient)
		deploymentSvc = deployment.NewService(logger, k8sClient, store, hookQueueSvc)
		podsSvc       = pod.NewService(logger, k8sClient, cf, hookQueueSvc)
		termilanSvc   = terminal.NewService(logger, cf, k8sClient)
		ingressSvc    = ingress.NewService(logger, cf, k8sClient, store, hookQueueSvc)
		configMapSvc  = configmap.NewService(logger, cf, jenkinsClient, k8sClient, store)
		discoverySvc  = discovery.NewService(logger, k8sClient, cf, store)
		toolsSvc      = tools.NewService(logger, cf, jenkinsClient, k8sClient, store)
		cronjobSvc    = cronjob.NewService(logger, cf, jenkinsClient, k8sClient, amqpClient, store)
		// end k8s rds
		templateSvc       = template.NewService(logger, store)
		groupSvc          = group.NewService(logger, cf, store)
		proclaimSvc       = proclaim.NewService(logger, cf, amqpClient, store)
		accountSvc        = account.NewService(logger, cf, store)
		hookSvc           = hooks.NewService(logger, store, hookQueueSvc)
		projectSvc        = project.NewService(logger, cf, rds, k8sClient, amqpClient, jenkinsClient, store, hookQueueSvc)
		publicSvc         = public.NewService(logger, cf, amqpClient, k8sClient, jenkinsClient, store)
		buildSvc          = build.NewService(logger, jenkinsClient, amqpClient, k8sClient, cf, store, hookQueueSvc)
		gitSvc            = git.NewService(logger, cf, gitClient, store)
		weixinSvc         = wechat.NewService(logger, cf, wxClient, store)
		permissionSvc     = permission.NewService(logger, casbinClient, store)
		msgAlarmSvc       = msgs.NewServiceAlarm(logger, cf, mailClient, amqpClient, store)
		msgNoticeSvc      = msgs.NewServiceNotice(logger, cf, mailClient, amqpClient, store)
		msgProclaimSvc    = msgs.NewServiceProclaim(logger, cf, mailClient, amqpClient, store)
		msgWechatSvc      = msgs.NewServiceWechatQueue(logger, cf, amqpClient, wxClient, store)
		roleSvc           = role.NewService(logger, casbinClient, store)
		memberSvc         = member.NewService(logger, cf, casbinClient, store)
		consulSvc         = consul.NewService(logger, cf, store)
		eventSvc          = event.NewService(logger, store)
		workspaceSvc      = workspace.NewService(logger, cf, store.Build())
		marketSvc         = market.NewService(logger, store)
		monitorSvc        = monitor.NewService(logger, cf, k8sClient, store)
		auditSvc          = audit.NewService(logger, cf, jenkinsClient, k8sClient, amqpClient, store, hookQueueSvc, buildSvc)
		statisticsSvc     = statistics.NewService(logger, cf, store)
		virtualServiceSvc = virtualservice.NewService(logger, istioClient)
		nodesSvc          = nodes.New(logger, "traceId", store, k8sClient)
	)
	// namespace service

	authSvc = auth.NewLoggingService(logger, authSvc)
	authSvc = auth.NewInstrumentingService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "auth",
			Subsystem: "auth_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "auth",
			Subsystem: "auth_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), authSvc)

	// k8s rds
	namespaceSvc = namespace.NewLoggingService(logger, namespaceSvc) // 日志
	storageSvc = storage.NewLoggingService(logger, storageSvc)
	pvcSvc = persistentvolumeclaim.NewLoggingService(logger, pvcSvc)
	pvSvc = persistentvolume.NewLoggingService(logger, pvSvc)
	deploymentSvc = deployment.NewLoggingService(logger, deploymentSvc)
	podsSvc = pod.NewLoggingService(logger, podsSvc)
	termilanSvc = terminal.NewLoggingService(logger, termilanSvc)
	ingressSvc = ingress.NewLoggingService(logger, ingressSvc)
	configMapSvc = configmap.NewLoggingService(logger, configMapSvc)
	discoverySvc = discovery.NewLoggingService(logger, discoverySvc)
	toolsSvc = tools.NewLoggingService(logger, toolsSvc)
	// end k8s rds

	templateSvc = template.NewLoggingService(logger, templateSvc)

	groupSvc = group.NewLoggingService(logger, groupSvc)
	proclaimSvc = proclaim.NewLoggingService(logger, proclaimSvc)
	noticeSvc = notice.NewLoggingService(logger, noticeSvc)
	accountSvc = account.NewLoggingService(logger, accountSvc)
	hookSvc = hooks.NewLoggingService(logger, hookSvc)
	projectSvc = project.NewLoggingService(logger, projectSvc)
	gitSvc = git.NewLoggingService(logger, gitSvc)
	cronjobSvc = cronjob.NewLoggingService(logger, cronjobSvc)
	weixinSvc = wechat.NewLoggingService(logger, weixinSvc)
	permissionSvc = permission.NewLoggingService(logger, permissionSvc)
	roleSvc = role.NewLoggingService(logger, roleSvc)
	memberSvc = member.NewLoggingService(logger, memberSvc)
	consulSvc = consul.NewLoggingService(logger, consulSvc)
	marketSvc = market.NewLoggingService(logger, marketSvc)
	monitorSvc = monitor.NewLoggingService(logger, monitorSvc)
	auditSvc = audit.NewLoggingService(logger, auditSvc)
	statisticsSvc = statistics.NewLoggingService(logger, statisticsSvc)
	buildSvc = build.NewLoggingService(logger, buildSvc)

	publicSvc = public.NewLoggingService(logger, publicSvc)
	publicSvc = public.NewInstrumentingService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "public",
			Subsystem: "public_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "public",
			Subsystem: "public_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), publicSvc)

	httpLogger := log.With(logger, "component", "http")
	{

		opts := []kithttp.ServerOption{
			kithttp.ServerErrorLogger(logger),
			kithttp.ServerErrorEncoder(encode.EncodeError),
			kithttp.ServerBefore(kithttp.PopulateRequestContext),
			kithttp.ServerBefore(kitjwt.HTTPToContext()),
			kithttp.ServerBefore(middleware.NamespaceToContext()),
			kithttp.ServerBefore(middleware.CasbinToContext()),
		}

		ems := []endpoint.Middleware{
			middleware.NamespaceMiddleware(logger), // 3
			middleware.CheckAuthMiddleware(logger),
			kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory),
		}

		//mux := http.NewServeMux()

		r := mux.NewRouter()

		r.Handle("/auth/", auth.MakeHandler(authSvc, httpLogger))

		// k8s rds
		r.Handle("/namespace", namespace.MakeHandler(namespaceSvc, httpLogger))
		r.Handle("/namespace/", namespace.MakeHandler(namespaceSvc, httpLogger))
		r.Handle("/storageclass", storage.MakeHandler(storageSvc, httpLogger))
		r.Handle("/storageclass/", storage.MakeHandler(storageSvc, httpLogger))
		r.Handle("/persistentvolumeclaim/", persistentvolumeclaim.MakeHandler(pvcSvc, httpLogger))
		r.Handle("/persistentvolume/", persistentvolume.MakeHandler(pvSvc, httpLogger))
		r.Handle("/ws/pods/console/exec/", sockjs.NewHandler("/ws/pods/console/exec", sockjs.DefaultOptions, func(session sockjs.Session) {
			termilanSvc.HandleTerminalSession(session)
		}))
		r.Handle("/terminal/", terminal.MakeHandler(termilanSvc, httpLogger, store))
		r.Handle("/deployment/", deployment.MakeHandler(deploymentSvc, httpLogger, cf, store))
		r.Handle("/pods/", pod.MakeHandler(podsSvc, httpLogger, store))
		r.Handle("/ingress/", ingress.MakeHandler(ingressSvc, httpLogger, store))
		r.Handle("/config/", configmap.MakeHandler(configMapSvc, logger, store))
		r.Handle("/configmap/", configmap.MakeHandler(configMapSvc, logger, store))
		r.Handle("/discovery/", discovery.MakeHandler(discoverySvc, logger))
		r.Handle("/tools/", tools.MakeHandler(toolsSvc, logger, store))
		r.Handle("/cronjob/", cronjob.MakeHandler(cronjobSvc, httpLogger, store, opts, ems))

		r.Handle("/virtualservice/", virtualservice.MakeHTTPHandler(virtualServiceSvc, logger, opts, ems))

		// project
		r.PathPrefix("/project").Handler(http.StripPrefix("/project", project.MakeHandler(projectSvc, httpLogger, store)))
		r.PathPrefix("/template").Handler(http.StripPrefix("/template", template.MakeHandler(templateSvc, httpLogger)))
		r.PathPrefix("/group").Handler(http.StripPrefix("/group", group.MakeHandler(groupSvc, httpLogger, store.Groups())))
		r.PathPrefix("/group/").Handler(http.StripPrefix("/group/", group.MakeHandler(groupSvc, httpLogger, store.Groups())))
		r.PathPrefix("/proclaim").Handler(http.StripPrefix("/proclaim", proclaim.MakeHandler(proclaimSvc, httpLogger)))
		r.PathPrefix("/proclaim/").Handler(http.StripPrefix("/proclaim/", proclaim.MakeHandler(proclaimSvc, httpLogger)))
		r.PathPrefix("/notice").Handler(http.StripPrefix("/notice", notice.MakeHandler(noticeSvc, httpLogger)))
		r.PathPrefix("/hooks").Handler(http.StripPrefix("/hooks", hooks.MakeHandler(hookSvc, httpLogger, store)))
		r.PathPrefix("/git").Handler(http.StripPrefix("/git", git.MakeHandler(gitSvc, httpLogger, store)))
		r.PathPrefix("/build").Handler(http.StripPrefix("/build", build.MakeHandler(buildSvc, httpLogger, store)))
		r.PathPrefix("/permission").Handler(http.StripPrefix("/permission", permission.MakeHandler(permissionSvc, httpLogger)))
		r.PathPrefix("/role/").Handler(http.StripPrefix("/role/", role.MakeHandler(roleSvc, httpLogger)))
		r.PathPrefix("/role").Handler(http.StripPrefix("/role", role.MakeHandler(roleSvc, httpLogger)))
		r.PathPrefix("/member").Handler(http.StripPrefix("/member", member.MakeHandler(memberSvc, httpLogger)))
		r.PathPrefix("/member/").Handler(http.StripPrefix("/member/", member.MakeHandler(memberSvc, httpLogger)))
		r.PathPrefix("/event").Handler(http.StripPrefix("/event", event.MakeHandler(eventSvc, httpLogger)))
		r.PathPrefix("/workspace").Handler(http.StripPrefix("/workspace", workspace.MakeHandler(workspaceSvc, httpLogger)))
		r.PathPrefix("/account").Handler(http.StripPrefix("/account", account.MakeHandler(accountSvc, httpLogger)))
		r.PathPrefix("/market").Handler(http.StripPrefix("/market", market.MakeHandler(marketSvc, httpLogger)))
		r.PathPrefix("/monitor").Handler(http.StripPrefix("/monitor", monitor.MakeHandler(monitorSvc, httpLogger)))
		r.PathPrefix("/audit").Handler(http.StripPrefix("/audit", audit.MakeHandler(auditSvc, httpLogger, store)))
		r.PathPrefix("/statistics").Handler(http.StripPrefix("/statistics", statistics.MakeHandler(statisticsSvc, httpLogger)))
		r.PathPrefix("/nodes").Handler(http.StripPrefix("/nodes", nodes.MakeHTTPHandler(nodesSvc, ems, opts)))

		// wechat receive
		r.PathPrefix("/wechat").Handler(http.StripPrefix("/wechat", wechat.MakeHandler(weixinSvc, httpLogger)))

		// public api
		r.PathPrefix("/public").Handler(http.StripPrefix("/public", public.MakeHandler(publicSvc, httpLogger)))
		r.PathPrefix("/consul").Handler(http.StripPrefix("/consul", consul.MakeHandler(consulSvc, httpLogger, cf)))

		r.PathPrefix("/").Handler(http.FileServer(http.Dir(cf.GetString("server", "http_static"))))

		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cf.GetString("server", "http_static")))))

		handlers := make(map[string]string, 3)
		if cf.GetBool("cors", "allow") {
			handlers["Access-Control-Allow-Origin"] = cf.GetString("cors", "origin")
			handlers["Access-Control-Allow-Methods"] = cf.GetString("cors", "methods")
			handlers["Access-Control-Allow-Headers"] = cf.GetString("cors", "headers")
		}
		http.Handle("/", accessControl(r, logger, handlers))
	}
	subscribeToQueue(logger, amqpClient, buildSvc, cronjobSvc, msgAlarmSvc, msgNoticeSvc, msgProclaimSvc, msgWechatSvc, hookQueueSvc)

	errs := make(chan error, 2)
	go func() {
		_ = level.Debug(logger).Log("transport", "http", "address", httpAddr, "config", configPath, "kubeconfig", kubeConfig, "msg", "listening")
		//errs <- http.ListenAndServe(httpAddr, addCors())
		errs <- http.ListenAndServe(httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = level.Error(logger).Log("terminated", <-errs)
}

func subscribeToQueue(logger log.Logger, client kplamqp.AmqpClient, buildSvc build.Service, cronSvc cronjob.Service,
	msgAlarmSvc msgs.ServiceAlarm,
	msgNoticeSvc msgs.ServiceNotice,
	msgProclaimSvc msgs.ServiceProclaim,
	msgWechatSvc msgs.ServiceWechatQueue,
	hookQueueSvc hooks.ServiceHookQueue,
) {
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.BuildTopic.String())
		go client.SubscribeToQueue(ctx, logger, buildSvc.ReceiverBuild)
	}
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.AlarmTopic.String())
		go client.SubscribeToQueue(ctx, logger, msgAlarmSvc.DistributeAlarm)
	}
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.NoticeTopic.String())
		go client.SubscribeToQueue(ctx, logger, msgNoticeSvc.DistributeNotice)
	}
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.ProclaimTopic.String())
		go client.SubscribeToQueue(ctx, logger, msgProclaimSvc.DistributeProclaim)
	}
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.MsgWechatTopic.String())
		go client.SubscribeToQueue(ctx, logger, msgWechatSvc.DistributeMsgWechat)
	}
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.CronJobTopic.String())
		go client.SubscribeToQueue(ctx, logger, cronSvc.CronJobQueuePop)
	}
	{
		ctx := context.WithValue(context.Background(), amqp.ContextKeyPublishKey, kplamqp.HookTopic.String())
		go client.SubscribeToQueue(ctx, logger, hookQueueSvc.HookReceiver)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func accessControl(h http.Handler, logger log.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		//requestId := r.Header.Get("X-Request-Id")
		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})
}

func addFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}
