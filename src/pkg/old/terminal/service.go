/**
 * @Time : 2019-06-27 18:05
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package terminal

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"strconv"
	"time"
)

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type TerminalSession struct {
	id            string
	sockJSSession sockjs.Session
	sizeChan      chan remotecommand.TerminalSize
}

type TerminalMessage struct {
	Op, Data, SessionID string
	Rows, Cols          uint16
}

type TerminalResult struct {
	SessionId string `json:"sessionId,omitempty"`
	Token     string `json:"token,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Pod       string `json:"pod,omitempty"`
	Container string `json:"container,omitempty"`
	Cmd       string `json:"cmd,omitempty"`
}

func (t TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	}
}

func (t TerminalSession) Read(p []byte) (int, error) {
	m, err := t.sockJSSession.Recv()
	if err != nil {
		return 0, err
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		fmt.Println(fmt.Sprintf("read msg (%s) form client error.%v", string(p), err))
		return 0, err
	}
	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{msg.Cols, msg.Rows}
		return 0, nil
	default:
		return 0, fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

func (t TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Op:   "stdout",
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = t.sockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t TerminalSession) Close(status uint32, reason string) {
	err := t.sockJSSession.Close(status, reason)
	fmt.Println(fmt.Sprintf("close socket (%s). %d, %s, %v", t.id, status, reason, err))

}

type IndexData struct {
	Namespace    string
	PodName      string
	Container    string
	ErrMsg       string
	SessionId    string
	Token        string
	BashStr      string
	TemplateFile string
}

var (
	ErrPodK8sGet         = errors.New("容器获取失败,该容器可能不存在")
	ErrSessionIdGenerate = errors.New("SessionId 生成失败,请刷新页面重试")
	//ErrConsoleTemplateFile = errors.New("控制台模版获取错误,请查询文件是否有权限获取")
)

type Service interface {
	//Attach(ctx context.Context, path string) http.Handler
	// 终端控制台页面
	Index(ctx context.Context, podName, container string) (*IndexData, error)

	// 处理客户端发来的ws建立请求
	HandleTerminalSession(session sockjs.Session)
}

type service struct {
	logger    log.Logger
	config    *config.Config
	k8sClient kubernetes.K8sClient
}

func NewService(logger log.Logger, config *config.Config, k8sClient kubernetes.K8sClient) Service {
	return &service{logger, config, k8sClient}
}

func (c *service) Index(ctx context.Context, podName, container string) (*IndexData, error) {
	project := ctx.Value(middleware.ProjectContext).(*types.Project)

	pod, err := c.k8sClient.Do().CoreV1().Pods(project.Namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Pods", "Get", "err", err.Error())
		return nil, ErrPodK8sGet
	}

	sessionId, err := genTerminalSessionId()
	if err != nil {
		_ = level.Error(c.logger).Log("genTerminalSessionId", "fail", "err", err.Error())
		return nil, ErrSessionIdGenerate
	}

	var errmsg string
	for _, val := range pod.Status.ContainerStatuses {
		if container != "filebeat" && container != "istio-proxy" && val.Name == project.Name && val.Ready == false {
			errmsg = val.State.Waiting.Message
			if errmsg == "" {
				errmsg = val.State.Waiting.Reason
			}
		}
	}

	token := generateToken(project.Namespace, podName, c.config.GetString("server", "app_key"))
	var bashStr = `starting container process caused "exec: \"bash\": executable file not found in $PATH": unknown`

	return &IndexData{
		Namespace:    project.Namespace,
		PodName:      podName,
		Container:    container,
		Token:        token,
		BashStr:      bashStr,
		SessionId:    sessionId,
		ErrMsg:       errmsg,
		TemplateFile: c.config.GetString("server", "http_static") + "/terminal.html",
	}, nil
}

/**
 * @Title websocket 连接
 */
//func (c *service) Attach(ctx context.Context, path string) http.Handler {
//	return sockjs.NewHandler(path, sockjs.DefaultOptions, c.HandleTerminalSession)
//}

func (c *service) HandleTerminalSession(session sockjs.Session) {
	var (
		buf string
		err error
		msg TerminalMessage
	)

	if buf, err = session.Recv(); err != nil {
		_ = level.Error(c.logger).Log("handleTerminalSession", "can't Recv:", "err", err.Error())
		return
	}

	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		_ = level.Error(c.logger).Log("handleTerminalSession", "can't UnMarshal", "err", err.Error(), "buf", buf)
		return
	}

	if msg.Op != "bind" {
		_ = level.Error(c.logger).Log("handleTerminalSession: expected 'bind' message, got:", buf)
		return
	}

	var tr TerminalResult
	if err := json.Unmarshal([]byte(msg.Data), &tr); err != nil {
		_ = level.Error(c.logger).Log("handleTerminalResult", "can't UnMarshal", "err", err.Error())
		return
	}

	err = c.checkShellToken(tr.Token, tr.Namespace, tr.Pod)
	if err != nil {
		_ = level.Error(c.logger).Log("http.status", http.StatusBadRequest, "token", "not valid", "token", tr.Token, "err", err.Error())
		return
	}
	ts := TerminalSession{
		id:            tr.SessionId,
		sockJSSession: session,
		sizeChan:      make(chan remotecommand.TerminalSize),
	}

	go WaitForTerminal(c.k8sClient.Do(), c.k8sClient.Config(), ts, tr.Namespace, tr.Pod, tr.Container, "")
	return
}

func (c *service) checkShellToken(token string, namespace string, podName string) error {
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

	appKey := c.config.GetString("server", "app_key")

	rawToken := namespace + podName + endTimeStr + appKey

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(rawToken))
	cipherToken := hex.EncodeToString(md5Ctx.Sum(nil))

	checkToken := string([]rune(cipherToken)[12:20]) + endTimeStr
	if checkToken != token {
		return errors.New("token not match")
	}
	return nil
}

func WaitForTerminal(k8sClient *k8s.Clientset, cfg *rest.Config, ts TerminalSession, namespace, pod, container, cmd string) {
	var err error
	validShells := []string{"bash", "sh"}

	if isValidShell(validShells, cmd) {
		cmds := []string{cmd}
		err = startProcess(k8sClient, cfg, cmds, ts, namespace, pod, container)
	} else {
		for _, testShell := range validShells {
			cmd := []string{testShell}
			if err = startProcess(k8sClient, cfg, cmd, ts, namespace, pod, container); err == nil {
				break
			}
		}
	}

	if err != nil {
		ts.Close(2, err.Error())
		return
	}

	ts.Close(1, "Process exited")
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
func startProcess(k8sClient *k8s.Clientset, cfg *rest.Config, cmd []string, ptyHandler PtyHandler, namespace, pod, container string) error {
	req := k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
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
		return err
	}

	return nil
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
func generateToken(namespace, pod, appKey string) string {
	endTime := time.Now().Unix() + 60*10
	rawTokenKey := namespace + pod + strconv.FormatInt(endTime, 10) + appKey
	md5Hash := md5.New()
	md5Hash.Write([]byte(rawTokenKey))
	cipher := md5Hash.Sum(nil)
	cipherStr := hex.EncodeToString(cipher)
	return cipherStr[12:20] + strconv.FormatInt(endTime, 10)
}
