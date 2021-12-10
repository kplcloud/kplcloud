/**
 * @Time : 2021/12/6 10:13 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitcache "github.com/icowan/kit-cache"
	"github.com/kplcloud/kplcloud/src/encode"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util"
	"github.com/pkg/errors"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	// HandleTerminalSession 处理客户端发来的ws建立请求
	HandleTerminalSession(session sockjs.Session)
	Token(ctx context.Context, userId, clusterId int64, namespace, svcName, podName string) (res tokenResult, err error)
}

type service struct {
	traceId, appKey string
	logger          log.Logger
	k8sClient       kubernetes.K8sClient
	repository      repository.Repository
	sessionTimeout  int64
	cache           kitcache.Service
}

func (s *service) Token(ctx context.Context, userId, clusterId int64, namespace, svcName, podName string) (res tokenResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	cluster, err := s.repository.Cluster(ctx).Find(ctx, clusterId)
	if err != nil {
		_ = level.Warn(logger).Log("repository.Cluster", "Find", "err", err.Error())
		err = encode.ErrClusterNotfound.Error()
		return
	}

	// 通过pod返查项目，得项目权限
	pod, err := s.k8sClient.Do(ctx).CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		_ = level.Warn(logger).Log("k8sClient.Do.CoreV1.Pods", "Get", "err", err.Error())
		err = encode.ErrPodNotfound.Wrap(err)
		return
	}

	pods, err := s.getPods(ctx, namespace, svcName)
	if err != nil {
		_ = level.Error(logger).Log("s", "getPods", "err", err.Error())
		err = encode.ErrPodNotfound.Wrap(err)
		return
	}
	if len(pods) == 0 {
		err = encode.ErrTerminalPodsNotfound.Error()
		return
	}

	for _, v := range pods {
		res.Pods = append(res.Pods, v.Name)
	}

	timeout := time.Duration(s.sessionTimeout) * time.Second
	expAt := time.Now().Add(timeout).Unix()

	// 创建声明
	claims := kpljwt.ArithmeticTerminalClaims{
		UserId:    userId,
		Cluster:   cluster.Name,
		Namespace: namespace,
		PodName:   podName,
		Container: "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//生成token
	tk, err := token.SignedString([]byte(kpljwt.GetJwtKey()))
	if err != nil {
		_ = level.Error(logger).Log("token", "SignedString", "err", err.Error())
		return
	}

	// userId 入redis，验证的时候查一下 时间为600秒
	_ = s.cache.Set(ctx, fmt.Sprintf("terminal:tk:%s", util.Md5Str(tk)), userId, timeout)

	var errMsg string
	for _, val := range pod.Status.ContainerStatuses {
		if svcName != "filebeat" && svcName != "istio-proxy" && val.Name == pod.Labels[types.LabelAppName.String()] && val.Ready == false {
			errMsg = val.State.Waiting.Message
			if errMsg == "" {
				errMsg = val.State.Waiting.Reason
			}
		}
	}

	var containers []string
	for _, val := range pod.Spec.Containers {
		containers = append(containers, val.Name)
	}

	res.Namespace = namespace
	res.Cluster = cluster.Name
	res.ErrMsg = errMsg
	res.SessionId = tk
	res.PodName = podName
	res.Containers = containers
	res.ServiceName = svcName
	res.Phase = string(pod.Status.Phase)
	res.HostIp = pod.Status.HostIP
	res.PodIp = pod.Status.PodIP
	if pod.Status.StartTime != nil {
		res.StartTime = pod.Status.StartTime.Time
	}

	return
}

func (s *service) getPods(ctx context.Context, namespace, svcName string) (res []v1.Pod, err error) {
	for _, name := range []string{
		types.LabelAppName.String(),
		"k8s-app",
		"app",
	} {
		pods, err := s.k8sClient.Do(ctx).CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set{
				name: svcName,
			}).String(),
		})
		if err != nil {
			err = encode.ErrPodNotfound.Wrap(err)
			return res, err
		}
		if len(pods.Items) == 0 {
			continue
		}
		return pods.Items, nil
	}

	return nil, err
}

func (s *service) HandleTerminalSession(session sockjs.Session) {
	var (
		buf string
		err error
		msg Message
		//terminalSession Session
	)

	if buf, err = session.Recv(); err != nil {
		_ = level.Error(s.logger).Log("handleTerminalSession", "can't Recv:", "err", err.Error())
		return
	}
	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		_ = level.Error(s.logger).Log("handleTerminalSession", "can't UnMarshal", "err", err.Error(), "buf", buf)
		return
	}

	if msg.Op != "bind" {
		_ = level.Error(s.logger).Log("handleTerminalSession: expected 'bind' message, got:", buf)
		return
	}

	var tr Result
	if err := json.Unmarshal([]byte(msg.Data), &tr); err != nil {
		_ = level.Error(s.logger).Log("handleTerminalResult", "can't UnMarshal", "err", err.Error())
		return
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.ContextKeyClusterName, tr.Cluster)

	// 验证token、权限、过期时间之类的
	err = s.checkShellToken(tr.Cluster, tr.Namespace, tr.PodName, tr.Container, tr.SessionId)
	if err != nil {
		_ = level.Error(s.logger).Log("http.status", http.StatusBadRequest, "token", "not valid", "token", tr.Token, "err", err.Error())
		err = encode.ErrAuthTimeout.Wrap(err)
		return
	}

	ts := Session{
		id:            tr.SessionId,
		sockJSSession: session,
		sizeChan:      make(chan remotecommand.TerminalSize),
	}

	go s.waitForTerminal(s.k8sClient.Do(ctx), s.k8sClient.Config(ctx), ts, tr.Namespace, tr.PodName, tr.Container, "")
	//terminalSession.sockJSSession = session
	//terminalSessions.Set(msg.SessionID, terminalSession)
	//terminalSession.bound <- nil
	return
}

func (s *service) checkShellToken(cluster, namespace, podName, container, token string) error {
	var atc kpljwt.ArithmeticTerminalClaims
	tk, err := jwt.ParseWithClaims(token, &atc, kpljwt.JwtKeyFunc)
	if err != nil || tk == nil {
		_ = level.Error(s.logger).Log("jwt", "ParseWithClaims", "err", err)
		err = encode.ErrAuthTimeout.Wrap(err)
		return err
	}

	claim, ok := tk.Claims.(*kpljwt.ArithmeticTerminalClaims)
	if !ok {
		_ = level.Error(s.logger).Log("tk", "Claims", "err", ok)
		err = encode.ErrAccountASD.Error()
		return err
	}

	if !strings.EqualFold(claim.Cluster, cluster) {
		return encode.ErrAccountASD.Error()
	}
	if !strings.EqualFold(claim.Namespace, namespace) {
		return encode.ErrAccountASD.Error()
	}
	if !strings.EqualFold(claim.PodName, podName) {
		return encode.ErrAccountASD.Error()
	}
	//if !strings.EqualFold(claim.Container, container) {
	//	return encode.ErrAccountASD.Error()
	//}
	ctx := context.Background()
	tkMd5 := util.Md5Str(token)
	// userId
	_, err = s.cache.Get(ctx, fmt.Sprintf("terminal:tk:%s", tkMd5), nil)
	if err != nil {
		return encode.ErrAuthTimeout.Wrap(err)
	}

	// 拿到用户ID之后的操作

	return nil
}

func (s *service) waitForTerminal(k8sClient *k8s.Clientset, cfg *rest.Config, ts Session, namespace, podName, container, cmd string) {
	var err error
	validShells := []string{"bash", "sh"}

	if isValidShell(validShells, cmd) {
		cmds := []string{cmd}
		err = startProcess(k8sClient, cfg, cmds, ts, namespace, podName, container)
	} else {
		for _, testShell := range validShells {
			cmd := []string{testShell}
			if err = startProcess(k8sClient, cfg, cmd, ts, namespace, podName, container); err == nil {
				break
			}
		}
	}
	if err != nil {
		_ = level.Error(s.logger).Log("namespace", namespace, "pod", podName, "container", container, "cmd", cmd, "service", "waitForTerminal", "err", err.Error())
		_ = ts.Toast(err.Error())
		ts.Close(2, err.Error())
		return
	}
	_ = ts.Toast("Process exited")
	ts.Close(1, "Process exited")
}

func New(logger log.Logger, traceId, appKey string, k8sClient kubernetes.K8sClient, repository repository.Repository, cacheSvc kitcache.Service, terminalSessionTimeout int64) Service {
	return &service{
		traceId,
		appKey,
		logger,
		k8sClient,
		repository,
		terminalSessionTimeout,
		cacheSvc,
	}
}

func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

// 开始建立ws连接
func startProcess(k8sClient *k8s.Clientset, cfg *rest.Config, cmd []string, ptyHandler PtyHandler, namespace, podName, container string) error {
	req := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: container,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(cfg, http.MethodPost, req.URL())
	if err != nil {
		err = errors.Wrap(err, "remotecommand.NewSPDYExecutor")
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})

	if err != nil {
		err = errors.Wrap(err, "exec.Stream")
		return err
	}

	return nil
}
