package istio

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/kplcloud/kplcloud/src/istio/types/v1beta1"
)

type IstioClient interface {
	Do() *rest.RESTClient
}

type istioClient struct {
	config *rest.Config
	client *rest.RESTClient
}

func NewClient(client *rest.Config) (IstioClient, error) {
	crdConfig := *client
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1beta1.GroupName, Version: v1beta1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}

	return &istioClient{config: &crdConfig, client: restClient}, nil
}

func (c *istioClient) Do() *rest.RESTClient {
	return c.client
}
