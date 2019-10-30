package pod

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/config"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/pkg/hooks"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/pods"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	ErrProjectPodsList   = errors.New("该项目的pods列表获取错误，请查看是否存在")
	ErrPodLogGet         = errors.New("容器日志获取错误,请重试")
	ErrPodGet            = errors.New("容器获取错误,可能不存在,请重试")
	ErrPodDelete         = errors.New("pod删除错误,请重试")
	ErrPodDownloadLogGet = errors.New("容器日志流获取错误,请重试")
	ErrPodDeploymentGet  = errors.New("项目获取错误,请重试")
)

var (
	lineReadLimit int64 = 5000
	byteReadLimit int64 = 5000000
)

type Service interface {
	// pod详情页数据
	Detail(ctx context.Context, podName string) (res map[string]interface{}, err error)

	// Pods 列表
	ProjectPods(ctx context.Context) (res []map[string]interface{}, err error)

	// 获取容器日志
	GetLog(ctx context.Context, podName, container string, previous bool) (res *LogDetails, err error)

	// 下载pod 容器日志
	DownloadLog(ctx context.Context, podName, container string, previous bool) (res io.ReadCloser, err error)

	// 删除pod
	Delete(ctx context.Context, podName string) (err error)

	// pods的内存及CPU使用
	PodsMetrics(ctx context.Context) (res map[string]interface{}, err error)

	PodsNetwork(ctx context.Context) (res map[string]interface{}, err error)
}

type service struct {
	logger       log.Logger
	k8sClient    kubernetes.K8sClient
	config       *config.Config
	hookQueueSvc hooks.ServiceHookQueue
}

func NewService(logger log.Logger, k8sClient kubernetes.K8sClient, config *config.Config, hookQueueSvc hooks.ServiceHookQueue) Service {
	return &service{logger, k8sClient, config, hookQueueSvc}
}

func (c *service) PodsNetwork(ctx context.Context) (res map[string]interface{}, err error) {
	return
}

func (c *service) PodsMetrics(ctx context.Context) (res map[string]interface{}, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	name := ctx.Value(middleware.NameContext).(string)

	dep, err := c.k8sClient.Do().AppsV1().Deployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return nil, ErrPodDeploymentGet
	}

	var selectorKey, selectorVal string
	for key, val := range dep.Spec.Selector.MatchLabels {
		selectorKey = key
		selectorVal = val
	}

	podList, err := c.k8sClient.Do().CoreV1().Pods(ns).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal),
	})

	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "List", "err", err.Error())
		return nil, ErrProjectPodsList
	}
	var currMemory int64
	var currCpu int64
	memory := make(map[string]int64)
	cpu := make(map[string]int64)

	for _, pod := range podList.Items {
		metrics := make(chan map[string]interface{})
		go pods.GetPodsMetrics(pod.Namespace, pod.Name, c.config.GetString("server", "heapster_url"), metrics)
		for {
			data, ok := <-metrics
			if !ok {
				break
			}
			if data["memory"] != nil {
				for _, v := range data["memory"].([]pods.XYRes) {
					memory[v.X] += v.Y
				}
			}
			if data["cpu"] != nil {
				for _, v := range data["cpu"].([]pods.XYRes) {
					cpu[v.X] += v.Y
				}
			}

			if data["curr_memory"] != nil {
				currMemory += data["curr_memory"].(int64)
			} else {
				currMemory += 0
			}
			if data["curr_cpu"] != nil {
				currCpu += data["curr_cpu"].(int64)
			} else {
				currCpu += 0
			}
		}
	}

	var mem []map[string]interface{}
	var cc []map[string]interface{}

	for k, v := range memory {
		mem = append(mem, map[string]interface{}{
			"x": k,
			"y": v,
		})
	}
	for k, v := range cpu {
		cc = append(cc, map[string]interface{}{
			"x": k,
			"y": v,
		})
	}

	return map[string]interface{}{
		"memory":      mem,
		"curr_memory": currMemory,
		"curr_cpu":    currCpu,
		"cpu":         cc,
	}, nil
}

func (c *service) Delete(ctx context.Context, podName string) (err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	name := ctx.Value(middleware.NameContext).(string)

	var gracePeriodSeconds = int64(0) // 1 强制删除
	var policy = metav1.DeletePropagationBackground
	if err = c.k8sClient.Do().CoreV1().Pods(ns).Delete(podName, &metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &policy,
	}); err != nil {
		_ = level.Error(c.logger).Log("Pods", "Delete", "err", err.Error())
		return ErrPodDelete
	}

	//操作记录  event
	go func() {
		if err := c.hookQueueSvc.SendHookQueue(ctx,
			repository.RebootEvent,
			name, ns,
			fmt.Sprintf("重启容器: %v.%v, 容器名称: %v", name, ns, podName)); err != nil {
			_ = level.Warn(c.logger).Log("hookQueueSvc", "SendHookQueue", "err", err.Error())
		}
	}()

	return
}

func (c *service) DownloadLog(ctx context.Context, podName, container string, previous bool) (res io.ReadCloser, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)

	logStream, err := c.k8sClient.Do().CoreV1().RESTClient().Get().
		Namespace(ns).
		Name(podName).
		Resource("pods").
		SubResource("log").
		VersionedParams(&v1.PodLogOptions{
			Container:  container,
			Follow:     false,
			Previous:   previous,
			Timestamps: false,
		}, scheme.ParameterCodec).Stream()
	if err != nil {
		_ = level.Error(c.logger).Log("pods.log", "Stream", "err", err.Error())
		return nil, ErrPodDownloadLogGet
	}

	return logStream, nil
}

