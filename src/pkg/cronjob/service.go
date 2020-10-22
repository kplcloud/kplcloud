package cronjob

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/yaml.v2"
	v13 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	resource2 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sTypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/icowan/config"
	amqpClient "github.com/kplcloud/kplcloud/src/amqp"
	"github.com/kplcloud/kplcloud/src/jenkins"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/configmapyaml"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/kplcloud/kplcloud/src/util/helper"
	"github.com/kplcloud/kplcloud/src/util/paginator"
)

const (
	version = "apps/v1"
	kind    = "CronJob"
)

var (
	ErrCronJobNameExists         = errors.New("定时任务名已存在")
	ErrCreateCronJobFailed       = errors.New("创建定时任务失败")
	ErrGetTempByKindFailed       = errors.New("获取CronJob yaml模板失败")
	ErrCreateCronYamlFailed      = errors.New("创建CronJob yaml模板失败")
	ErrCronJobInfoFailed         = errors.New("从K8S获取定时任务信息失败")
	ErrStrconvFailed             = errors.New("格式转换失败")
	ErrIfIsInGroupFailed         = errors.New("判断是否在组里查询失败")
	ErrUserIsNotInGroup          = errors.New("用户不在该组里")
	ErrCronJobCountFailed        = errors.New("定时任务总数获取失败")
	ErrCronJobListFailed         = errors.New("定时任务列表数据获取失败")
	ErrGetCronJobFromK8sFailed   = errors.New("从k8s拉取定时任务信息失败")
	ErrCronJobNotExists          = errors.New("定时任务不存在")
	ErrDelCronJobFailed          = errors.New("删除定时任务失败")
	ErrK8sDelCronJobFailed       = errors.New("从k8s删除job失败")
	ErrConfigMapDelFailed        = errors.New("从configMap删除数据失败")
	ErrCronJobUpdateFailed       = errors.New("定时任务数据库更新失败")
	ErrBuildJenkinsJob           = errors.New("Jenkins Build Job错误:")
	ErrCronJobClientCreateFailed = errors.New("cronjob client 创建失败")
	ErrBuildQueuePublish         = errors.New("构建入列出错了,请联系管理员")
	ErrCreateConfMapFailed       = errors.New("创建config map 失败")
	ErrConfMapYamlFailed         = errors.New("更新远程configMapYaml")
	ErrGetTemplateFailed         = errors.New("获取模板信息错误")
	ErrFileBeatYamlFailed        = errors.New("FileBeatYaml 错误")
	ErrGetConfigDataFailed       = errors.New("获取configData信息错误")
	ErrCreateConfMapDataFailed   = errors.New("创建configData信息失败")
	ErrUpdateConfMapDataFailed   = errors.New("更新configData信息失败")
	ErrExchangeCronJobTemp       = errors.New("转换模板 失败")
	ErrJenkinsCreateJob          = errors.New("jenkins创建job失败")
	ErrJenkinsBuildFailed        = errors.New("jenkins构建失败")
)

