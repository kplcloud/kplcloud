/**
 * @Time : 2019-06-21 10:41
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package kubernetes

import (
	"github.com/icowan/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

const (
	defaultQPS   = 1e2
	defaultBurst = 1e2
)

type K8sClient interface {
	Do() *kubernetes.Clientset
	Config() *rest.Config
}

type client struct {
	clientSet *kubernetes.Clientset
	config    *rest.Config
}

func NewClient(cf *config.Config) (cli K8sClient, err error) {
	cliConfig, err := clientcmd.BuildConfigFromFlags("", cf.GetString("server", "kube_config"))
	if err != nil {
		return
	}

	cliConfig.QPS = defaultQPS
	cliConfig.Burst = defaultBurst
	cliConfig.Timeout = time.Second * 10

	// create the clientset
	clientset, err := kubernetes.NewForConfig(cliConfig)
	if err != nil {
		return
	}

	return &client{clientSet: clientset, config: cliConfig}, nil
}

func (c *client) Do() *kubernetes.Clientset {
	return c.clientSet
}

func (c *client) Config() *rest.Config {
	return c.config
}
