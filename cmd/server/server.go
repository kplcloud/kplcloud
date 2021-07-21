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
	"github.com/icowan/config"
	kitcache "github.com/icowan/kit-cache"
	redisclient "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"

	"github.com/kplcloud/kplcloud/src/api"
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
		Use:               "kplcloud",
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

	installCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的http端口")

	addFlags(rootCmd)
	rootCmd.AddCommand(startCmd, generateCmd, settingCmd, installCmd, resetCmd)
}

func Run() {
	httpAddr = envString("HTTP_ADDR", httpAddr)
	configPath = envString("CONFIG_PATH", configPath)
	env = envString("ENV", env)
	webPath = envString("WEB_PATH", DefaultWebPath)
	namespace = envString("POD_NAMESPACE", envString("NAMESPACE", namespace))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute", err.Error())
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
