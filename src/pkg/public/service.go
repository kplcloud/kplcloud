package public

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/util/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/url"
	"strings"
	"time"
)

var (
	ErrProjectFind             = errors.New("项目查询失败")
	ErrProjectDeploymentK8sGet = errors.New("项目查询失败,请检查Kubernetes是否存在")
	ErrProjectToken            = errors.New("项目Token校验错误")
	ErrProjectEventNotPush     = errors.New("EventName not push")
	ErrProjectDeploymentGet    = errors.New("项目Deployment获取错误")
	ErrProjectMemberGet        = errors.New("build用户获取错误,可能不存在")
	ErrProjectJenkinsBuild     = errors.New("调用Jenkins Build错误")
)

type Service interface {
	// gitlab post过来的数据
	GitPost(ctx context.Context, namespace, name, token, keyWord, branch string, req gitlabHook) (err error)

	// prometheus 的告警数据
	PrometheusAlert(ctx context.Context, req *prometheusAlerts) error

	// 获取配置信息
	Config(ctx context.Context) (res map[string]interface{}, err error)
}

type service struct {
	logger     log.Logger
	cf         *config.Config
	amqpClient amqpClient.AmqpClient
	k8sClient  kubernetes.K8sClient
	jenkins    jenkins.Jenkins
	repository repository.Repository
}

func NewService(logger log.Logger,
	cf *config.Config,
	amqpClient amqpClient.AmqpClient,
	k8sClient kubernetes.K8sClient,
	jenkins jenkins.Jenkins,
	repository repository.Repository) Service {
	return &service{
		logger,
		cf,
		amqpClient,
		k8sClient,
		jenkins,
		repository}
}

func (c *service) PrometheusAlert(ctx context.Context, req *prometheusAlerts) (err error) {
	//接收到报警后，存入Mq,由消息分发中心接管消息处理流程
	type alarm struct {
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Name      string    `json:"name"`
		Namespace string    `json:"namespace"`
		Desc      string    `json:"desc"`
		CreatedAt time.Time `json:"created_at"`
	}

	var name, ns string

	name = req.GetName(&req.prom.CommonLabels)
	ns = req.GetNamespace(&req.prom.CommonLabels)

	alerts := req.Alerts()

	if name == "" {
		for _, v := range alerts {
			name = req.GetName(&v.Labels)
			if name != "" {
				break
			}
		}
	}

	if ns == "" {
		for _, v := range alerts {
			ns = req.GetNamespace(&v.Labels)
			if ns != "" {
				break
			}
		}
	}

	data := alarm{
		Title:     req.GetAlertName(),
		Content:   req.String(),
		Name:      name,
		Namespace: ns,
		Desc:      req.GetDesc(),
		CreatedAt: time.Now(),
	}

	b, _ := json.Marshal(data)

	//存入mq
	defer func() {
		if err := c.amqpClient.PublishOnQueue(amqpClient.AlarmTopic, func() []byte {
			return []byte(b)
		}); err != nil {
			_ = level.Error(c.logger).Log("amqpClient", "PublicAlarmQueue", "err", err.Error())
		}
	}()

	return
}

func (c *service) GitPost(ctx context.Context, namespace, name, token, keyWord, branch string, req gitlabHook) (err error) {
	project, err := c.repository.Project().FindByNsName(namespace, name)
	if err != nil {
		_ = level.Error(c.logger).Log("projectRepository", "FindByNsName", "err", err.Error())
		return ErrProjectFind
	}

	deployment, err := c.k8sClient.Do().AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("Deployments", "Get", "err", err.Error())
		return ErrProjectDeploymentK8sGet
	}

	if deployment.ObjectMeta.GetUID() != types.UID(token) {
		_ = level.Error(c.logger).Log("dep", deployment.ObjectMeta.GetUID(), "token", token)
		return ErrProjectToken
	}

	if req.EventName != "push" {
		_ = level.Error(c.logger).Log("req", "EventName", "err", "not push")
		return ErrProjectEventNotPush
	}

	//var build bool
	var email string
	for _, val := range req.Commits {
		if strings.Contains(val.Message, keyWord) {
			//build = true
			email = val.Author.Email
			break
		}
	}

	//if !build {
	//	_ = level.Error(c.logger).Log("req", "build", "err", "not build")
	//	return
	//}

	pt, err := c.repository.ProjectTemplate().FindByProjectId(project.ID, repository.Deployment)
	if err != nil {
		_ = level.Error(c.logger).Log("templateRepository", "FindByProjectId", "err", err.Error())
		return ErrProjectDeploymentGet
	}

	if branch == "" {
		branch = pt.FieldStruct.Branch
	}

	refs := strings.Split(req.Ref, "/")
	if refs[len(refs)-1] != branch {
		_ = level.Error(c.logger).Log("refs", refs[len(refs)-1], "branch", branch)
		return
	}

	member, err := c.repository.Member().Find(email)
	if err != nil {
		_ = level.Error(c.logger).Log("memberRepository", "Find", "err", err.Error())
		return ErrProjectMemberGet
	}

	_ = level.Info(c.logger).Log("email", email, "project", project.ID, "member", member.Username)

	jobName := project.Name + "." + project.Namespace
	var tagName string
	params := url.Values{
		"TAGNAME": []string{tagName},
	}
	// jenkins build
	if err = c.jenkins.Build(jobName, params); err != nil {
		_ = level.Error(c.logger).Log("jenkins", "Build", "err", err.Error())
		return ErrProjectJenkinsBuild
	}

	return
}

func (c *service) Config(ctx context.Context) (res map[string]interface{}, err error) {
	gitAddr := c.cf.GetString("git", "git_addr")
	res = map[string]interface{}{
		"git_type": c.cf.GetString("git", "git_type"),
		"git_addr": helper.GitUrl(gitAddr),
		"domain":   c.cf.GetString("server", "domain_suffix"),
	}
	return
}
