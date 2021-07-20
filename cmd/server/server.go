/**
 * @Time: 2020/4/23 14:07
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	kitcache "github.com/icowan/kit-cache"
	mysqlclient "github.com/icowan/mysql-client"
	redisclient "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"

	"github.com/kplcloud/kplcloud/src/api"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/redis"
	"github.com/kplcloud/kplcloud/src/repository"
)

const (
	DefaultHttpPort   = ":8080"
	DefaultGRPCPort   = ":8082"
	DefaultConfigPath = "/usr/local/activity/etc/app.cfg"
	DefaultEnv        = "dev"
	DefaultWebPath    = "./web"
)

var (
	httpAddr, grpcAddr, configPath, env, desc string
	webPath                                   string
	logger                                    log.Logger
	db                                        *gorm.DB
	cf                                        *config.Config
	err                                       error
	store                                     repository.Repository
	appName, namespace                        string

	rootCmd = &cobra.Command{
		Use:               "kit-admin",
		Short:             "",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# 活动平台
可用的配置类型：
[start, generate, setting]
有关本系统的相关概述，请参阅 http://github.com/kplcloud/kplcloud
`,
	}
)

var (
	rateBucketNum = 5000

	rds      redisclient.RedisClient
	cacheSvc kitcache.Service
	apiSvc   api.Service
	//hashId   hashids.HashIds
)

func init() {
	startCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的http端口")
	startCmd.PersistentFlags().StringVarP(&grpcAddr, "grpc.port", "g", DefaultGRPCPort, "服务启动的gRPC端口")

	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", DefaultConfigPath, "配置文件路径")
	rootCmd.PersistentFlags().StringVarP(&env, "config.env", "e", DefaultEnv, "环境: test,prod,dev")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "app", "命名空间")
	rootCmd.PersistentFlags().StringVarP(&appName, "app.name", "a", "", "应用名称")

	settingCmd.PersistentFlags().StringVar(&desc, "desc", "", "描述")

	generateCmd.AddCommand(genTableCmd, genInitDataCmd)

	settingCmd.AddCommand(settingAddCmd, settingDelCmd, settingUpdateCmd, settingGetCmd)

	addFlags(rootCmd)
	rootCmd.AddCommand(startCmd, generateCmd, settingCmd)
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

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		cf.GetString(config.SectionMysql, "user"),
		cf.GetString(config.SectionMysql, "password"),
		cf.GetString(config.SectionMysql, "host"),
		cf.GetString(config.SectionMysql, "port"),
		cf.GetString(config.SectionMysql, "database"))

	// 连接数据库
	db, err = mysqlclient.NewMysql(dbUrl, cf.GetBool(config.SectionServer, "app.debug"))
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		return err
	}

	if err != nil {
		_ = level.Error(logger).Log("redis", "connect", "err", err)
		return err
	}
	_ = level.Info(logger).Log("rds", "connect", "success", true)

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
	// 实例化cache
	cacheSvc = kitcache.New(logger, logging.TraceId, rds)
	cacheSvc = kitcache.NewLoggingServer(logger, cacheSvc, logging.TraceId)
	// 实例化仓库
	store = repository.NewRepository(db, logger, logging.TraceId, tracer, rds)
	// 实例化外部API
	apiSvc = api.NewApi(logger, logging.TraceId, tracer, cf, cacheSvc)

	return err
}

func Run() {
	httpAddr = envString("HTTP_ADDR", httpAddr)
	configPath = envString("CONFIG_PATH", configPath)
	env = envString("ENV", env)
	webPath = envString("WEB_PATH", DefaultWebPath)
	namespace = envString("POD_NAMESPACE", envString("NAMESPACE", namespace))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func addFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}
