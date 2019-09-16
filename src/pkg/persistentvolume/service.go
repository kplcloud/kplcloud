/**
 * @Time : 2019-06-25 19:24
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package persistentvolume

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrPersistentVolumeK8sGet = errors.New("获取存储卷错误")
)

type Service interface {
	// 获取pv详情
	Get(ctx context.Context, name string) (rs interface{}, err error)

	// 获取pv列表
	List(ctx context.Context, ns string)

	// 删除pv
	Delete(ctx context.Context, ns, name string) (err error)

	// 同步pv
	Sync(ctx context.Context, ns string) (err error)
}

type service struct {
	logger    log.Logger
	k8sClient kubernetes.K8sClient
}

func NewService(logger log.Logger, client kubernetes.K8sClient) Service {
	return &service{logger, client}
}

func (c *service) Get(ctx context.Context, name string) (rs interface{}, err error) {
	pv, err := c.k8sClient.Do().CoreV1().PersistentVolumes().Get(name, metav1.GetOptions{})
	c.k8sClient.Do().CoreV1()
	if err != nil {
		_ = level.Error(c.logger).Log("PersistentVolumes", "Get", "err", err.Error())
		return nil, ErrPersistentVolumeK8sGet
	}

	return pv, nil
}

func (c *service) List(ctx context.Context, ns string) {

}

func (c *service) Delete(ctx context.Context, ns, name string) (err error) {

	return
}

func (c *service) Sync(ctx context.Context, ns string) (err error) {
	return
}
