/**
 * @Time : 8/11/21 11:43 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
)

type Middleware func(Service) Service

type Service interface {
	Sync(ctx context.Context, clusterId int64, ns string) (err error)
	// PutImage 手动更新Image
	PutImage(ctx context.Context, clusterId int64, ns, name, image string) (err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func (s *service) PutImage(ctx context.Context, clusterId int64, ns, name, image string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	get, err := s.k8sClient.Do(ctx).AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.AppsV1.Deployments", "Get", "err", err.Error())
		return encode.ErrDeploymentGetNotfound.Wrap(err)
	}
	for k, v := range get.Spec.Template.Spec.Containers {
		if strings.EqualFold(v.Name, name) {
			get.Spec.Template.Spec.Containers[k].Image = image
			break
		}
	}

	update, err := s.k8sClient.Do(ctx).AppsV1().Deployments(ns).Update(ctx, get, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	b, _ := json.Marshal(update)
	fmt.Println(string(b))

	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().AsInt64())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().AsInt64())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().AsInt64())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().AsInt64())
	fmt.Println(update.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())

	return
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	var nss *appv1.DeploymentList
	if nss, err = s.k8sClient.Do(ctx).AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{}); err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.AppsV1.Deployments", "List", "err", err.Error())
		return encode.ErrDeploymentSyncList.Wrap(err)
	}

	fmt.Println(ns)

	for _, v := range nss.Items {
		b, _ := json.Marshal(v)
		fmt.Println(string(b))
		fmt.Println("request.Cpu", v.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
		fmt.Println("limit.Cpu", v.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().Value())
		fmt.Println("request.Memory", v.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value())
		fmt.Println("limit.Memory", v.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().Value())
		fmt.Println("---")
		fmt.Println(v.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().AsInt64())
		fmt.Println(v.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().AsInt64())
		fmt.Println(v.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().Value())
		fmt.Println(v.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().Value())
		fmt.Println(v.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String())
		fmt.Println(v.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
	}

	return
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, repository repository.Repository) Service {
	logger = log.With(logger, "deployment", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		k8sClient:  client,
		repository: repository,
	}
}
