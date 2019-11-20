/**
 * @Time : 2019-06-26 14:34
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"errors"
	"github.com/ghodss/yaml"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var (
	ErrPvcK8sList            = errors.New("存储卷声明列表获取失败")
	ErrPvcGet                = errors.New("存储卷声明获取错误")
	ErrPvcPost               = errors.New("存储卷声明创建失败")
	ErrPvcDelete             = errors.New("存储卷声明删除错误,或许已经删除了")
	ErrPvcTemplateGet        = errors.New("存储卷声明模版获取错误")
	ErrPvcTemplateEncode     = errors.New("存储卷声明模版解析错误")
	ErrPvGet                 = errors.New("存储卷获取错误")
	ErrStorageClassNotExists = errors.New("存储类不存在")
	ErrPvcListCount          = errors.New("存储卷声明统计出错")
	ErrPvcList               = errors.New("存储卷声明获取出错")
)

type Service interface {
	// 同步pvc
	Sync(ctx context.Context, ns string) (err error)

	// 获取pvc详情
	Get(ctx context.Context, ns, name string) (rs interface{}, err error)

	// 删除存储卷声明
	Delete(ctx context.Context, ns, name string) (err error)

	// 创建持久化存储卷
	Post(ctx context.Context, ns, name, storage, storageClassName string, accessModes []string) (err error)

	// 持久化存储卷列表
	List(ctx context.Context, ns string, page, limit int) (resp map[string]interface{}, err error)

	// 当前空间下所有的pvc
	All(ctx context.Context) (resp map[string]interface{}, err error)
}

type service struct {
	logger     log.Logger
	k8sClient  kubernetes.K8sClient
	repository repository.Repository
}

func NewService(logger log.Logger, client kubernetes.K8sClient,
	repository repository.Repository) Service {
	return &service{
		logger: logger, k8sClient: client,
		repository: repository,
	}
}

func (c *service) All(ctx context.Context) (resp map[string]interface{}, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	list, err := c.repository.Pvc().FindBy(ns, 0, 100)
	if err != nil {
		_ = level.Error(c.logger).Log("pvcRepository", "FindBy", "err", err.Error())
		return nil, ErrPvcList
	}
	return map[string]interface{}{
		"items": list,
	}, nil
}

func (c *service) Post(ctx context.Context, ns, name, storage, storageClassName string, accessModes []string) (err error) {
	if class, err := c.repository.StorageClass().Find(storageClassName); err != nil || class.Name == "" {
		_ = level.Error(c.logger).Log("classRepository", "Find", "err", err)
		return ErrStorageClassNotExists
	}

	var pvc *v1.PersistentVolumeClaim

	defer func() {
		b, _ := yaml.Marshal(pvc)
		if err == nil && c.repository.Pvc().FirstOrCreate(ns, name, accessModes[0], storage, storageClassName, string(b), pvc.Spec.Selector.String(), map[string]string{}) != nil {
			_ = level.Warn(c.logger).Log("pvcRepository", "FirstOrCreate", "err", "FirstOrCreate err")
		}
	}()

	tpl, err := c.repository.Template().FindByKindType(repository.PersistentVolumeClaimKind)
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByKindType", "err", err.Error())
		return ErrPvcTemplateGet
	}

	enTpl, err := encode.EncodeTemplate(repository.StorageClassKind.ToString(), tpl.Detail, map[string]interface{}{
		"name":             name,
		"namespace":        ns,
		"accessModes":      accessModes,
		"storage":          storage,
		"storageClassName": storageClassName,
	})

	if err != nil {
		_ = level.Error(c.logger).Log("encode", "EncodeTemplate", "err", err.Error())
		return ErrPvcTemplateEncode
	}
	err = yaml.Unmarshal([]byte(enTpl), &pvc)
	if err != nil {
		_ = level.Error(c.logger).Log("yaml", "Unmarshal", "err", err.Error())
		return
	}

	if pvc, err = c.k8sClient.Do().CoreV1().PersistentVolumeClaims(ns).Create(pvc); err != nil {
		_ = level.Error(c.logger).Log("PersistentVolumeClaims", "Create", "err", err.Error())
		return ErrPvcPost
	}

	return
}

func (c *service) List(ctx context.Context, ns string, page, limit int) (rs map[string]interface{}, err error) {
	total, err := c.repository.Pvc().Count(ns)
	if err != nil {
		_ = level.Error(c.logger).Log("pvcRepository", "Count", "err", err.Error())
		return rs, ErrPvcListCount
	}

	p := paginator.NewPaginator(page, limit, int(total))

	res, err := c.repository.Pvc().FindBy(ns, p.Offset(), p.PerPageNums())
	if err != nil {
		_ = level.Error(c.logger).Log("pvcRepository", "FindBy", "err", err.Error())
		return rs, ErrPvcList
	}

	var resp v1.PersistentVolumeClaimList

	for _, v := range res {
		var pv v1.PersistentVolumeClaim
		_ = yaml.Unmarshal([]byte(v.Detail.String), &pv)
		resp.Items = append(resp.Items, pv)
	}

	return map[string]interface{}{
		"items": resp.Items,
		"page":  p.Result(),
	}, nil
}

func (c *service) Sync(ctx context.Context, ns string) (err error) {
	pvcs, err := c.k8sClient.Do().CoreV1().PersistentVolumeClaims(ns).List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("PersistentVolumeClaims", "List", "err", err.Error())
		return ErrPvcK8sList
	}

	for _, pvc := range pvcs.Items {
		b, _ := yaml.Marshal(pvc)
		storage, _ := pvc.Spec.Resources.Requests[v1.ResourceStorage].MarshalJSON()
		if err = c.repository.Pvc().FirstOrCreate(ns, pvc.Name,
			string(pvc.Spec.AccessModes[0]), strings.Trim(string(storage), `"`),
			*pvc.Spec.StorageClassName,
			string(b),
			pvc.Spec.Selector.String(), pvc.Labels); err != nil {
			_ = level.Warn(c.logger).Log("pvcRepository", "FirstOrCreate", "err", err)
		}
	}

	return
}

func (c *service) Get(ctx context.Context, ns, name string) (rs interface{}, err error) {
	_, err = c.repository.Pvc().Find(ns, name)
	if err != nil {
		_ = level.Error(c.logger).Log("pvcRepository", "Find", "err", err.Error())
		return nil, ErrPvcGet
	}

	p, err := c.k8sClient.Do().CoreV1().PersistentVolumeClaims(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("PersistentVolumeClaims", "Get", "err", err.Error())
		return nil, ErrPvcGet
	}

	pv, err := c.k8sClient.Do().CoreV1().PersistentVolumes().Get(p.Spec.VolumeName, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("PersistentVolumes", "Get", "err", err.Error())
		return nil, ErrPvGet
	}

	return map[string]interface{}{
		"pvc": p,
		"pv":  pv,
	}, nil
}

func (c *service) Delete(ctx context.Context, ns, name string) (err error) {
	defer func() {
		if err == nil {
			if e := c.repository.Pvc().Delete(ns, name); e != nil {
				_ = level.Warn(c.logger).Log("pvcRepository", "Delete", "err", e.Error())
			}
		}
	}()

	err = c.k8sClient.Do().CoreV1().PersistentVolumeClaims(ns).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("PersistentVolumeClaims", "Delete", "err", err.Error())
		return ErrPvcDelete
	}

	return
}
