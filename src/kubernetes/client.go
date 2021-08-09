/**
 * @Time : 2019-06-21 10:41
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package kubernetes

import (
	"context"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

const (
	defaultQPS   = 1e2
	defaultBurst = 1e2
)

type K8sClient interface {
	Do(ctx context.Context) *kubernetes.Clientset
	Config(ctx context.Context) *rest.Config
}

type client struct {
	clientSet map[string]*kubernetes.Clientset
	config    map[string]*rest.Config
}

// TODO: logging
func NewClient(store repository.Repository) (cli K8sClient, err error) {

	clusters, err := store.Cluster(context.Background()).FindAll(context.Background(), 1)
	if err != nil {
		err = errors.Wrap(err, "store.Cluster")
		return
	}

	clientSetMap := map[string]*kubernetes.Clientset{}
	clientConfigMap := map[string]*rest.Config{}

	for _, v := range clusters {
		cliConfig, err := clientcmd.BuildConfigFromFlags("", v.ConfigData)
		if err != nil {
			log.Print("err", "clientcmd.BuildConfigFromFlags")
			continue
		}
		cliConfig.QPS = defaultQPS
		cliConfig.Burst = defaultBurst
		cliConfig.Timeout = time.Second * 10

		clientSet, err := kubernetes.NewForConfig(cliConfig)
		if err != nil {
			log.Print("err", "kubernetes.NewForConfig")
			continue
		}
		clientSetMap[v.Name] = clientSet
		clientConfigMap[v.Name] = cliConfig
	}

	return &client{clientSet: clientSetMap, config: cliConfig}, nil
}

func (c *client) Do(ctx context.Context) *kubernetes.Clientset {
	cluster, ok := ctx.Value(middleware.ClusterContextKey).(string)
	if !ok {
		cluster = "c1"
	}
	return c.clientSet[cluster]
}

func (c *client) Config(ctx context.Context) *rest.Config {
	cluster, ok := ctx.Value(middleware.ClusterContextKey).(string)
	if !ok {
		cluster = "c1"
	}
	return c.config[cluster]
}