type Service interface {
	// 添加定时任务
	AddCronJob(ctx context.Context, acj addCronJob) error

	// 定时任务列表
	List(ctx context.Context, name string, ns string, group string, page int, limit int) (map[string]interface{}, error)

	// 定时任务详情
	Detail(ctx context.Context, name string, ns string) (res *DetailReturnData, err error)

	// 定时任务修改
	Put(ctx context.Context, name string, acj addCronJob) error

	// 删除定时任务
	Delete(ctx context.Context, name string, ns string) error

	// 删除所有
	DeleteJobAll(ctx context.Context, ns string) error

	// 定时任务处理队列
	CronJobQueuePop(ctx context.Context, data string) error

	// 更新日志
	UpdateLog(ctx context.Context, req cronJobLogUpdate) error

	// 手动触发
	Trigger(ctx context.Context, name, ns string) (err error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	jenkins    jenkins.Jenkins
	k8sClient  kubernetes.K8sClient
	amqpClient amqpClient.AmqpClient
	repository repository.Repository
}

// /api/v1/cronjob/operations/invite-jon-p2p/trigger
func (c *service) Trigger(ctx context.Context, name, ns string) (err error) {
	logger := log.With(c.logger, "request-id", ctx.Value("request-id"), "namespace", ns, "name", name)
	cronJob, err := c.k8sClient.Do().BatchV1beta1().CronJobs(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do", "BatchV1beta1", "CronJobs", "Get", "err", err.Error())
		return
	}

	annotations := make(map[string]string)
	annotations["cronjob.kubernetes.io/instantiate"] = "manual"

	labels := make(map[string]string)
	for k, v := range cronJob.Spec.JobTemplate.Labels {
		labels[k] = v
	}

	var newJobName string
	if len(cronJob.Name) < 42 {
		newJobName = cronJob.Name + "-manual-" + rand.String(3)
	} else {
		newJobName = cronJob.Name[0:41] + "-manual-" + rand.String(3)
	}

	jobToCreate := &v13.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        newJobName,
			Namespace:   ns,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	_, err = c.k8sClient.Do().BatchV1().Jobs(ns).Create(jobToCreate)
	if err != nil {
		_ = level.Error(logger).Log("k8sClient.Do", "BatchV1beta1", "Jobs", "Create", "err", err.Error())
		return err
	}

	return nil
}

type Param struct {
	Username  string
	UserToken string
	GitAddr   string
	GitToken  string
	GitType   string
	Branch    string
}

func (c *service) AddCronJob(ctx context.Context, acj addCronJob) (err error) {

	// 判断是否重名
	_, isExists := c.repository.CronJob().GetCronJobByNameAndNs(acj.Name, acj.Namespace)
	if !isExists {
		_ = level.Error(c.logger).Log("cronjob", "AddCronJob cronjob name s exists  ")
		return ErrCronJobNameExists
	}

	memberId := ctx.Value(middleware.UserIdContext).(int64)

	var cronjobModelArgs, _ = json.Marshal(acj.Args)
	cronJobCreate, err := c.repository.CronJob().Create(&types.Cronjob{
		Name:        acj.Name,
		Namespace:   acj.Namespace,
		Schedule:    acj.Schedule,
		Image:       acj.Image,
		GitPath:     acj.GitPath,
		GitType:     acj.GitType,
		ConfMapName: acj.ConfMap,
		Args:        string(cronjobModelArgs),
		LogPath:     acj.LogPath,
		AddType:     acj.AddType,
		MemberID:    memberId,
	})

	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "AddCronJob create cronjob failed ", "err", err.Error())
		return ErrCreateCronJobFailed
	}

	defer func(jobId int64) {
		// 如果后面的流程失败了需要删除当前这个job 所以 id 还是需要的
		if err != nil {
			_ = level.Error(c.logger).Log("cronjob", "AddCronJob failed ", "err", err.Error())
			if err = c.repository.CronJob().Delete(jobId); err != nil {
				_ = level.Error(c.logger).Log("cronjob", "AddCronJob delete cronjob failed ", "err", err.Error())
			}
		}
	}(cronJobCreate.ID)

	template, err := c.repository.Template().GetTemplateByKind("CronJob")
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "AddCronJob get template by kind failed ", "err", err.Error())
		return ErrGetTempByKindFailed
	}

	image := acj.Image
	if acj.AddType == "Script" {
		image = acj.Namespace + "/" + acj.Name + ":" + image
	}
	param := map[string]interface{}{
		"name":        acj.Name,
		"namespace":   acj.Namespace,
		"image":       image,
		"schedule":    acj.Schedule,
		"args":        acj.Args,
		"configMap":   acj.ConfMap, // 默认configMap为空
		"logPath":     acj.LogPath, // 默认日志采集目录为空
		"isConfigMap": "",          //首次创建不需要configMap
	}
	exchangeTemplate, err := encode.EncodeTemplate("cronjob", template.Detail, param)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "AddCronJob create cronjob yaml template failed ", "err", err.Error())
		return ErrCreateCronYamlFailed
	}

	if acj.AddType == "Script" {
		var gitAddr, gitPath, gits string

		if acj.GitPath != "" {
			gitAddr = acj.GitPath
			gitPath = helper.GitName(gitAddr)
			_, _, gits = parseGitAddr(gitAddr)
		}

		// 1. jenkins create project
		// 2. k8s create cronjob
		// 3. insert db
		var jenkinsParam jenkins.Params
		jenkinsParam.Name = acj.Name + "-cronjob"
		jenkinsParam.Namespace = acj.Namespace
		jenkinsParam.GitAddr = acj.GitPath
		jenkinsParam.GitType = cronJobCreate.GitType
		jenkinsParam.GitVersion = cronJobCreate.Image

		tmp, err := c.repository.Template().GetTemplateByKind("JenkinsCommand")
		if err != nil {
			_ = level.Error(c.logger).Log("cronjob", "AddCronJob get template by kind failed ", "err", err.Error())
			return ErrGetTempByKindFailed
		}

		command, err := encode.EncodeTemplate("JenkinsCommand", tmp.Detail, map[string]string{
			"app_name":  acj.Name,
			"git_name":  acj.Name + "-cronjob",
			"git_path":  gits + "/" + gitPath,
			"namespace": acj.Namespace,
		})

		if err != nil {
			_ = level.Error(c.logger).Log("cronjob", "AddCronJob encode template failed ", "err", err.Error())
			return err
		}

		jenkinsParam.Command = command

		if err := c.jenkins.CreateJobParams(jenkinsParam); err != nil {
			_ = level.Error(c.logger).Log("cronjob", "AddCronJob create project failed ", "err", err.Error())
			return ErrJenkinsCreateJob
		}

		var tagName string
		names := strings.Split(acj.Image, ":")
		tagName = names[len(names)-1]

		// jenkins build
		params := url.Values{
			"TAGNAME": []string{tagName},
		}
		name := acj.Name + "-cronjob"
		name += "." + acj.Namespace
		if err = c.jenkins.Build(name, params); err != nil {
			_ = level.Error(c.logger).Log("cronjob", "AddCronJob jenkins build failed ", "err", err.Error())
			return ErrJenkinsBuildFailed
		}

		memberId := ctx.Value(middleware.UserIdContext).(int64)

		//add builds database
		builds := types.Build{
			Name:      acj.Name + "-cronjob",
			Namespace: acj.Namespace,
			Version:   acj.Image,
			Status:    null.StringFrom(repository.Building),
			GitType:   null.StringFrom(acj.GitType),
			Address:   null.StringFrom(acj.GitPath),
			BuilderID: memberId,
		}

		if buildss, err := c.repository.Build().CreateBuild(&builds); err == nil {
			if err := c.amqpClient.PublishOnQueue(amqpClient.CronJobTopic, func() []byte {
				b, _ := json.Marshal(JenkinsCronjobData{
					Name:      acj.Name + "-cronjob",
					Namespace: acj.Namespace,
					BuildId:   buildss.ID,
					BuildTime: time.Now(),
				})
				return b
			}); err != nil {
				_ = level.Error(c.logger).Log("cronjob", "AddCronJob PublishOnQueue", "err", err.Error())
				//return ErrBuildQueuePublish
			}
		} else {
			_ = level.Error(c.logger).Log("cronjob", "AddCronJob create build", "err", err.Error())
		}

	}

	cronjob, err := convertToV1Beta(exchangeTemplate)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "AddCronJob create build", "err", err.Error())
		return
	}

	cron, err := c.k8sClient.Do().BatchV1beta1().CronJobs(acj.Namespace).Create(cronjob)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "AddCronJob k8s create failed", "err", err.Error())
		return
	}
	cron.Kind = kind
	cron.APIVersion = version

	// 脚本模式直接返回，不进行后续操作
	if acj.AddType != "Script" {
		return
	}

	return nil
}

