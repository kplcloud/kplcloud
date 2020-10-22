package consul

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hashicorp/consul/api"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/consul"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/paginator"
	"strings"
)

var (
	ErrConsulGet      = errors.New("获取consul信息失败")
	ErrConsulConnect  = errors.New("consul链接失败")
	ErrConsulDelete   = errors.New("consul删除失败")
	ErrConsulExist    = errors.New("已存在同名consul")
	ErrConsulCreate   = errors.New("consul创建失败")
	ErrConsulUpdate   = errors.New("consul更新失败")
	ErrConsulNotExist = errors.New("consul不存在")
	ErrConsulKVGet    = errors.New("KV信息获取失败")
	ErrConsulKVParams = errors.New("参数错误")
	ErrConsulKVPost   = errors.New("KV创建失败")
	ErrConsulKVDelete = errors.New("KV删除失败")
)

type Service interface {
	// ACL信息同步到DB
	Sync(ctx context.Context) (err error)

	// 获取ACL详情
	Detail(ctx context.Context, ns, name string) (res map[string]interface{}, err error)

	// 获取ACL列表
	List(ctx context.Context, ns, name string, page, limit int) (res map[string]interface{}, err error)

	// 创建ACL
	Post(ctx context.Context, ns, name, clientType, rules string) error

	// 更新ACL
	Update(ctx context.Context, ns, name, clientType, rules string) error

	// 删除ACL
	Delete(ctx context.Context, ns, name string) error

	// 获取KV详情
	KVDetail(ctx context.Context, ns, name, prefix string) (res map[string]interface{}, err error)

	// 获取KV目录列表
	KVList(ctx context.Context, ns, name, prefix string) (res map[string]interface{}, err error)

	// 获取KV
	KVPost(ctx context.Context, ns, name, key, value string) (err error)

	// 删除KV或目录
	KVDelete(ctx context.Context, ns, name, prefix string, filderState bool) (err error)
}

type service struct {
	logger     log.Logger
	cf         *config.Config
	repository repository.Repository
}

func NewService(logger log.Logger, cf *config.Config, store repository.Repository) Service {
	return &service{
		logger,
		cf,
		store,
	}
}

func (c *service) Sync(ctx context.Context) (err error) {
	consulClient, err := consul.NewConsulClient(c.cf)
	if err != nil {
		_ = level.Error(c.logger).Log("Sync", "NewConsulClient", "err", err.Error())
		return
	}
	aclList, err := consulClient.ACLList(nil)
	if err != nil {
		_ = level.Error(c.logger).Log("Sync", "ACLList", "err", err.Error())
		return
	}

	go func() {
		for _, v := range aclList {
			nameStr := strings.Split(v.Name, ".")
			if len(nameStr) < 2 {
				_ = level.Error(c.logger).Log("Consul", "Sync", "Refused", v.Name)
				continue
			}
			if err = c.repository.Consul().UpdateOrCreate(v, nameStr[0], nameStr[1]); err != nil {
				_ = level.Error(c.logger).Log("Sync", "UpdateOrCreate", err.Error())
			}

		}
	}()

	return nil
}

func (c *service) Detail(ctx context.Context, ns, name string) (res map[string]interface{}, err error) {
	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		err = ErrConsulGet
		return
	}

	if consulInfo.Token != "" {
		consulInfo.EncryptToken = strings.Replace(consulInfo.Token, consulInfo.Token[10:24], "**************", -1)
	}
	res = map[string]interface{}{
		"id":         consulInfo.ID,
		"name":       consulInfo.Name,
		"namespace":  consulInfo.Namespace,
		"fullToken":  consulInfo.Token,
		"type":       consulInfo.Type,
		"rules":      consulInfo.Rules,
		"token":      consulInfo.EncryptToken,
		"consulRule": nil,
	}

	consulRule := map[string]interface{}{}
	err = json.Unmarshal([]byte(consulInfo.Rules), &consulRule)
	if err == nil {
		res["consulRule"] = consulRule
	}

	return
}

func (c *service) List(ctx context.Context, ns, name string, page, limit int) (res map[string]interface{}, err error) {
	count, err := c.repository.Consul().Count(ns, name)
	if err != nil {
		_ = level.Error(c.logger).Log("Consul", "List Count", "err", err.Error())
		return nil, ErrConsulGet
	}

	p := paginator.NewPaginator(page, limit, count)

	list, err := c.repository.Consul().FindOffsetLimit(ns, name, p.Offset(), limit)

	if err != nil {
		_ = level.Error(c.logger).Log("Consul", "List", "err", err.Error())
		return nil, ErrConsulGet
	}

	var dat []map[string]interface{}
	for _, val := range list {
		var token string
		if val.Token != "" {
			token = strings.Replace(val.Token, val.Token[10:24], "**************", -1)
		}
		dat = append(dat, map[string]interface{}{
			"id":        val.ID,
			"name":      val.Name,
			"namespace": val.Namespace,
			"token":     token,
			"type":      val.Type,
		})
	}
	res = map[string]interface{}{
		"items": dat,
		"page":  p.Result(),
	}
	return
}

func (c *service) Post(ctx context.Context, ns, name, clientType, rules string) error {
	if _, notExist := c.repository.Consul().Find(ns, name); notExist == false {
		return ErrConsulExist
	}

	consulClient, err := consul.NewConsulClient(c.cf)
	if err != nil {
		_ = level.Error(c.logger).Log("Post", "NewConsulClient", "err", err.Error())
		return ErrConsulConnect
	}

	acl, err := consulClient.ACLCreate(&api.ACLEntry{
		Name:  ns + "." + name,
		Type:  clientType,
		Rules: rules,
	})
	if err != nil {
		_ = level.Error(c.logger).Log("Post", "ACLCreate", "err", err.Error())
		return ErrConsulCreate
	}

	if _, err := c.repository.Consul().FirstOrCreate(acl, ns, name); err != nil {
		_ = level.Error(c.logger).Log("Post", "FirstOrCreate", "err", err.Error())
		return ErrConsulCreate
	}

	return nil
}