func (c *service) GetLog(ctx context.Context, podName, container string, previous bool) (res *LogDetails, err error) {
	ns := ctx.Value(middleware.NamespaceContext).(string)
	pod, err := c.k8sClient.Do().CoreV1().Pods(ns).Get(podName, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "List", "err", err.Error())
		return nil, ErrPodGet
	}

	var refLineNum = 0
	var offsetFrom = 2000000000
	var offsetTo = 2000000100

	refTimestamp := NewestTimestamp

	logSelector := DefaultSelection

	logSelector = &Selection{
		ReferencePoint: LogLineId{
			LogTimestamp: LogTimestamp(refTimestamp),
			LineNum:      refLineNum,
		},
		OffsetFrom:      offsetFrom,
		OffsetTo:        offsetTo,
		LogFilePosition: "end",
	}

	if container == "" {
		container = pod.Spec.Containers[0].Name
	}

	result, err := c.getLogDetails(ns, podName, container, logSelector, previous)
	if err != nil {
		_ = level.Error(c.logger).Log("c", "getLogDetails", "err", err.Error())
		return nil, ErrPodLogGet
	}

	return result, nil
}

func (c *service) getLogDetails(ns, podId, container string, logSelector *Selection, usePreviousLogs bool) (*LogDetails, error) {
	logOptions := c.mapToLogOptions(container, logSelector, usePreviousLogs)
	rawLogs, err := c.readRawLogs(ns, podId, logOptions)
	if err != nil {
		return nil, err
	}

	deatils := c.constructLogDetails(podId, rawLogs, container, logSelector)
	return deatils, nil
}

func (c *service) mapToLogOptions(container string, logSelector *Selection, previous bool) *v1.PodLogOptions {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   previous,
		Timestamps: true,
	}

	if logSelector.LogFilePosition == Beginning {
		logOptions.LimitBytes = &byteReadLimit
	} else {
		logOptions.TailLines = &lineReadLimit
	}

	return logOptions
}

func (c *service) readRawLogs(namespace, podID string, logOptions *v1.PodLogOptions) (
	string, error) {
	readCloser, err := c.openStream(namespace, podID, logOptions)
	if err != nil {
		return err.Error(), nil
	}

	defer func() {
		_ = readCloser.Close()
	}()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c *service) openStream(namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return c.k8sClient.Do().CoreV1().RESTClient().Get().Namespace(namespace).
		Name(podID).Resource("pods").
		SubResource("log").VersionedParams(logOptions, scheme.ParameterCodec).
		Stream()
}

func (c *service) constructLogDetails(podID string, rawLogs string, container string, logSelector *Selection) *LogDetails {
	parsedLines := ToLogLines(rawLogs)
	logLines, fromDate, toDate, logSelection, lastPage := parsedLines.SelectLogs(logSelector)

	readLimitReached := isReadLimitReached(int64(len(rawLogs)), int64(len(parsedLines)), logSelector.LogFilePosition)
	truncated := readLimitReached && lastPage

	info := LogInfo{
		PodName:       podID,
		ContainerName: container,
		FromDate:      fromDate,
		ToDate:        toDate,
		Truncated:     truncated,
	}
	return &LogDetails{
		Info:      info,
		Selection: logSelection,
		LogLines:  logLines,
	}
}

func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == End && linesLoaded >= lineReadLimit)
}

func (c *service) ProjectPods(ctx context.Context) (res []map[string]interface{}, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	dep, err := c.k8sClient.Do().AppsV1().Deployments(project.Namespace).Get(project.Name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return nil, ErrPodDeploymentGet
	}

	var selectorKey, selectorVal string
	for key, val := range dep.Spec.Selector.MatchLabels {
		selectorKey = key
		selectorVal = val
	}

	podList, err := c.k8sClient.Do().CoreV1().Pods(project.Namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", selectorKey, selectorVal),
	})

	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "List", "err", err.Error())
		return nil, ErrProjectPodsList
	}

	for _, pod := range podList.Items {
		var restartCount int32
		for _, container := range pod.Status.ContainerStatuses {
			if container.Name == project.Name {
				restartCount = container.RestartCount
				break
			}
		}
		res = append(res, map[string]interface{}{
			"name":          pod.Name,
			"node_name":     pod.Spec.NodeName,
			"status":        pod.Status.Phase,
			"restart_count": restartCount,
			"create_at":     pod.CreationTimestamp,
		})
	}

	return res, nil
}

func (c *service) Detail(ctx context.Context, podName string) (res map[string]interface{}, err error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	pod, err := c.k8sClient.Do().CoreV1().Pods(project.Namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "Get", "err", err.Error())
		return nil, ErrPodGet
	}

	metrics := make(chan map[string]interface{})
	go pods.GetPodsMetrics(pod.Namespace, pod.Name, c.config.GetString("server", "heapster_url"), metrics)

	var restartCount int32
	for _, v := range pod.Status.ContainerStatuses {
		if v.Name == project.Name {
			restartCount = v.RestartCount
			break
		}
	}

	metricData := map[string]interface{}{}
	for {
		data, ok := <-metrics
		if !ok {
			break
		}
		if m, ok := data["memory"]; ok {
			metricData["memory"] = m
		}
		if u, ok := data["cpu"]; ok {
			metricData["cpu"] = u
		}
		if m, ok := data["curr_memory"]; ok {
			metricData["curr_memory"] = m
		}
		if u, ok := data["curr_cpu"]; ok {
			metricData["curr_cpu"] = u
		}
	}

	res = map[string]interface{}{
		"name":          pod.Name,
		"node_name":     pod.Spec.NodeName,
		"status":        pod.Status.Phase,
		"restart_count": restartCount,
		"create_at":     pod.CreationTimestamp,
		"metrics":       metricData,
		"pod":           pod,
	}

	return
}