func (c *service) CronJobQueuePop(ctx context.Context, data string) (err error) {
	if len(data) <= 0 {
		return nil
	}
	defer func() {
		if err != nil {
			if e := c.amqpClient.PublishOnQueue(amqpClient.CronJobTopic, func() []byte {
				return []byte(data)
			}); e != nil {
				_ = level.Error(c.logger).Log("amqpClient", "PublishOnQueue", "err", err.Error())
			}
		}
		time.Sleep(time.Second * 2)
	}()

	err = c.handleBuildCronJob(ctx, data)
	if err != nil {
		var dat *JenkinsCronjobData
		if er := json.Unmarshal([]byte(data), &dat); er != nil {
			_ = level.Error(c.logger).Log("cronjob", "CronJobQueuePop", "err", er.Error())
		}
		time.Sleep(4 * time.Second)
	}
	return
}

func (c *service) handleBuildCronJob(ctx context.Context, data string) (err error) {

	var dat JenkinsCronjobData
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		return err
	}

	buildsInfo, err := c.repository.Build().FindById(dat.Namespace, dat.Name, dat.BuildId)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "CronJobQueuePop", "err", err.Error())
		return err
	}

	var build jenkins.Build
	job, err := c.jenkins.GetJob(dat.Name + "." + dat.Namespace)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "CronJobQueuePop", "err", err.Error())
		return err
	}
	if buildsInfo.BuildID.Int64 > 0 {
		build, err = c.jenkins.GetBuild(job, int(buildsInfo.BuildID.Int64))
	} else {
		build, err = c.jenkins.GetLastBuild(job)
	}

	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "CronJobQueuePop", "err", err.Error())
		return err
	}

	resBody, err := c.jenkins.GetBuildConsoleOutput(build)
	if err != nil {
		_ = level.Error(c.logger).Log("jenkins", "GetBuildConsoleOutput", "err", err.Error())
		return err
	}

	id, _ := strconv.ParseInt(build.Id, 10, 64)
	buildsInfo.BuildID = null.IntFrom(id)
	buildsInfo.Status = null.StringFrom(build.Result)
	buildsInfo.Output = null.StringFrom(string(resBody))
	if buildsInfo.Status == null.StringFrom("") {
		buildsInfo.Status = null.StringFrom(repository.Building)
	}
	buildsInfo.BuildTime = null.TimeFrom(time.Now())

	go func(buildsInfo *types.Build) {
		if err = c.repository.Build().Update(buildsInfo); err != nil {
			_ = level.Error(c.logger).Log("cronjob", "CronJobQueuePop", "err", err.Error())
		}
	}(&buildsInfo)

	if build.Building || build.Result == repository.Building {
		return errors.New("build...")
	}

	//TODO:: 通知
	return nil
}