func (c *service) Update(ctx context.Context, ns, name, clientType, rules string) error {
	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		return ErrConsulNotExist
	}

	consulClient, err := consul.NewConsulClient(c.cf)
	if err != nil {
		_ = level.Error(c.logger).Log("Update", "NewConsulClient", "err", err.Error())
		return ErrConsulConnect
	}
	ae := api.ACLEntry{
		Name:  ns + "." + name,
		ID:    consulInfo.Token,
		Type:  clientType,
		Rules: rules,
	}
	if err := consulClient.ACLUpdate(&ae); err != nil {
		_ = level.Error(c.logger).Log("Update", "ACLUpdate", "err", err.Error())
		return ErrConsulUpdate
	}

	if err = c.repository.Consul().UpdateOrCreate(&ae, ns, name); err != nil {
		_ = level.Error(c.logger).Log("Update", "UpdateOrCreate", "err", err.Error())
		return ErrConsulUpdate
	}

	return nil
}

func (c *service) Delete(ctx context.Context, ns, name string) error {
	consulClient, err := consul.NewConsulClient(c.cf)
	if err != nil {
		_ = level.Error(c.logger).Log("Delete", "NewConsulClient", "err", err.Error())
		return ErrConsulConnect
	}

	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		_ = level.Error(c.logger).Log("Delete", "Find Not Found")
		return ErrConsulGet
	}

	if err = consulClient.ACLDelete(consulInfo.Token); err != nil {
		_ = level.Error(c.logger).Log("Delete", "ACLDelete", "err", err.Error())
		return ErrConsulDelete
	}

	if err = c.repository.Consul().Delete(ns, name); err != nil {
		_ = level.Error(c.logger).Log("Delete", "DB", "err", err.Error())
		return ErrConsulDelete
	}

	return nil
}

func (c *service) KVDetail(ctx context.Context, ns, name, prefix string) (res map[string]interface{}, err error) {
	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		err = ErrConsulNotExist
		return
	}

	kvClient, err := consul.NewKVClient(c.cf, consulInfo.Token)
	if err != nil {
		_ = level.Error(c.logger).Log("KVDetail", "NewKVClient", "err", err.Error())
		err = ErrConsulConnect
		return
	}

	pairs, err := kvClient.KVGet(prefix)
	if err != nil {
		_ = level.Error(c.logger).Log("KVDetail", "KVGet", "err", err.Error())
		err = ErrConsulKVGet
		return
	}
	return map[string]interface{}{
		"key":   pairs.Key,
		"value": string(pairs.Value),
	}, nil
}

func (c *service) KVList(ctx context.Context, ns, name, prefix string) (res map[string]interface{}, err error) {
	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		err = ErrConsulNotExist
		return
	}

	kvClient, err := consul.NewKVClient(c.cf, consulInfo.Token)
	if err != nil {
		_ = level.Error(c.logger).Log("KVList", "NewKVClient", "err", err.Error())
		err = ErrConsulConnect
		return
	}

	list, err := kvClient.KVList(prefix)
	if err != nil {
		_ = level.Error(c.logger).Log("KVList", "KVList", "err", err.Error())
		err = ErrConsulKVGet
		return
	}

	var data []string
	for _, value := range list {
		v := strings.Replace(string(value.Key), prefix, "", 1)
		if v == "" {
			continue
		}
		if strings.Contains(v, "/") == false {
			data = append(data, v)
			continue
		}
		index := strings.Index(v, "/")
		v = v[0 : index+1]
		if len(data) <= 0 {
			data = append(data, v)
			continue
		}

		var inData bool
		for _, vv := range data {
			if vv == v {
				inData = true
			}
		}
		if inData != true {
			data = append(data, v)
		}
	}
	return map[string]interface{}{"prefix": prefix, "detail": data}, nil
}

func (c *service) KVPost(ctx context.Context, ns, name, key, value string) (err error) {
	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		err = ErrConsulNotExist
		return
	}

	kvClient, err := consul.NewKVClient(c.cf, consulInfo.Token)
	if err != nil {
		_ = level.Error(c.logger).Log("KVPost", "NewKVClient", "err", err.Error())
		err = ErrConsulConnect
		return
	}

	if !strings.Contains(key, ns+"."+name) {
		err = ErrConsulKVParams
		return
	}

	err = kvClient.KVPut(key, value)
	if err != nil {
		_ = level.Error(c.logger).Log("KVPost", "KVPut", "err", err.Error())
		err = ErrConsulKVPost
		return
	}
	return
}

func (c *service) KVDelete(ctx context.Context, ns, name, prefix string, filderState bool) (err error) {
	consulInfo, notExist := c.repository.Consul().Find(ns, name)
	if notExist == true {
		err = ErrConsulNotExist
		return
	}

	kvClient, err := consul.NewKVClient(c.cf, consulInfo.Token)
	if err != nil {
		_ = level.Error(c.logger).Log("KVDelete", "NewKVClient", "err", err.Error())
		err = ErrConsulConnect
		return
	}
	if filderState == true {
		err = kvClient.KVDeleteTree(prefix)
	} else {
		err = kvClient.KVDelete(prefix)
	}
	if err != nil {
		_ = level.Error(c.logger).Log("KVDelete", "KVDelete", err, err.Error())
		err = ErrConsulKVDelete
		return
	}
	return nil
}
