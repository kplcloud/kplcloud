/**
 * @Time: 2021/12/5 17:17
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package pod

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	lineReadLimit int64 = 5000    // 最多读取行数
	byteReadLimit int64 = 5000000 // 最多读取字节
)

// Service Pod 相关操作
type Service interface {
	// GetLog 获取容器日志
	GetLog(ctx context.Context, clusterId int64, namespace, podName, container string, previous bool) (res LogDetails, err error)
	// DownloadLog 下载pod 容器日志
	DownloadLog(ctx context.Context, clusterId int64, namespace, podName, container string, previous bool) (res io.ReadCloser, err error)
}

type service struct {
	traceId    string
	logger     log.Logger
	repository repository.Repository
	k8sClient  kubernetes.K8sClient
}

func (s *service) GetLog(ctx context.Context, clusterId int64, namespace, podName, container string, previous bool) (res LogDetails, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))
	pod, err := s.k8sClient.Do(ctx).CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do.CoreV1.Pods", "Get", "err", err.Error())
		err = encode.ErrPodNotfound.Wrap(err)
		return
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

	res, err = s.getLogDetails(ctx, namespace, podName, container, logSelector, previous)
	if err != nil {
		_ = level.Error(logger).Log("c", "getLogDetails", "err", err.Error())
		return
	}

	return
}

func (s *service) DownloadLog(ctx context.Context, clusterId int64, namespace, podName, container string, previous bool) (res io.ReadCloser, err error) {
	panic("implement me")
}

func (s *service) getLogDetails(ctx context.Context, ns, podId, container string, logSelector *Selection, usePreviousLogs bool) (res LogDetails, err error) {
	logOptions := s.mapToLogOptions(container, logSelector, usePreviousLogs)
	rawLogs, err := s.readRawLogs(ctx, ns, podId, logOptions)
	if err != nil {
		return
	}

	deatils := s.constructLogDetails(podId, rawLogs, container, logSelector)
	return deatils, nil
}

func (s *service) mapToLogOptions(container string, logSelector *Selection, previous bool) *v1.PodLogOptions {
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

func (s *service) readRawLogs(ctx context.Context, namespace, podID string, logOptions *v1.PodLogOptions) (
	string, error) {
	readCloser, err := s.openStream(ctx, namespace, podID, logOptions)
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

func (s *service) openStream(ctx context.Context, namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return s.k8sClient.Do(ctx).CoreV1().RESTClient().Get().Namespace(namespace).
		Name(podID).Resource("pods").SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).
		Stream(ctx)
}

func (s *service) constructLogDetails(podID string, rawLogs string, container string, logSelector *Selection) LogDetails {
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
	return LogDetails{
		Info:      info,
		Selection: logSelection,
		LogLines:  logLines,
	}
}

func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == End && linesLoaded >= lineReadLimit)
}

func New(logger log.Logger, traceId string, client kubernetes.K8sClient, repository repository.Repository) Service {
	return &service{
		logger: logger, k8sClient: client,
		repository: repository,
		traceId:    traceId,
	}
}
