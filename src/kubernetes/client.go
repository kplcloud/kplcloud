/**
 * @Time : 2019-06-21 10:41
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package kubernetes

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

type Middleware func(K8sClient) K8sClient

type K8sClient interface {
	Do(ctx context.Context) *kubernetes.Clientset
	Config(ctx context.Context) *rest.Config
	Reload(ctx context.Context) (err error)
	Connect(ctx context.Context, name, configData string) (err error)
}

type client struct {
	clientSet      map[string]*kubernetes.Clientset
	config         map[string]*rest.Config
	defaultCluster string
	store          repository.Repository
}

func (c *client) Connect(ctx context.Context, name, configData string) (err error) {
	if strings.EqualFold(configData, "") {
		cluster, err := c.store.Cluster(ctx).FindByName(ctx, name)
		if err != nil {
			return errors.Wrap(err, "store.Cluster.FindByName")
		}
		configData = cluster.ConfigData
	}

	defer func() {
		if e := os.RemoveAll("/tmp/config"); e != nil {
			log.Println("os.RemoveAll", "err", err.Error())
		}
	}()
	_ = ioutil.WriteFile("/tmp/config", []byte(configData), os.ModePerm)
	//cliConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (config *clientcmdapi.Config, e error) {
	//	e = yaml.Unmarshal([]byte(configData), &config)
	//	return
	//})
	cliConfig, err := clientcmd.BuildConfigFromFlags("", "/tmp/config")
	if err != nil {
		err = errors.Wrap(err, "clientcmd.BuildConfigFromFlags")
		return
	}
	cliConfig.QPS = defaultQPS
	cliConfig.Burst = defaultBurst
	cliConfig.Timeout = time.Second * 30

	clientSet, err := kubernetes.NewForConfig(cliConfig)
	if err != nil {
		err = errors.Wrap(err, "kubernetes.NewForConfig")
		return
	}

	c.clientSet[name] = clientSet
	c.config[name] = cliConfig

	return
}

func (c *client) Reload(ctx context.Context) (err error) {
	clusters, err := c.store.Cluster(context.Background()).FindAll(context.Background(), 1)
	if err != nil {
		err = errors.Wrap(err, "store.Cluster")
		return
	}
	for _, v := range clusters {
		if err = c.Connect(ctx, v.Name, v.ConfigData); err != nil {
			log.Println(fmt.Sprintf("cluserName: %s err: %s", v.Name, err.Error()))
		}
	}
	return nil
}

// NewClient TODO: logging
func NewClient(store repository.Repository) (cli K8sClient, err error) {
	clusters, err := store.Cluster(context.Background()).FindAll(context.Background(), 1)
	if err != nil {
		err = errors.Wrap(err, "store.Cluster")
		return
	}
	var defaultCluster string
	if clusters != nil && len(clusters) > 0 {
		defaultCluster = clusters[0].Name
	}

	clientSetMap := map[string]*kubernetes.Clientset{}
	clientConfigMap := map[string]*rest.Config{}

	for _, v := range clusters {
		_ = ioutil.WriteFile("/tmp/config", []byte(v.ConfigData), os.ModePerm)
		cliConfig, err := clientcmd.BuildConfigFromFlags("", "/tmp/config")
		if err != nil {
			log.Print("err", "clientcmd.BuildConfigFromFlags")
			continue
		}
		cliConfig.QPS = defaultQPS
		cliConfig.Burst = defaultBurst
		cliConfig.Timeout = time.Second * 30

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
	cluster, ok := ctx.Value(middleware.ContextKeyClusterName).(string)
	if !ok {
		cluster = "c1"
	}
	return c.clientSet[cluster]
}

func (c *client) Config(ctx context.Context) *rest.Config {
	cluster, ok := ctx.Value(middleware.ContextKeyClusterName).(string)
	if !ok {
		cluster = "c1"
	}
	return c.config[cluster]
}