func (c *service) getOnePull(name string, ns string) (res map[string]interface{}, pods []map[string]interface{}, events []map[string]interface{}, cronjobYaml *v1beta1.CronJob, err error) {
	cronjobInfo, err := c.k8sClient.Do().BatchV1beta1().CronJobs(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "getOnePull get cronjob info failed ", "err", err.Error())
		return nil, nil, nil, &v1beta1.CronJob{}, ErrCronJobInfoFailed
	}
	cronjobInfo.Kind = kind
	cronjobInfo.APIVersion = version
	cronjobYaml = cronjobInfo

	res = map[string]interface{}{
		"name":              cronjobInfo.Name,
		"namespace":         cronjobInfo.Namespace,
		"status":            cronjobInfo.Status,
		"cron":              cronjobInfo.Spec,
		"schedule":          cronjobInfo.Spec.Schedule,
		"suspend":           cronjobInfo.Spec.Suspend,            //是否挂起
		"lastScheduleTime":  cronjobInfo.Status.LastScheduleTime, //最近调度时间
		"creationTimestamp": cronjobInfo.CreationTimestamp,
		"active":            len(cronjobInfo.Status.Active), //活跃中
	}

	ops := new(metav1.ListOptions)

	for _, p := range cronjobInfo.Status.Active {
		ops.LabelSelector = "job-name=" + p.Name
		list, err := c.k8sClient.Do().CoreV1().Pods(ns).List(*ops)
		if err != nil {
			_ = level.Error(c.logger).Log("cronjob", "getOnePull get pods list failed ", "err", err.Error())
			continue
		}
		for _, val := range list.Items {
			podInfo, err := c.k8sClient.Do().CoreV1().Pods(ns).Get(val.Name, v12.GetOptions{})
			if err != nil {
				_ = level.Error(c.logger).Log("cronjob", "getOnePull get pods by name  failed ", "err", err.Error())
				continue
			}

			pods = append(pods, map[string]interface{}{
				"name":       podInfo.Name,
				"node_name":  podInfo.Spec.NodeName,
				"status":     podInfo.Status.Phase,
				"create_at":  podInfo.CreationTimestamp,
				"uid":        val.Labels["controller-uid"],
				"containers": podInfo.Spec.Containers,
			})
		}
	}

	eventList, err := c.k8sClient.Do().EventsV1beta1().Events(ns).List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "getOnePull get event list failed ", "err", err.Error())
		return nil, nil, nil, &v1beta1.CronJob{}, err
	}
	for _, e := range eventList.Items {
		if e.Regarding.Name != name {
			continue
		}
		events = append(events, map[string]interface{}{
			"note":            e.Note,
			"component":       e.DeprecatedSource.Component,
			"deprecatedCount": e.DeprecatedCount,
			"firstTime":       e.DeprecatedFirstTimestamp,
			"lastTime":        e.DeprecatedLastTimestamp,
			"kind":            e.Regarding.Kind,
		})
	}
	return
}

func (c *service) List(ctx context.Context, name string, ns string, group string, page int, limit int) (map[string]interface{}, error) {
	isAdmin := ctx.Value(middleware.IsAdmin).(bool)
	memberId := ctx.Value(middleware.UserIdContext).(int64)

	groupInt := 0
	if group != "" && group != "undefined" {
		var err error
		groupInt, err = strconv.Atoi(group)
		if err != nil {
			_ = level.Error(c.logger).Log("cronjob", "List", "err", err.Error())
			return nil, ErrStrconvFailed
		}
		// 如果不是自己组,不能看
		if !isAdmin {
			res, err := c.repository.Groups().IsInGroup(int64(groupInt), memberId)
			if err != nil {
				_ = level.Error(c.logger).Log("cronjob", "List", "err", err.Error())
				return nil, ErrIfIsInGroupFailed
			}
			if !res {
				// 不让看
				_ = level.Error(c.logger).Log("cronjob", "List ", "err", "User is not in the group")
				return nil, ErrUserIsNotInGroup
			}
		}
	}

	cnt, err := c.repository.CronJob().CronJobCountWithGroup(name, ns, int64(groupInt))
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "List ", "err", err.Error())
		return nil, ErrCronJobCountFailed
	}
	p := paginator.NewPaginator(page, limit, int(cnt))

	cronjobs, err := c.repository.CronJob().CronJobPaginateWithGroup(name, ns, int64(groupInt), p.Offset(), limit)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "List ", "err", err.Error())
		return nil, ErrCronJobListFailed
	}

	var res []map[string]interface{}
	for _, v := range cronjobs {
		cjInfo, _, _, _, err := c.getOnePull(v.Name, v.Namespace)

		if err != nil {
			_ = level.Error(c.logger).Log("cronjob", "List ", "err", err.Error())
			return nil, ErrGetCronJobFromK8sFailed
		}

		res = append(res, map[string]interface{}{
			"name":               v.Name,
			"namespace":          v.Namespace,
			"add_type":           v.AddType,
			"status":             cjInfo["status"].(v1beta1.CronJobStatus),
			"cron":               cjInfo["cron"].(v1beta1.CronJobSpec),
			"schedule":           cjInfo["schedule"].(string),
			"suspend":            cjInfo["suspend"].(*bool),                 //是否挂起
			"last_schedule_time": cjInfo["lastScheduleTime"].(*metav1.Time), //最近调度时间
			"active":             cjInfo["active"].(int),                    //活跃中
			"created_at":         cjInfo["creationTimestamp"].(metav1.Time),
		})
	}

	var returnData = map[string]interface{}{
		"list": res,
		"page": map[string]interface{}{
			"total":     cnt,
			"pageTotal": p.PageTotal(),
			"pageSize":  limit,
			"page":      p.Page(),
		},
	}
	return returnData, nil
}

