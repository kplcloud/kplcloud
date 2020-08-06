package virtualservice

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/kplcloud/kplcloud/src/istio"
	"github.com/kplcloud/kplcloud/src/istio/types/v1beta1"
)

type Service interface {
	List(ctx context.Context)

	Get(ctx context.Context, ns, name string)
}

type service struct {
	logger      log.Logger
	istioClient istio.IstioClient
}

func (s *service) List(ctx context.Context) {

}

func (s *service) Get(ctx context.Context, ns, name string) {
	logger := log.With(s.logger, "namespace", ns, "name", name)

	var vs v1beta1.VirtualService

	opts := metav1.GetOptions{
		ResourceVersion: v1beta1.GroupVersion,
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: v1beta1.GroupVersion,
		},
	}

	err := s.istioClient.Do().Get().Namespace(ns).
		Resource("virtualservices").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().Into(&vs)

	fmt.Println(err)
	if err != nil {
		_ = level.Error(logger).Log("istioClient.Do", "Get", ns, "name", name, "resource", "virtualservices", "err", err.Error())
		return
	}

	fmt.Println(vs)
	fmt.Println("-----")
}

func NewService(logger log.Logger, istioClient istio.IstioClient) Service {
	logger = log.With(logger, "service", "virtualservice")
	return &service{
		logger:      logger,
		istioClient: istioClient,
	}
}
