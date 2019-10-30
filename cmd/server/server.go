package server

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/transport/amqp"
	"github.com/jinzhu/gorm"
	kplamqp "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/casbin"
	"github.com/kplcloud/kplcloud/src/cmd"
	"github.com/kplcloud/kplcloud/src/config"
	"github.com/kplcloud/kplcloud/src/email"
	"github.com/kplcloud/kplcloud/src/git-repo"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/mysql"
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
	"github.com/kplcloud/kplcloud/src/pkg/wechat"
	"github.com/kplcloud/kplcloud/src/pkg/workspace"
	"github.com/kplcloud/kplcloud/src/redis"
	"github.com/kplcloud/kplcloud/src/repository"
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

	cmd.AddFlags(rootCmd)
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

	logger = logging.SetLogging(logger, cf)

	// mysql client
	db, err = mysql.NewDb(cf)
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		panic(err)
	}

	store = repository.NewRepository(db)
	if err = importToDb(cf.GetString("server", "app_key")); err != nil {
		_ = level.Error(logger).Log("import", "database", "err", err)
		panic(err)
	}

	// redis client
	rds, err := redis.NewRedisClient(cf)
	if err != nil {
		_ = level.Error(logger).Log("redis", "connect", "err", err)
		panic(err)
	}

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
	//istioClient, err := istio.NewClient(cf)
	//if err != nil {
	//	_ = logger.Log("istio", "NewClient", "err", err)
	//	panic(err)
	//}

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
		templateSvc    = template.NewService(logger, store)
		groupSvc       = group.NewService(logger, cf, store)
		proclaimSvc    = proclaim.NewService(logger, cf, amqpClient, store)
		accountSvc     = account.NewService(logger, cf, store)
		hookSvc        = hooks.NewService(logger, store, hookQueueSvc)
		projectSvc     = project.NewService(logger, cf, rds, k8sClient, amqpClient, jenkinsClient, store, hookQueueSvc)
		publicSvc      = public.NewService(logger, cf, amqpClient, k8sClient, jenkinsClient, store)
		buildSvc       = build.NewService(logger, jenkinsClient, amqpClient, k8sClient, cf, store, hookQueueSvc)
		gitSvc         = git.NewService(logger, cf, gitClient, store)
		weixinSvc      = wechat.NewService(logger, cf, wxClient, store)
		permissionSvc  = permission.NewService(logger, casbinClient, store)
		msgAlarmSvc    = msgs.NewServiceAlarm(logger, cf, mailClient, amqpClient, store)
		msgNoticeSvc   = msgs.NewServiceNotice(logger, cf, mailClient, amqpClient, store)
		msgProclaimSvc = msgs.NewServiceProclaim(logger, cf, mailClient, amqpClient, store)
		msgWechatSvc   = msgs.NewServiceWechatQueue(logger, cf, amqpClient, wxClient, store)
		roleSvc        = role.NewService(logger, casbinClient, store)
		memberSvc      = member.NewService(logger, cf, casbinClient, store)
		consulSvc      = consul.NewService(logger, cf, store)
		eventSvc       = event.NewService(logger, store)
		workspaceSvc   = workspace.NewService(logger, cf, store.Build())
		marketSvc      = market.NewService(logger, store)
		monitorSvc     = monitor.NewService(logger, cf, k8sClient, store)
		auditSvc       = audit.NewService(logger, cf, jenkinsClient, k8sClient, amqpClient, store, hookQueueSvc, buildSvc)
		statisticsSvc  = statistics.NewService(logger, cf, store)
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
		mux := http.NewServeMux()

		mux.Handle("/auth/", auth.MakeHandler(authSvc, httpLogger))

		// k8s rds
		mux.Handle("/namespace", namespace.MakeHandler(namespaceSvc, httpLogger))
		mux.Handle("/namespace/", namespace.MakeHandler(namespaceSvc, httpLogger))
		mux.Handle("/storageclass", storage.MakeHandler(storageSvc, httpLogger))
		mux.Handle("/storageclass/", storage.MakeHandler(storageSvc, httpLogger))
		mux.Handle("/persistentvolumeclaim/", persistentvolumeclaim.MakeHandler(pvcSvc, httpLogger))
		mux.Handle("/persistentvolume/", persistentvolume.MakeHandler(pvSvc, httpLogger))
		mux.Handle("/ws/pods/console/exec/", sockjs.NewHandler("/ws/pods/console/exec", sockjs.DefaultOptions, func(session sockjs.Session) {
			termilanSvc.HandleTerminalSession(session)
		}))
		mux.Handle("/terminal/", terminal.MakeHandler(termilanSvc, httpLogger, store))
		mux.Handle("/deployment/", deployment.MakeHandler(deploymentSvc, httpLogger, cf, store))
		mux.Handle("/pods/", pod.MakeHandler(podsSvc, httpLogger, store))
		mux.Handle("/ingress/", ingress.MakeHandler(ingressSvc, httpLogger, store))
		mux.Handle("/config/", configmap.MakeHandler(configMapSvc, logger, store))
		mux.Handle("/configmap/", configmap.MakeHandler(configMapSvc, logger, store))
		mux.Handle("/discovery/", discovery.MakeHandler(discoverySvc, logger))
		mux.Handle("/tools/", tools.MakeHandler(toolsSvc, logger, store))
		mux.Handle("/cronjob/", cronjob.MakeHandler(cronjobSvc, httpLogger, store))

		// project
		mux.Handle("/project/", project.MakeHandler(projectSvc, httpLogger, store))
		mux.Handle("/template/", template.MakeHandler(templateSvc, httpLogger))
		mux.Handle("/template", template.MakeHandler(templateSvc, httpLogger))
		mux.Handle("/group/", group.MakeHandler(groupSvc, httpLogger, store.Groups()))
		mux.Handle("/group", group.MakeHandler(groupSvc, httpLogger, store.Groups()))
		mux.Handle("/proclaim/", proclaim.MakeHandler(proclaimSvc, httpLogger))
		mux.Handle("/proclaim", proclaim.MakeHandler(proclaimSvc, httpLogger))
		mux.Handle("/notice/", notice.MakeHandler(noticeSvc, httpLogger))
		mux.Handle("/notice", notice.MakeHandler(noticeSvc, httpLogger))
		mux.Handle("/hooks/", hooks.MakeHandler(hookSvc, httpLogger, store))
		mux.Handle("/git/", git.MakeHandler(gitSvc, httpLogger, store))
		mux.Handle("/build/", build.MakeHandler(buildSvc, httpLogger, store))
		mux.Handle("/permission/", permission.MakeHandler(permissionSvc, httpLogger))
		mux.Handle("/role", role.MakeHandler(roleSvc, httpLogger))
		mux.Handle("/role/", role.MakeHandler(roleSvc, httpLogger))
		mux.Handle("/member", member.MakeHandler(memberSvc, httpLogger))
		mux.Handle("/member/", member.MakeHandler(memberSvc, httpLogger))
		mux.Handle("/event/", event.MakeHandler(eventSvc, httpLogger))
		mux.Handle("/workspace/", workspace.MakeHandler(workspaceSvc, httpLogger))
		mux.Handle("/account/", account.MakeHandler(accountSvc, httpLogger))
		mux.Handle("/account", account.MakeHandler(accountSvc, httpLogger))
		mux.Handle("/market/", market.MakeHandler(marketSvc, httpLogger))
		mux.Handle("/monitor/", monitor.MakeHandler(monitorSvc, httpLogger))
		mux.Handle("/audit/", audit.MakeHandler(auditSvc, httpLogger, store))
		mux.Handle("/statistics/", statistics.MakeHandler(statisticsSvc, httpLogger))

		// wechat receive
		mux.Handle("/wechat/", wechat.MakeHandler(weixinSvc, httpLogger))

		// public api
		mux.Handle("/public/", public.MakeHandler(publicSvc, httpLogger))
		mux.Handle("/consul/", consul.MakeHandler(consulSvc, httpLogger, cf))

		mux.Handle("/", http.FileServer(http.Dir(cf.GetString("server", "http_static"))))

		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cf.GetString("server", "http_static")))))

		handlers := make(map[string]string, 3)
		if cf.GetBool("cors", "allow") {
			handlers["Access-Control-Allow-Origin"] = cf.GetString("cors", "origin")
			handlers["Access-Control-Allow-Methods"] = cf.GetString("cors", "methods")
			handlers["Access-Control-Allow-Headers"] = cf.GetString("cors", "headers")
		}
		http.Handle("/", accessControl(mux, logger, handlers))
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
