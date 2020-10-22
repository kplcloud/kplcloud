/**
 * @Time : 2019/7/17 2:57 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : client
 * @Software: GoLand
 */

package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/icowan/config"
)

type Consul interface {
	TokenList(query *api.QueryOptions) (tokens []*api.ACLTokenListEntry, err error)
	ACLList(query *api.QueryOptions) (aclTokenList []*api.ACLEntry, err error)
	ACLCreate(query *api.ACLEntry) (acl *api.ACLEntry, err error)
	ACLUpdate(query *api.ACLEntry) (err error)
	ACLDelete(token string) (err error)
}
type consul struct {
	client *api.Client
}

func NewConsulClient(cf *config.Config) (consulClient Consul, err error) {
	conf := api.DefaultConfig()
	conf.Address = cf.GetString("consul", "consul_addr")
	conf.Token = cf.GetString("consul", "consul_token")

	cli, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &consul{client: cli}, nil
}

func (c *consul) TokenList(query *api.QueryOptions) (tokens []*api.ACLTokenListEntry, err error) {
	tokens, _, err = c.client.ACL().TokenList(&api.QueryOptions{Token: "398073a8-5091-4d9c-871a-bbbeb030d1f6"})
	return
}
func (c *consul) ACLList(query *api.QueryOptions) (acls []*api.ACLEntry, err error) {
	acls, qm, err := c.client.ACL().List(query)
	if err != nil {
		return
	}

	if len(acls) < 1 {
		err = errors.New(fmt.Sprintf("err: %v", err))
		return
	}

	if qm.LastIndex == 0 {
		err = errors.New(fmt.Sprintf("bad: %v", qm))
		return
	}

	if !qm.KnownLeader {
		err = errors.New(fmt.Sprintf("bad: %v", qm))
		return
	}

	return
}

func (c *consul) ACLCreate(query *api.ACLEntry) (acl *api.ACLEntry, err error) {
	id, wm, err := c.client.ACL().Create(query, nil)
	if err != nil {
		return
	}
	if wm.RequestTime == 0 {
		err = errors.New(fmt.Sprintf("bad: %v", wm))
		return
	}
	if id == "" {
		err = errors.New(fmt.Sprintf("invalid: %v", id))
		return
	}
	acl, _, err = c.client.ACL().Info(id, nil)
	if err != nil {
		return
	}
	return
}

func (c *consul) ACLUpdate(query *api.ACLEntry) (err error) {
	wm, err := c.client.ACL().Update(query, nil)
	if err != nil {
		return
	}
	if wm.RequestTime == 0 {
		err = errors.New(fmt.Sprintf("bad: %v", wm))
		return
	}
	return
}

func (c *consul) ACLDelete(token string) (err error) {
	wm, err := c.client.ACL().Destroy(token, nil)
	if err != nil {
		return
	}

	if wm.RequestTime == 0 {
		err = errors.New(fmt.Sprintf("bad: %v", wm))
		return
	}
	return
}
