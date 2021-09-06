/**
 * @Time : 2019-06-26 14:34
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//var (
//	ErrPvcK8sList            = errors.New("存储卷声明列表获取失败")
//	ErrPvcGet                = errors.New("存储卷声明获取错误")
//	ErrPvcPost               = errors.New("存储卷声明创建失败")
//	ErrPvcDelete             = errors.New("存储卷声明删除错误,或许已经删除了")
//	ErrPvcTemplateGet        = errors.New("存储卷声明模版获取错误")
//	ErrPvcTemplateEncode     = errors.New("存储卷声明模版解析错误")
//	ErrPvGet                 = errors.New("存储卷获取错误")
//	ErrStorageClassNotExists = errors.New("存储类不存在")
//	ErrPvcListCount          = errors.New("存储卷声明统计出错")
//	ErrPvcList               = errors.New("存储卷声明获取出错")
//)

type Middleware func(Service) Service

type Service interface {
	// Sync 同步pvc
	Sync(ctx context.Context, clusterId int64, ns string) (err error)
	// Get 获取pvc详情
	Get(ctx context.Context, clusterId int64, ns, name string) (rs interface{}, err error)
	// Delete 删除存储卷声明
	Delete(ctx context.Context, clusterId int64, ns, name string) (err error)
	// Create 创建持久化存储卷
	Create(ctx context.Context, clusterId int64, ns, name, storage, storageClassName string, accessModes []string) (err error)
	// List 持久化存储卷列表
	List(ctx context.Context, clusterId int64, ns string, page, pageSize int) (resp map[string]interface{}, err error)
	// All 当前空间下所有的pvc
	All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error)
}

type service struct {
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
	traceId    string
}

func (s *service) Sync(ctx context.Context, clusterId int64, ns string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	pvcs, err := s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		_ = level.Error(logger).Log("PersistentVolumeClaims", "List", "err", err.Error())
		return encode.ErrPersistentVolumeClaimList.Wrap(err)
	}

	for _, pvc := range pvcs.Items {
		b, _ := yaml.Marshal(pvc)
		//storage, _ := pvc.Spec.Resources.Requests[v1.ResourceStorage].MarshalJSON()
		//if err = c.repository.Pvc().FirstOrCreate(ns, pvc.Name,
		//	string(pvc.Spec.AccessModes[0]), strings.Trim(string(storage), `"`),
		//	*pvc.Spec.StorageClassName,
		//	string(b),
		//	pvc.Spec.Selector.String(), pvc.Labels); err != nil {
		//	_ = level.Warn(c.logger).Log("pvcRepository", "FirstOrCreate", "err", err)
		//}
		fmt.Println(string(b))
	}

	return
}

func (s *service) Get(ctx context.Context, clusterId int64, ns, name string) (rs interface{}, err error) {
	//_, err = c.repository.Pvc().Find(ns, name)
	//if err != nil {
	//	_ = level.Error(c.logger).Log("pvcRepository", "Find", "err", err.Error())
	//	return nil, ErrPvcGet
	//}
	//
	//p, err := c.k8sClient.Do().CoreV1().PersistentVolumeClaims(ns).Get(name, metav1.GetOptions{})
	//if err != nil {
	//	_ = level.Error(c.logger).Log("PersistentVolumeClaims", "Get", "err", err.Error())
	//	return nil, ErrPvcGet
	//}
	//
	//pv, err := c.k8sClient.Do().CoreV1().PersistentVolumes().Get(p.Spec.VolumeName, metav1.GetOptions{})
	//if err != nil {
	//	_ = level.Error(c.logger).Log("PersistentVolumes", "Get", "err", err.Error())
	//	return nil, ErrPvGet
	//}
	//
	//return map[string]interface{}{
	//	"pvc": p,
	//	"pv":  pv,
	//}, nil
	return
}

func (s *service) Delete(ctx context.Context, clusterId int64, ns, name string) (err error) {
	panic("implement me")
}

func (s *service) Create(ctx context.Context, clusterId int64, ns, name, storage, storageClassName string, accessModes []string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	sc, err := s.repository.StorageClass(ctx).FindName(ctx, clusterId, storageClassName)
	if err != nil {
		_ = level.Error(logger).Log("repository.StorageClass", "FindName", "err", err)
		return encode.ErrStorageClassNotfound.Error()
	}
	// todo 查询是否存在pvc
	var pvc *corev1.PersistentVolumeClaim
	tpl, err := s.repository.K8sTpl(ctx).EncodeTemplate(ctx, types.KindPersistentVolumeClaim, map[string]interface{}{
		"name":             name,
		"namespace":        ns,
		"accessModes":      accessModes,
		"storage":          storage,
		"storageClassName": sc.Name,
	}, &pvc)
	if err != nil {
		return encode.ErrPersistentVolumeClaimCreate.Wrap(err)
	}
	fmt.Println(string(tpl))
	fmt.Println(pvc)

	pvc, err = s.k8sClient.Do(ctx).CoreV1().PersistentVolumeClaims(ns).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		_ = level.Error(logger).Log("CoreV1.PersistentVolumeClaims", "Create", "err", err.Error())
		return encode.ErrPersistentVolumeClaimCreate.Wrap(err)
	}

	// todo 保存到数据库

	return
}

func (s *service) List(ctx context.Context, clusterId int64, ns string, page, pageSize int) (resp map[string]interface{}, err error) {
	panic("implement me")
}

func (s *service) All(ctx context.Context, clusterId int64) (resp map[string]interface{}, err error) {
	panic("implement me")
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, repository repository.Repository) Service {
	return &service{
		logger: logger, k8sClient: client,
		repository: repository,
		traceId:    traceId,
	}
}
