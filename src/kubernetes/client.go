/**
 * @Time : 2019-06-21 10:41
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package kubernetes

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
)

const (
	defaultQPS   = 1e2
	defaultBurst = 1e2
)

type K8sClient interface {
	Do(ctx context.Context) *kubernetes.Clientset
	Config(ctx context.Context) *rest.Config
	Reload(ctx context.Context) (err error)
	Connect(ctx context.Context, name string) (err error)
}

type client struct {
	clientSet      map[string]*kubernetes.Clientset
	config         map[string]*rest.Config
	defaultCluster string
	store          repository.Repository
}

func (c *client) Connect(ctx context.Context, name string) (err error) {
	cluster, err := c.store.Cluster(ctx).FindByName(ctx, name)
	if err != nil {
		err = errors.Wrap(err, "store.Cluster.FindByName")
		return
	}
	cliConfig, err := clientcmd.BuildConfigFromFlags("", cluster.ConfigData)
	if err != nil {
		err = errors.Wrap(err, "clientcmd.BuildConfigFromFlags")
		return
	}
	cliConfig.QPS = defaultQPS
	cliConfig.Burst = defaultBurst
	cliConfig.Timeout = time.Second * 10

	clientSet, err := kubernetes.NewForConfig(cliConfig)
	if err != nil {
		err = errors.Wrap(err, "kubernetes.NewForConfig")
		return
	}

	c.clientSet[cluster.Name] = clientSet
	c.config[cluster.Name] = cliConfig

	return
}

func (c *client) Reload(ctx context.Context) (err error) {
	clusters, err := c.store.Cluster(context.Background()).FindAll(context.Background(), 1)
	if err != nil {
		err = errors.Wrap(err, "store.Cluster")
		return
	}
	var defaultCluster string
	if clusters != nil {
		defaultCluster = clusters[0].Name
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
			log.Println("err", "kubernetes.NewForConfig")
			continue
		}
		clientSetMap[v.Name] = clientSet
		clientConfigMap[v.Name] = cliConfig
	}
	c.clientSet = clientSetMap
	c.config = clientConfigMap
	c.defaultCluster = defaultCluster
	return nil
}

// TODO: logging
func NewClient(store repository.Repository) (cli K8sClient, err error) {

	clusters, err := store.Cluster(context.Background()).FindAll(context.Background(), 1)
	if err != nil {
		err = errors.Wrap(err, "store.Cluster")
		return
	}
	var defaultCluster string
	if clusters != nil {
		defaultCluster = clusters[0].Name
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
			log.Println("err", "kubernetes.NewForConfig")
			continue
		}
		clientSetMap[v.Name] = clientSet
		clientConfigMap[v.Name] = cliConfig
	}

	return &client{store: store, clientSet: clientSetMap, config: clientConfigMap, defaultCluster: defaultCluster}, nil
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