func (c *service) Detail(ctx context.Context, name string, ns string) (res *DetailReturnData, err error) {
	cronjobInfo, isExists := c.repository.CronJob().GetCronJobByNameAndNs(name, ns)
	if isExists {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", "cronjob not exists")
		return nil, ErrCronJobNotExists
	}

	var data = new(DetailReturnData)
	argsArr := strings.Split(cronjobInfo.Args, ",")
	json.Unmarshal([]byte(cronjobInfo.Args), &data.Args)
	data.Name = cronjobInfo.Name
	data.Namespace = cronjobInfo.Namespace
	data.Schedule = cronjobInfo.Schedule
	data.GitType = cronjobInfo.GitType
	data.GitPath = cronjobInfo.GitPath
	data.Image = cronjobInfo.Image
	if cronjobInfo.Suspend == 1 {
		data.Suspend = true
	} else {
		data.Suspend = false
	}
	data.Active = int64(cronjobInfo.Active)
	data.LastSchedule = cronjobInfo.LastSchedule
	data.ConfMapName = cronjobInfo.ConfMapName
	data.LogPath = cronjobInfo.LogPath
	data.AddType = cronjobInfo.AddType
	data.CronjobInfo, data.CronjobPods, data.CronjobEvents, data.CronjobYaml, err = c.getOnePull(name, ns)

	data.Command = strings.Replace(strings.Replace(argsArr[2], "\"", "", -1), "]", "", -1)
	return data, err
}

func (c *service) Put(ctx context.Context, name string, acj addCronJob) error {

	cronjobInfo, isExists := c.repository.CronJob().GetCronJobByNameAndNs(name, acj.Namespace)
	if isExists {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", "cronjob not exists")
		return ErrCronJobNotExists
	}

	oldImage := cronjobInfo.Image
	newImage := acj.Image

	//add cronjob table

	var cronjobModelArgs, _ = json.Marshal(acj.Args)
	cronjobInfo.Name = acj.Name
	cronjobInfo.Namespace = acj.Namespace
	cronjobInfo.Schedule = acj.Schedule
	cronjobInfo.Image = acj.Image
	cronjobInfo.GitPath = acj.GitPath
	cronjobInfo.GitType = acj.GitType
	//cronjobInfo.ConfMapName = nr.ConfMap
	cronjobInfo.Args = string(cronjobModelArgs)
	cronjobInfo.LogPath = acj.LogPath

	err := c.repository.CronJob().Update(cronjobInfo, cronjobInfo.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "PUT ", "err", "cronjob not exists")
		return ErrCronJobUpdateFailed
	}

	// 如果版本没变化，不build
	if acj.AddType == "Script" {
		var tagName string
		names := strings.Split(acj.Image, ":")
		tagName = names[len(names)-1]
		params := url.Values{
			"TAGNAME": []string{tagName},
		}
		name := acj.Name + "-cronjob"
		name += "." + acj.Namespace

		if oldImage != newImage {
			// jenkins build
			if err := c.jenkins.Build(name, params); err != nil {
				return errors.New(ErrBuildJenkinsJob.Error() + err.Error())
			}

			cronjobInfos, err := c.k8sClient.Do().BatchV1beta1().CronJobs(acj.Namespace).Get(acj.Name, metav1.GetOptions{})
			if err != nil {
				_ = level.Error(c.logger).Log("cronJob", "PUT ", "err", err.Error())
				return ErrGetCronJobFromK8sFailed
			}

			originImage := cronjobInfos.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
			s := strings.Split(originImage, ":")
			cronjobInfos.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = s[0] + ":" + newImage
			_, err = c.k8sClient.Do().BatchV1beta1().CronJobs(acj.Namespace).Update(cronjobInfos)
			if err != nil {
				_ = level.Error(c.logger).Log("cronJob", "PUT ", "err", err.Error())
				return err
			}
			memberId := ctx.Value(middleware.UserIdContext).(int64)

			//add builds database
			builds := types.Build{
				Name:      acj.Name + "-cronjob",
				Namespace: acj.Namespace,
				Version:   acj.Image,
				Status:    null.StringFrom(repository.Building),
				GitType:   null.StringFrom(acj.GitType),
				Address:   null.StringFrom(acj.GitPath),
				BuilderID: memberId,
			}

			if buildss, err := c.repository.Build().CreateBuild(&builds); err == nil {
				if err := c.amqpClient.PublishOnQueue(amqpClient.CronJobTopic, func() []byte {
					b, _ := json.Marshal(JenkinsCronjobData{
						Name:      acj.Name + "-cronjob",
						Namespace: acj.Namespace,
						BuildId:   buildss.ID,
						BuildTime: time.Now(),
					})
					return b
				}); err != nil {
					_ = level.Error(c.logger).Log("cronjob", "PUT PublishOnQueue", "err", err.Error())
					return ErrBuildQueuePublish
				}
			} else {
				_ = level.Error(c.logger).Log("cronjob", "PUT create build", "err", err.Error())
			}
		}
	}

	_, err = ExchangeCronJobTemp(acj.Name, acj.Namespace, c.k8sClient, c.config, c.repository)
	if err != nil {
		_ = level.Error(c.logger).Log("cronjob", "PUT exchangeCronJobTemp", "err", err.Error())
		return ErrExchangeCronJobTemp
	}

	// TODO:: 通知

	//exchange template
	return nil
}

