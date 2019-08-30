/**
 * @Time : 2019-07-04 15:36
 * @Author : solacowa@gmail.com
 * @File : metrics
 * @Software: GoLand
 */

package pods

import (
	"github.com/kplcloud/request"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
	"fmt"
)

type jsonRes struct {
	Metrics []struct {
		Timestamp metav1.Time `json:"timestamp"` // todo 时间考虑处理一下 加8小时
		Value     int64       `json:"value"`
	} `json:"metrics"`
	LatestTimestamp metav1.Time `json:"latest_timestamp"`
}

type XYRes struct {
	X string `json:"x"`
	Y int64  `json:"y"`
}

func GetPodsMetrics(ns, name, httpUrl string, metrics chan map[string]interface{}) {
	if httpUrl == "" {
		httpUrl = "heapster.kube-system"
	}
	if res, err := getMetrics(ns, name, httpUrl, "cpu/usage_rate"); err == nil {
		var cpuRes []XYRes
		for _, v := range res.Metrics {
			curTime := v.Timestamp.Local().In(time.Local).Unix()
			cpuRes = append(cpuRes, XYRes{
				X: time.Unix(curTime, 0).Format("2006-01-02 15:04:05"),
				Y: v.Value,
			})
		}
		var i int
		var currCpu int64
		if len(cpuRes)-1 >= 0 {
			i = len(cpuRes) - 1
			currCpu = cpuRes[i].Y
		}
		metrics <- map[string]interface{}{
			"cpu":      cpuRes,
			"curr_cpu": currCpu,
		}
	}

	if res, err := getMetrics(ns, name, httpUrl, "memory/usage"); err == nil {
		var mRes []XYRes
		for _, v := range res.Metrics {
			curTime := v.Timestamp.Local().In(time.Local).Unix()
			mRes = append(mRes, XYRes{
				X: time.Unix(curTime, 0).Format("2006-01-02 15:04:05"),
				Y: v.Value / 1024 / 1024, // Mi
			})
		}
		var i int
		var currMem int64
		if len(mRes)-1 >= 0 {
			i = len(mRes) - 1
			currMem = mRes[i].Y
		}
		metrics <- map[string]interface{}{
			"memory":      mRes,
			"curr_memory": currMem,
		}
	}
	close(metrics)
	return
}

func getMetrics(ns string, podName string, httpUrl, metricName string) (res jsonRes, err error) {
	var uri string
	if podName == "" {
		uri = fmt.Sprintf("%s/api/v1/model/namespaces/%s/metrics/%s",
			httpUrl, ns, metricName)
	}else{
		uri = fmt.Sprintf("%s/api/v1/model/namespaces/%s/pods/%s/metrics/%s",
			httpUrl, ns, podName, metricName)
	}


	req := request.NewRequest(uri, "GET")
	// 集群内部不需要代理先注释掉
	//if c.config.GetString("server", "http_proxy") != "" {
	//	dialer := &net.Dialer{
	//		Timeout:   time.Duration(30 * time.Second),
	//		KeepAlive: time.Duration(30 * time.Second),
	//	}
	//	req.HttpClient(&http.Client{
	//		Transport: &http.Transport{
	//			Proxy: func(_ *http.Request) (*url.URL, error) {
	//				return url.Parse(c.config.GetString("server", "http_proxy"))
	//			},
	//			DialContext: dialer.DialContext,
	//		},
	//	})
	//}
	if err = req.Do().Into(&res); err != nil {
		return
	}

	return
}
