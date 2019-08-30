package istio

//type IstioClient interface {
//	Do() *crd.Client
//}
//
//type istioClient struct {
//	client *crd.Client
//}
//
//func NewClient(cf *config.Config) (IstioClient, error) {
//	client, err := crd.NewClient(cf.GetString("server", "kube_config"), "", model.IstioConfigTypes, "")
//	if err != nil {
//		return nil, err
//	}
//	return &istioClient{client: client}, nil
//}
//
//func (c *istioClient) Do() *crd.Client {
//	return c.client
//}