func (c *service) Delete(ctx context.Context, name string, ns string) error {

	cronJob, isExists := c.repository.CronJob().GetCronJobByNameAndNs(name, ns)
	if isExists {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", "cronjob not exists")
		return ErrCronJobNotExists
	}

	err := c.repository.CronJob().Delete(cronJob.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
		return ErrDelCronJobFailed
	}

	// 删除 jenkins job
	// 删除 k8s cronJob, job, pod
	// 删除 k8s configmap
	// 写入历史

	job, err := c.jenkins.GetJob(cronJob.Name + "-cronjob." + cronJob.Namespace)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
		return ErrCronJobInfoFailed
	}

	if err = c.jenkins.DeleteJob(job); err != nil {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
		return ErrK8sDelCronJobFailed
	}

	//del cronJob
	//先获取CronJob,
	cronjobInfo, err := c.k8sClient.Do().BatchV1beta1().CronJobs(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
		return ErrGetCronJobFromK8sFailed
	}
	cronjobInfo.Kind = kind
	cronjobInfo.APIVersion = version

	// 获得pods  循环删除job,pod
	ops := new(metav1.ListOptions)
	for _, p := range cronjobInfo.Status.Active {
		ops.LabelSelector = "job-name=" + p.Name
		list, err := c.k8sClient.Do().CoreV1().Pods(ns).List(*ops)
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
			continue
		}
		for _, val := range list.Items {
			podInfo, err := c.k8sClient.Do().CoreV1().Pods(ns).Get(val.Name, v12.GetOptions{})
			if err != nil {
				_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
				continue
			}
			podName := podInfo.Name
			c.k8sClient.Do().CoreV1().Pods(ns).Delete(podName, &metav1.DeleteOptions{})
		}

		//del job
		err = c.k8sClient.Do().BatchV1().Jobs(ns).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
			continue
		}
	}

	//del cronJob
	cron, err := c.k8sClient.Do().BatchV1beta1().CronJobs(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
		return ErrCronJobInfoFailed
	}
	cron.Kind = kind
	cron.APIVersion = version

	err = c.k8sClient.Do().BatchV1beta1().CronJobs(ns).Delete(name, nil)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
		return ErrCronJobInfoFailed
	}

	//删除数据库configmap
	if configMap, isExists := c.repository.ConfigMap().Find(ns, name); !isExists {
		conf, err := c.k8sClient.Do().CoreV1().ConfigMaps(ns).Get(name, v12.GetOptions{})
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
			return err
		}
		conf.Kind = "ConfigMap"
		err = c.k8sClient.Do().CoreV1().ConfigMaps(ns).Delete(name, &v12.DeleteOptions{})
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
			return err
		}
		err = c.repository.ConfigMap().Delete(configMap.ID)
		err = c.repository.ConfigData().Delete(configMap.ID)
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob", "Delete ", "err", err.Error())
			return ErrConfigMapDelFailed
		}
	}

	return nil
}

func (c *service) DeleteJobAll(ctx context.Context, ns string) error {

	list, err := c.k8sClient.Do().BatchV1().Jobs(ns).List(metav1.ListOptions{})
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob", "DeleteJobAll ", "err", err.Error())
		return ErrCronJobInfoFailed
	}

	for _, v := range list.Items {
		if err = c.k8sClient.Do().BatchV1().Jobs(ns).Delete(v.Name, &metav1.DeleteOptions{}); err != nil {
			_ = level.Error(c.logger).Log("cronJob", "DeleteJobAll ", "err", err.Error())
			return ErrDelCronJobFailed
		}
	}
	return nil
}

