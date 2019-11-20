package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/cmd"
	"github.com/kplcloud/kplcloud/src/config"
	jenkins2 "github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/mysql"
	"github.com/kplcloud/kplcloud/src/pkg/configmap"
	"github.com/kplcloud/kplcloud/src/pkg/consul"
	"github.com/kplcloud/kplcloud/src/pkg/namespace"
	"github.com/kplcloud/kplcloud/src/pkg/persistentvolumeclaim"
	"github.com/kplcloud/kplcloud/src/pkg/storage"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/spf13/cobra"
	"os"
)

var logger log.Logger

var (
	configPath = envString("CONFIG_PATH", "./app.cfg")
	kubeConfig = envString("KUBE_CONFIG", "~/.kube/config")
	ns         string
	store      repository.Repository
	k8sClient  kubernetes.K8sClient
	db         *gorm.DB
	cf         *config.Config
	jenkins    jenkins2.Jenkins

	rootCmd = &cobra.Command{
		Use:               "client",
		Short:             "开普勒平台服务端",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# 开普勒平台服务端

您可以通过改命令来启动您的服务

可用的配置类型：

[client]

有关开普勒平台的相关概述，请参阅 https://github.com/kplcloud/kplcloud
`,
	}

	syncCmd = &cobra.Command{
		Use:   "sync [OPTIONS]",
		Short: "同步kubernetes数据",
		Example: `## 同步kubernetes数据到数据库
client sync namespace -c app.cfg -k k8s-config.yaml
client sync pvc -c app.cfg -k k8s-config.yaml -n app
client sync storage -c app.cfg -k k8s-config.yaml -n app
client sync consul -c app.cfg -k k8s-config.yaml -n app
client sync configmap -c app.cfg -k k8s-config.yaml -n app
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("flag is nil")
			}
			var err error
			cf = getConfig()
			if k8sClient, err = kubernetes.NewClient(cf); err != nil {
				return err
			}
			if db, err = mysql.NewDb(cf); err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			store = repository.NewRepository(db)

			if len(args) == 1 && args[0] == "namespace" {
				return syncNamespace()
			}

			switch args[0] {
			case "storage":
				return syncStorage()
			case "pvc":
				return syncPvc()
			case "consul":
				return syncConsul()
			case "config":
				return syncConfigMap()
			}

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", "./app.cfg", "配置文件路径: ./app.cfg")
	rootCmd.PersistentFlags().StringVarP(&kubeConfig, "kube.config", "k", "~/.kube/config", "kubeconfig文件路径: ~/.kube/config")

	syncCmd.PersistentFlags().StringVarP(&ns, "namespace", "n", "default", "操作的空间")

	cmd.AddFlags(rootCmd)
	rootCmd.AddCommand(syncCmd)
}

func main() {
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

func getConfig() *config.Config {
	cf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	_ = cf.SetValue("server", "kube_config", kubeConfig)

	return cf
}

func syncNamespace() error {
	if err := namespace.NewService(logger, cf, k8sClient, store).Sync(context.Background()); err != nil {
		_ = logger.Log("svc", "sync", "err", err.Error())
		return err
	}
	return nil
}

func syncStorage() error {
	return storage.NewService(logger, k8sClient, store).Sync(context.Background())
}

func syncPvc() error {
	return persistentvolumeclaim.NewService(logger, k8sClient, store).Sync(context.Background(), ns)
}

func syncConsul() error {
	return consul.NewService(logger, cf, store).Sync(context.Background())
}

func syncConfigMap() error {
	return configmap.NewService(logger, cf, jenkins, k8sClient, store).Sync(context.Background(), ns)
}
