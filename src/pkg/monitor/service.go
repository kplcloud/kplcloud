/**
 * @Time : 2019-07-29 14:22
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package monitor

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/request"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"sync"
	"time"
)

type Service interface {
	// 获取全局网络请求
	QueryNetwork(ctx context.Context) (data []map[string]interface{}, err error)

	// 全局OPS
	Ops(ctx context.Context) (rs interface{}, err error)

	// 全局内存有CPU消耗
	Metrics(ctx context.Context) (map[string]interface{}, error)
}

var (
	ErrMonitorPrometheusGet = errors.New("Prometheus获取数据出错: ")
	ErrMonitorNodesList     = errors.New("节点列表信息获取错误")
)

type service struct {
	logger                     log.Logger
	config                     *config.Config
	k8sClient                  kubernetes.K8sClient
	repository                 repository.Repository
	prometheusUrl, heapsterUrl string
}

func NewService(logger log.Logger, config *config.Config,
	k8sClient kubernetes.K8sClient, store repository.Repository) Service {

	prometheusUrl := config.GetString("server", "prometheus_url")
	if prometheusUrl == "" {
		prometheusUrl = "http://prometheus.istio-system:9090"
	}
	prometheusUrl += "/api/v1/query_range"

	heapsterUrl := config.GetString("server", "heapster_url")
	if heapsterUrl == "" {
		heapsterUrl = "http://heapster.kube-system"
	}

	return &service{logger, config, k8sClient,
		store,
		prometheusUrl, heapsterUrl}
}

func (c *service) Metrics(ctx context.Context) (map[string]interface{}, error) {
	nodes, err := c.k8sClient.Do().CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Nodes", "List", "err", err.Error())
		return nil, ErrMonitorNodesList
	}

	language := map[string]int64{}
	var memory int64
	var cpu int64
	var usedMem int64
	var usedCpu int64

	for _, node := range nodes.Items {
		memory += node.Status.Capacity.Memory().Value()
		cpu += node.Status.Capacity.Cpu().Value()

		metricsRes := c.getNodeMetrics(node.Name, []string{
			"memory/usage",
			"cpu/usage_rate",
		})
		for _, v := range metricsRes {
			if v["memory/usage"] != nil {
				usedMem += v["memory/usage"][len(v["memory/usage"])-1].Y
			}
			if v["cpu/usage_rate"] != nil {
				usedCpu += v["cpu/usage_rate"][len(v["cpu/usage_rate"])-1].Y
			}
		}
	}

	memory = memory / 1024 / 1024 / 1024
	usedMem = usedMem / 1024 / 1024 / 1024

	if res, err := c.repository.Project().CountLanguage(); err == nil {
		var total int64
		for _, v := range res {
			if v.Language != "Total" {
				continue
			}
			total = v.Total
			break
		}
		for _, v := range res {
			if v.Language == "Total" {
				continue
			}
			language[v.Language] = v.Total * 100 / total
		}
		language["Total"] = total
	}

	return map[string]interface{}{
		"language": language,
		"memory": map[string]interface{}{
			"total": memory,
			"used":  usedMem,
		},
		"cpu": map[string]interface{}{
			"total": cpu,
			"used":  usedCpu / 1000,
		},
	}, nil
}

type res struct {
	X string `json:"x"`
	Y int64  `json:"y"`
}

type jsonRes struct {
	Metrics []struct {
		Timestamp metav1.Time `json:"timestamp"` // 时间考虑处理一下 加8小时
		Value     int64       `json:"value"`
	} `json:"metrics"`
}

func (c *service) getNodeMetrics(name string, metricsNames []string) []map[string][]res {
	var wg sync.WaitGroup
	wg.Add(len(metricsNames))
	rs := make([]map[string][]res, 0)

	for _, v := range metricsNames {
		var path = "/api/v1/model/nodes/" + name + "/metrics/" + v
		go func(path, v string) {
			if resp, err := c.getMetrics(path); err == nil {
				rs = append(rs, map[string][]res{
					v: resp,
				})
			} else {
				_ = level.Error(c.logger).Log("c", "getMetrics", "err", err.Error())
			}
			wg.Done()
		}(path, v)
	}

	wg.Wait()

	return rs
}

func (c *service) getMetrics(path string) (rs []res, err error) {
	var jsonrs jsonRes
	if err = request.NewRequest(c.heapsterUrl+path, "GET").Do().Into(&jsonrs); err != nil {
		return
	}

	for _, v := range jsonrs.Metrics {
		curTime := v.Timestamp.Local().In(time.Local).Unix()
		rs = append(rs, res{
			X: time.Unix(curTime, 0).Format("2006-01-02 15:04:05"),
			Y: v.Value,
		})
	}
	return
}

func (c *service) Ops(ctx context.Context) (rs interface{}, err error) {

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric struct{}        `json:"metric"`
				Values [][]interface{} `json:"values"`
			} `json:"result"`
		} `json:"data"`
	}

	if err = request.NewRequest(c.prometheusUrl, "GET").
		Param("query", "round(sum(irate(istio_requests_total{reporter=\"destination\"}[1m])), 0.001)").
		Param("start", strconv.FormatInt(time.Now().Unix()-900, 10)).
		Param("end", strconv.FormatInt(time.Now().Unix(), 10)).
		Param("step", "14").Do().Into(&resp); err != nil {
		err = errors.New(ErrMonitorPrometheusGet.Error() + err.Error())
		return
	}

	type xy struct {
		X string  `json:"x"`
		Y float64 `json:"y"`
	}
	var res []xy

	if len(resp.Data.Result) > 0 {
		for _, v := range resp.Data.Result[0].Values {
			var x float64 = v[0].(float64)
			var y string = v[1].(string)

			yy, _ := strconv.ParseFloat(y, 64)
			xx := strconv.FormatFloat(x, 'f', 0, 64)

			unix, _ := strconv.Atoi(xx)

			res = append(res, xy{
				X: time.Unix(int64(unix), 0).Format("2006/01/02 15:04:05"),
				Y: yy,
			})
		}
	}

	return res, nil
}

func (c *service) QueryNetwork(ctx context.Context) (data []map[string]interface{}, err error) {

	type res struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Metric struct{}        `json:"metric"`
				Values [][]interface{} `json:"values"`
			} `json:"result"`
			ResultType string `json:"resultType"`
		} `json:"data"`
	}

	var receiveRes res
	var transmitRes res

	query := request.NewRequest(c.prometheusUrl, "GET").
		Param("start", strconv.FormatInt(time.Now().Unix()-900, 10)).
		Param("end", strconv.FormatInt(time.Now().Unix(), 10)).
		Param("step", "14")
	{
		if err = query.Param("query", "sum(rate(container_network_receive_bytes_total{}[1m]))").
			Do().Into(&receiveRes); err != nil {
			err = errors.New(ErrMonitorPrometheusGet.Error() + err.Error())
			return
		}
	}
	{
		if err = query.Param("query", "sum(rate(container_network_transmit_bytes_total{}[1m]))").
			Do().Into(&transmitRes); err != nil {
			err = errors.New(ErrMonitorPrometheusGet.Error() + err.Error())
			return
		}
	}

	for key, val := range receiveRes.Data.Result[0].Values {
		receive, _ := strconv.ParseFloat(val[1].(string), 64)
		transmit, _ := strconv.ParseFloat(transmitRes.Data.Result[0].Values[key][1].(string), 64)
		data = append(data, map[string]interface{}{
			"time":     time.Unix(int64(val[0].(float64)), 0).Format("2006/01/02 15:04:05"),
			"receive":  receive / 1024 / 1024,
			"transmit": transmit / 1024 / 1024,
		})
	}

	return
}