func (c *service) UpdateLog(ctx context.Context, req cronJobLogUpdate) error {
	//获取定时任务信息
	nr, notFound := c.repository.CronJob().GetCronJobByNameAndNs(req.Name, req.Namespace)
	if notFound {
		_ = level.Error(c.logger).Log("cronJob.UpdateLog", "GetCronJobByNameAndNs ", "err", "data not found")
		return ErrCronJobNameExists
	}

	//更新日志采集目录
	nr.LogPath = req.LogPath
	err := c.repository.CronJob().Update(nr, nr.ID)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob.UpdateLog", "Update ", "err", err.Error())
		return ErrCronJobUpdateFailed
	}

	_, notFound = c.repository.ConfigMap().Find(req.Namespace, req.Name)
	if notFound {
		_, err = c.repository.ConfigMap().Create(&types.ConfigMap{
			Name:      req.Name,
			Namespace: req.Namespace,
			Type:      null.IntFrom(2),
		})
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob.UpdateLog", "Create ", "err", err.Error())
			return ErrCreateConfMapFailed
		}
	} else {
		err = configmapyaml.SyncConfigMapYaml(req.Namespace, req.Name, c.logger, c.k8sClient, c.repository)
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob.UpdateLog", "SyncConfigMapYaml ", "err", err.Error())
			return ErrConfMapYamlFailed
		}
	}

	//日志对应key-value
	log_key := "filebeat.yml"
	var fileBeatYaml string
	var fileBeat = new(helper.FileBeat)
	fileBeat.Namespace = nr.Namespace
	fileBeat.Name = nr.Name
	fileBeat.LogPath = "/" + strings.Trim(nr.LogPath, "/") + "/"
	template, err := c.repository.Template().FindByKindType("FileBeat")
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob.UpdateLog", "FindByKindType ", "err", err.Error())
		return ErrGetTemplateFailed
	}
	fileBeatYaml, err = helper.FileBeatYaml(fileBeat, template)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob.UpdateLog", "FileBeatYaml ", "err", err.Error())
		return ErrFileBeatYamlFailed
	}

	//更新configMapData
	configMap, _ := c.repository.ConfigMap().Find(req.Namespace, req.Name)
	var configMapId int64
	configMapId = configMap.ID

	//获取configData信息
	configDataLog, notFound := c.repository.ConfigData().FindByConfMapIdAndKey(configMapId, log_key)

	if notFound {
		err := c.repository.ConfigData().Create(&types.ConfigData{
			Key:         log_key,
			Value:       fileBeatYaml,
			ConfigMapID: configMap.ID,
		})
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob.UpdateLog", "CreateConfigMapData", "err", err.Error())
			return ErrCreateConfMapDataFailed
		}
	} else {
		// todo 之后实现
		err = c.repository.ConfigData().Update(configDataLog.ID, fileBeatYaml, "")
		if err != nil {
			_ = level.Error(c.logger).Log("cronJob.UpdateLog", "UpdateConfigMapData", "err", err.Error())
			return ErrUpdateConfMapDataFailed
		}
	}

	//更新configMap yaml
	err = configmapyaml.UpdateConfigMapYaml(configMapId, c.logger, c.k8sClient, c.repository)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob.UpdateLog", "UpdateConfigMapYaml", "err", err.Error())
		return ErrConfMapYamlFailed
	}

	_, err = ExchangeCronJobTemp(req.Name, req.Namespace, c.k8sClient, c.config, c.repository)
	if err != nil {
		_ = level.Error(c.logger).Log("cronJob.UpdateLog", "ExchangeCronJobTemp", "err", err.Error())
		return ErrExchangeCronJobTemp
	}

	//todo:: 通知

	return nil
}

func ExchangeCronJobTemp(name, namespace string,
	k8sClient kubernetes.K8sClient,
	confClient *config.Config,
	repository repository.Repository) (cronjobInfo *v1beta1.CronJob, err error) {
	var isConfigMap string

	nr, isExists := repository.CronJob().GetCronJobByNameAndNs(name, namespace)
	if isExists {
		return &v1beta1.CronJob{}, ErrCronJobInfoFailed
	}

	_, isExists = repository.ConfigMap().Find(namespace, name)
	if !isExists {
		isConfigMap = "ok"
	}

	envList, err := repository.ConfigEnv().GetConfigEnvByNameNs(name, namespace)
	if err != nil {
		return &v1beta1.CronJob{}, ErrCronJobInfoFailed
	}

	template, err := repository.Template().GetTemplateByKind("CronJob")
	if err != nil {
		return &v1beta1.CronJob{}, ErrGetTempByKindFailed
	}
	image := nr.Image
	if nr.AddType == "Script" {
		image = nr.Namespace + "/" + nr.Name + ":" + image
	}

	argStr := strings.Replace(strings.Replace(strings.Replace(nr.Args, "\"", "", -1), "]", "", -1), "[", "", -1)
	argsArr := strings.Split(argStr, ",")

	if nr.LogPath != "" { //如果有configMapData,则挂载configMap 文件
		isConfigMap = "ok"
	}

	param := map[string]interface{}{
		"name":        nr.Name,
		"namespace":   nr.Namespace,
		"image":       image,
		"schedule":    nr.Schedule,
		"args":        argsArr,
		"configMap":   nr.ConfMapName,
		"logPath":     nr.LogPath,
		"envs":        envList,
		"isConfigMap": isConfigMap,
	}

	exchangeTemplate, err := encode.EncodeTemplate("cronjob", template.Detail, param)
	if err != nil {
		return &v1beta1.CronJob{}, err
	}

	//更新cronJob Yaml ,先更新环境变量等基本信息
	cronjobInfo, err = patch(nr.Namespace, exchangeTemplate, k8sClient)
	if err != nil {
		return &v1beta1.CronJob{}, err
	}

	//如果日志采集目录不为空，则处理日志采集
	if nr.LogPath != "" {
		//日志采集
		var filebeatContainer v1.Container
		var index int

		for k, container := range cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers {
			if container.Name == "filebeat" {
				filebeatContainer = container
				index = k
				continue
			}
			if container.Name == name {
				var volumeMount v1.VolumeMount
				volumeIndex := -1
				for kk, v := range container.VolumeMounts {
					if v.Name == "app-logs" {
						volumeIndex = kk
						break
					}
				}
				volumeMount.Name = "app-logs"
				volumeMount.MountPath = nr.LogPath
				if volumeIndex == -1 {
					cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[k].VolumeMounts = append(cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[k].VolumeMounts, volumeMount)
				} else {
					cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[k].VolumeMounts[volumeIndex] = volumeMount
				}
			}
		}

		jobLogPaths := []string{nr.LogPath}

		if filebeatContainer.Name != "" {
			cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[index] = putContainer(filebeatContainer, name, jobLogPaths, confClient)
		} else {
			cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers = append(cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers, putContainer(filebeatContainer, name, jobLogPaths, confClient))
		}

		if cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[0].EnvFrom == nil || len(cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[0].EnvFrom) == 0 {
			var envFrom = new(v1.EnvFromSource)
			envFrom.ConfigMapRef = &v1.ConfigMapEnvSource{}
			envFrom.ConfigMapRef.LocalObjectReference.Name = name
			cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[0].EnvFrom = append(cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Containers[0].EnvFrom, *envFrom)
		}

		var volumes []v1.Volume
		volumes = append(volumes, v1.Volume{
			Name:         "app-logs",
			VolumeSource: v1.VolumeSource{EmptyDir: nil},
		})

		var configMapExists bool
		for _, v := range cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Volumes {
			if v.Name == name {
				configMapExists = true
				break
			}
			if v.Name == "app-logs" {
				volumes = []v1.Volume{}
			}
		}

		if !configMapExists {
			mode := int32(420)
			volumes = append(volumes, v1.Volume{
				Name: name,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: name,
						},
						DefaultMode: &mode,
					},
				},
			})
		}

		cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Volumes = append(cronjobInfo.Spec.JobTemplate.Spec.Template.Spec.Volumes, volumes...)

		//update

		cronjobInfo, err = cronjobUpdate(nr.Namespace, cronjobInfo, k8sClient)
		if err != nil {
			return &v1beta1.CronJob{}, ErrCronJobClientCreateFailed
		}
	}

	return
}

