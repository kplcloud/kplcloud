/**
 * @Time : 2019-07-25 15:11
 * @Author : soupzhb@gmail.com
 * @File : endpoint.go
 * @Software: GoLand
 */

package workspace

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	utilpods "github.com/kplcloud/kplcloud/src/util/pods"
)

var (
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrProclaimParamsRefused = errors.New("参数校验未通过.")
)

type Service interface {
	// 工作台 当前空间的CPU及内存使用
	Metrice(ctx context.Context, ns string) (res map[string]interface{}, err error)
	Active(ctx context.Context) (res []map[string]interface{}, err error)
}

type service struct {
	build  repository.BuildRepository
	logger log.Logger
	cf     *config.Config
}

func NewService(logger log.Logger, cf *config.Config, build repository.BuildRepository) Service {
	return &service{
		logger: logger,
		cf:     cf,
		build:  build,
	}
}

func (c *service) Metrice(ctx context.Context, ns string) (res map[string]interface{}, err error) {
	metrics := make(chan map[string]interface{})
	go utilpods.GetPodsMetrics(ns, "", c.cf.GetString("server", "heapster_url"), metrics)
	var memory, currMemory, currCpu, cpu interface{}
	for {
		data, ok := <-metrics
		if !ok {
			break
		}
		if m, ok := data["memory"]; ok {
			memory = m
		}
		if u, ok := data["cpu"]; ok {
			cpu = u
		}
		if m, ok := data["curr_memory"]; ok {
			currMemory = m
		}
		if u, ok := data["curr_cpu"]; ok {
			currCpu = u
		}
	}

	return map[string]interface{}{
		"memory":      memory,
		"cpu":         cpu,
		"curr_memory": currMemory,
		"curr_cpu":    currCpu,
	}, nil
}

func (c *service) Active(ctx context.Context) (res []map[string]interface{}, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	builds, _ := c.build.FindOffsetLimit(ns, "", 0, 5)

	for _, v := range builds {
		res = append(res, map[string]interface{}{
			"id":        v.ID,
			"updatedAt": v.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			"user": map[string]interface{}{
				"name":   v.Member.Username,
				"avatar": "https://niu.yirendai.com/kpl-logo-blue.png",
			},
			"group": map[string]interface{}{
				"name":   v.Namespace,
				"avatar": "https://niu.yirendai.com/kpl-logo-blue.png",
			},
			"project": map[string]interface{}{
				"name": v.Name,
				//"link": "#/project/detail/" + val.Project,
			},
			"template": "在 @{group} Build 了应用 @{project}",
		})
	}
	return
}
