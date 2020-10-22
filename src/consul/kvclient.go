/**
 * @Time : 2019/7/17 3:47 PM
 * @Author : yuntinghu1003@gmail.com
 * @File : kvclient
 * @Software: GoLand
 */

package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/icowan/config"
)

type KVClient interface {
	KVList(prefix string) (pairs api.KVPairs, err error)
	KVGet(prefix string) (pairs *api.KVPair, err error)
	KVPut(key string, value string) (err error)
	KVDelete(key string) (err error)
	KVDeleteTree(prefix string) (err error)
}

type kvClient struct {
	client *api.Client
}

func NewKVClient(cf *config.Config, tokenId string) (KVClient, error) {
	conf := api.DefaultConfig()
	conf.Address = cf.GetString("consul", "consul_addr")
	conf.Token = tokenId

	cli, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &kvClient{client: cli}, nil
}

func (c *kvClient) KVList(prefix string) (pairs api.KVPairs, err error) {
	pairs, meta, err := c.client.KV().List(prefix, &api.QueryOptions{})
	if err != nil {
		err = errors.New(fmt.Sprintf("err: %v", err))
		return
	}
	if meta.LastIndex == 0 {
		err = errors.New(fmt.Sprintf("unexpected value: %#v", meta))
		return
	}
	return
}

func (c *kvClient) KVGet(prefix string) (pair *api.KVPair, err error) {
	pair, meta, err := c.client.KV().Get(prefix, nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("err: %v", err))
		return
	}
	if meta.LastIndex == 0 {
		err = errors.New(fmt.Sprintf("unexpected value: %#v", meta))
		return
	}
	return
}

func (c *kvClient) KVPut(key string, value string) (err error) {
	_, err = c.client.KV().Put(&api.KVPair{Key: key, Value: []byte(value)}, nil)
	return
}

func (c *kvClient) KVDelete(key string) (err error) {
	_, err = c.client.KV().Delete(key, nil)
	return
}

func (c *kvClient) KVDeleteTree(prefix string) (err error) {
	_, err = c.client.KV().DeleteTree(prefix, nil)
	return
}