func cronjobUpdate(ns string, template interface{}, k8sClient kubernetes.K8sClient) (cron *v1beta1.CronJob, err error) {
	var cronjob *v1beta1.CronJob

	switch reflect.TypeOf(template).String() {
	case "string":
		s, ok := template.(string)
		if !ok {
			return
		}
		cronjob, err = convertToV1Beta(s)
		if err != nil {
			return
		}
	case "*v1beta1.CronJob":
		cronjob = template.(*v1beta1.CronJob)
	}

	cron, err = k8sClient.Do().BatchV1beta1().CronJobs(ns).Update(cronjob)
	if err != nil {
		return
	}
	cron.Kind = kind
	cron.APIVersion = version

	return
}

func putContainer(container v1.Container, name string, paths []string, confClient *config.Config) v1.Container {
	container.Name = "filebeat"
	dockerRepo := confClient.GetString("server", "docker_repo")
	container.Image = fmt.Sprintf(dockerRepo + "/filebeat:6.2.4")
	container.Args = []string{"-c", "/etc/filebeat/filebeat.yml", "-e"}

	resourceRequest := v1.ResourceList{}
	resourceLimit := v1.ResourceList{}
	resourceRequest[v1.ResourceMemory] = resource2.MustParse("128Mi")
	resourceLimit[v1.ResourceMemory] = resource2.MustParse("256Mi")

	container.Resources = v1.ResourceRequirements{
		Limits:   resourceLimit,
		Requests: resourceRequest,
	}

	var mounts []v1.VolumeMount
	for _, path := range paths {
		mounts = append(mounts, v1.VolumeMount{
			Name:      "app-logs",
			MountPath: path,
		})
		// todo 考虑要不要支持多个路径
		break
	}

	mounts = append(mounts, v1.VolumeMount{
		Name:      name,
		ReadOnly:  true,
		MountPath: "/etc/filebeat/filebeat.yml",
		SubPath:   "filebeat.yml",
	})

	container.VolumeMounts = mounts
	container.ImagePullPolicy = v1.PullIfNotPresent

	return container
}

func patch(ns string, tmp string, k8sClient kubernetes.K8sClient) (cron *v1beta1.CronJob, err error) {
	cron, err = convertToV1Beta(tmp)
	if err != nil {
		return
	}
	sv, err := json.Marshal(cron)
	cron, err = k8sClient.Do().BatchV1beta1().CronJobs(ns).Patch(cron.Name, k8sTypes.MergePatchType, sv)
	if err != nil {
		return
	}
	cron.Kind = kind
	cron.APIVersion = version
	return
}

// yaml -> json -> struct
func convertToV1Beta(tmp string) (cronjob *v1beta1.CronJob, err error) {
	var body interface{}

	if err = yaml.Unmarshal([]byte(tmp), &body); err != nil {
		return
	}

	body = helper.Convert(body)
	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &cronjob)
	if err != nil {
		return
	}
	return
}

func parseGitAddr(gitAddr string) (owner string, repo string, git string) {
	gitAddr = strings.Replace(gitAddr, ".git", "", -1)
	addr := strings.Split(gitAddr, ":")
	gits := strings.Split(addr[0], "@")
	names := strings.Split(addr[1], "/")
	owner = names[len(names)-2]
	repo = names[len(names)-1]

	return owner, repo, gits[1]
}

func NewService(logger log.Logger, config *config.Config,
	jenkins jenkins.Jenkins,
	k8sClient kubernetes.K8sClient,
	amqpClient amqpClient.AmqpClient,
	repository repository.Repository) Service {
	return &service{logger, config,
		jenkins,
		k8sClient,
		amqpClient,
		repository}
}
