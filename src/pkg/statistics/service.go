/**
 * @Time : 2019-07-29 11:30
 * @Author : soupzhb@gmail.com
 * @File : service.go
 * @Software: GoLand
 */

package statistics

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/paginator"
)

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrNotInGroupRefused = errors.New("没有此组权限")
	ErrProjectRefused    = errors.New("获取组内项目失败")
	ErrBuildRefused      = errors.New("获取build列表失败")
	ErrGroupBuildRefused = errors.New("获取groupBuild列表失败")
)

type Service interface {
	// 获取 build 统计
	Build(ctx context.Context, req buildRequest) (res map[string]interface{}, err error)
}

type service struct {
	logger log.Logger
	cf     *config.Config
	store  repository.Repository
}

func NewService(logger log.Logger, cf *config.Config, store repository.Repository) Service {
	return &service{
		logger: logger,
		cf:     cf,
		store:  store,
	}
}

func (c *service) Build(ctx context.Context, req buildRequest) (res map[string]interface{}, err error) {
	//ns := ctx.Value(middleware.NamespaceContext).(string)
	isAdmin := ctx.Value(middleware.IsAdmin).(bool)
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	var project []*types.Project

	if req.GroupID > 0 {
		if isAdmin == false {
			//非超管，需要验证组权限
			isIn, err := c.store.Groups().IsInGroup(int64(req.GroupID), memberId)

			if isIn == false {
				_ = level.Error(c.logger).Log("statistics", "c.group.IsInGroup", "err", err.Error())
				return nil, ErrNotInGroupRefused
			}

			project, err = c.store.Project().GetProjectByGroupId(int64(req.GroupID))
			if err != nil {
				_ = level.Error(c.logger).Log("statistics", "c.project.GetProjectByGroupId", "err", err.Error())
				return nil, ErrProjectRefused
			}

		}
	} else {
		if isAdmin == false {
			//获取权限内的项目
			project, err = c.store.Project().GetProjectByMidAndNs(memberId, req.Namespace)
			if err != nil {
				_ = level.Error(c.logger).Log("statistics", "c.project.GetProjectByMidAndNs", "err", err.Error())
				return nil, ErrProjectRefused
			}

		} else {
			project, err = c.store.Project().GetProjectByNs(req.Namespace)
			if err != nil {
				_ = level.Error(c.logger).Log("statistics", "c.project.GetProjectByNs", "err", err.Error())
				return nil, ErrProjectRefused
			}
		}
	}

	var projectNames []string

	if len(project) > 0 {
		for _, v := range project {
			projectNames = append(projectNames, v.Name)
		}
	}

	reqBuild := repository.StatisticsRequest{
		Namespace:   req.Namespace,
		Name:        req.Name,
		ProjectName: projectNames,
		STime:       req.STime,
		ETime:       req.ETime,
		BuildID:     0,
		GitType:     "",
	}

	count, _ := c.store.Build().CountStatistics(reqBuild)
	p := paginator.NewPaginator(req.Page, req.Limit, int(count))

	list, err := c.store.Build().FindStatisticsOffsetLimit(reqBuild, p.Offset(), req.Limit)

	if err != nil {
		_ = level.Error(c.logger).Log("statistics", "c.build.FindStatisticsOffsetLimit", "err", err.Error())
		return nil, ErrBuildRefused
	}

	type ResponseList struct {
		Id        int64  `json:"id"`
		Namespace string `json:"namespace"`
		Ns        string `json:"ns"`
		Name      string `json:"name"`
		Success   int    `json:"success"`
		Failure   int    `json:"failure"`
		Aborted   int    `json:"aborted"`
		Rollback  int    `json:"rollback"`
		Building  int    `json:"building"`
	}

	var ResponseListMap []ResponseList

	for _, v := range list {
		//每个项目的build记录
		ress, err := c.store.Build().GetGroupByBuilds(repository.StatisticsRequest{
			Namespace:   v.Namespace,
			Name:        v.Name,
			ProjectName: projectNames,
			STime:       req.STime,
			ETime:       req.ETime,
			BuildID:     0,
			GitType:     "",
		}, "namespace,name,status")

		if err != nil {
			_ = level.Error(c.logger).Log("statistics", "c.build.GetGroupByBuilds", "err", err.Error())
			return nil, ErrGroupBuildRefused
		}
		list := ResponseList{
			Id:        v.ID,
			Namespace: v.Namespace,
			Ns:        v.Namespace,
			Name:      v.Name,
			Success:   0,
			Failure:   0,
			Aborted:   0,
			Rollback:  0,
			Building:  0,
		}
		for _, v1 := range ress {
			if v1.Build.Status.String == "success" || v1.Build.Status.String == "SUCCESS" {
				list.Success = v1.Count
			}

			if v1.Build.Status.String == "failure" || v1.Build.Status.String == "FAILURE" {
				list.Failure = v1.Count
			}

			if v1.Build.Status.String == "aborted" || v1.Build.Status.String == "ABORTED" {
				list.Aborted = v1.Count
			}

			if v1.Build.Status.String == "rollback" || v1.Build.Status.String == "ROLLBACK" {
				list.Rollback = v1.Count
			}

			if v1.Build.Status.String == "building" || v1.Build.Status.String == "BUILDING" {
				list.Building = v1.Count
			}
		}
		ResponseListMap = append(ResponseListMap, list)
	}

	type AllBuilds struct {
		X     string `json:"x"`
		Y     int    `json:"y"`
		Color string `json:"color"`
	}

	//所有符合条件项目的总数
	allBuilds, err := c.store.Build().GetGroupByBuilds(repository.StatisticsRequest{
		Namespace:   req.Namespace,
		Name:        req.Name,
		ProjectName: projectNames,
		STime:       req.STime,
		ETime:       req.ETime,
		BuildID:     0,
		GitType:     "",
	}, "status")

	if err != nil {
		_ = level.Error(c.logger).Log("statistics", "c.build.GetGroupByBuilds-allBuilds", "err", err.Error())
		return nil, ErrGroupBuildRefused
	}

	var aBuilds []AllBuilds
	for _, v2 := range allBuilds {
		if v2.Build.Status.String == "success" || v2.Build.Status.String == "SUCCESS" {
			a := AllBuilds{
				X:     "success",
				Y:     v2.Count,
				Color: "#36c966", //绿色
			}
			aBuilds = append(aBuilds, a)
		}

		if v2.Build.Status.String == "failure" || v2.Build.Status.String == "FAILURE" {
			a := AllBuilds{
				X:     "failure",
				Y:     v2.Count,
				Color: "#f2526f", //红色
			}
			aBuilds = append(aBuilds, a)
		}

		if v2.Build.Status.String == "aborted" || v2.Build.Status.String == "ABORTED" {
			a := AllBuilds{
				X:     "aborted",
				Y:     v2.Count,
				Color: "#cbceb5bd", // 灰色
			}
			aBuilds = append(aBuilds, a)
		}

		if v2.Build.Status.String == "rollback" || v2.Build.Status.String == "ROLLBACK" {
			a := AllBuilds{
				X:     "rollback",
				Y:     v2.Count,
				Color: "#d08993bd", // 紫色
			}
			aBuilds = append(aBuilds, a)
		}

		if v2.Build.Status.String == "building" || v2.Build.Status.String == "BUILDING" {
			a := AllBuilds{
				X:     "building",
				Y:     v2.Count,
				Color: "#fad214", //黄色
			}
			aBuilds = append(aBuilds, a)
		}
	}

	var returnData = map[string]interface{}{
		"list": ResponseListMap,
		"all":  aBuilds,
		"page": map[string]interface{}{
			"total":     count,
			"pageTotal": p.PageTotal(),
			"pageSize":  req.Limit,
			"page":      req.Page,
		},
	}

	return returnData, nil
}
