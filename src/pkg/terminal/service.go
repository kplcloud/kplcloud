/**
 * @Time : 2021/12/6 10:13 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package terminal

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/pkg/errors"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	// HandleTerminalSession 处理客户端发来的ws建立请求
	HandleTerminalSession(session sockjs.Session)
	Token(ctx context.Context, clusterId int64, namespace, podName, container string) (res tokenResult, err error)
}

type service struct {
	traceId, appKey string
	logger          log.Logger
	k8sClient       kubernetes.K8sClient
	repository      repository.Repository
}

func (s *service) Token(ctx context.Context, clusterId int64, namespace, podName, container string) (res tokenResult, err error) {
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

	// TODO: 通过pod返查项目，得项目权限
	//s.repository.Application(ctx)

	sessionId, _ := genTerminalSessionId()
	var errMsg string
	for _, val := range pod.Status.ContainerStatuses {
		if container != "filebeat" && container != "istio-proxy" && val.Name == pod.Labels["app"] && val.Ready == false {
			errMsg = val.State.Waiting.Message
			if errMsg == "" {
				errMsg = val.State.Waiting.Reason
			}
		}
	}

	token := generateToken(cluster.Name, namespace, podName, s.appKey)
	var bashStr = `starting container process caused "exec: \"bash\": executable file not found in $PATH": unknown`

	res.Token = token
	res.BashStr = bashStr
	res.Namespace = namespace
	res.Container = container
	res.Cluster = cluster.Name
	res.ErrMsg = errMsg
	res.SessionId = sessionId
	res.PodName = podName

	return
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

	// TODO: 验证token、权限、过期时间之类的
	//err = s.checkShellToken(tr.Cluster, tr.Namespace, tr.Pod, tr.Container, tr.Token)
	//if err != nil {
	//	_ = level.Error(s.logger).Log("http.status", http.StatusBadRequest, "token", "not valid", "token", tr.Token, "err", err.Error())
	//	return
	//}
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
	endTimeRaw := []rune(token)
	var endTime int64
	var endTimeStr string
	var err error

	if len(endTimeRaw) > 8 {
		endTimeStr = string(endTimeRaw[8:])
		endTime, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			return err
		}
	}
	ntime := time.Now().Unix()

	if ntime > endTime {
		return errors.New("token time expired")
	}

	rawToken := namespace + podName + endTimeStr + s.appKey

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(rawToken))
	cipherToken := hex.EncodeToString(md5Ctx.Sum(nil))

	checkToken := string([]rune(cipherToken)[12:20]) + endTimeStr
	if checkToken != token {
		return errors.New("token not match")
	}
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
		ts.Close(2, err.Error())
		return
	}

	ts.Close(1, "Process exited")
}

func New(logger log.Logger, traceId, appKey string, k8sClient kubernetes.K8sClient, repository repository.Repository) Service {
	return &service{
		traceId,
		appKey,
		logger,
		k8sClient,
		repository,
	}
}

func genTerminalSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

// token生成规则
// 1. 600秒时限，平台appkey，并进行md5加密
// 2. 取生成的32位加密字符串第12-20位，于unixtime进行拼接生成token
func generateToken(clusterName, namespace, pod, appKey string) string {
	endTime := time.Now().Unix() + 60*10
	rawTokenKey := fmt.Sprintf("%s:%s:%s:%s:%s", clusterName, namespace, pod, appKey, strconv.FormatInt(endTime, 10))
	//rawTokenKey := namespace + pod + strconv.FormatInt(endTime, 10) + appKey
	md5Hash := md5.New()
	md5Hash.Write([]byte(rawTokenKey))
	cipher := md5Hash.Sum(nil)
	cipherStr := hex.EncodeToString(cipher)
	return cipherStr[12:20] + strconv.FormatInt(endTime, 10)
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

	fmt.Println(namespace, podName, container)

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
