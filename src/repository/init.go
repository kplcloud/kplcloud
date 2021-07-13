/**
 * @Time : 2019-08-09 12:16
 * @Author : solacowa@gmail.com
 * @File : init
 * @Software: GoLand
 */

package repository

import "github.com/jinzhu/gorm"

type Repository interface {
	Build() BuildRepository
	Notice() NoticesRepository
	Proclaim() NoticesRepository
	NoticeReceive() NoticeReceiveRepository
	Namespace() NamespaceRepository
	Member() MemberRepository
	Template() TemplateRepository
	Groups() GroupsRepository
	StorageClass() StorageClassRepository
	Pvc() PvcRepository
	Project() ProjectRepository
	ProjectTemplate() ProjectTemplateRepository
	Webhook() WebhookRepository
	Event() EventRepository
	EventHistory() EventHistoryRepository
	CronJob() CronjobRepository
	ConfigMap() ConfigMapRepository
	ConfigData() ConfigDataRepository
	WechatUser() WechatUserRepository
	Permission() PermissionRepository
	Role() RoleRepository
	ProjectJenkins() ProjectJenkinsRepository
	Consul() ConsulRepository
	Dockerfile() DockerfileRepository
	NoticeMember() NoticeMemberRepository
	ConfigEnv() ConfigEnvRepository
}

type store struct {
	db *gorm.DB

	build           BuildRepository
	notice          NoticesRepository
	proclaim        NoticesRepository
	noticeReceive   NoticeReceiveRepository
	namespace       NamespaceRepository
	member          MemberRepository
	template        TemplateRepository
	groups          GroupsRepository
	storageClass    StorageClassRepository
	pvc             PvcRepository
	project         ProjectRepository
	projectTemplate ProjectTemplateRepository
	webhook         WebhookRepository
	event           EventRepository
	eventHistory    EventHistoryRepository
	cronJob         CronjobRepository
	configMap       ConfigMapRepository
	configData      ConfigDataRepository
	wechatUser      WechatUserRepository
	permission      PermissionRepository
	role            RoleRepository
	projectJenkins  ProjectJenkinsRepository
	consul          ConsulRepository
	dockerfile      DockerfileRepository
	noticeMember    NoticeMemberRepository
	configEnv       ConfigEnvRepository
}

func NewRepository(db *gorm.DB) Repository {
	return &store{
		db:              db,
		build:           NewBuildRepository(db),
		notice:          NewNoticesRepository(db),
		proclaim:        NewNoticesRepository(db),
		noticeReceive:   NewNoticeReceiveRepository(db),
		namespace:       NewNamespaceRepository(db),
		member:          NewMemberRepository(db),
		template:        NewTemplateRepository(db),
		groups:          NewGroupsRepository(db),
		storageClass:    NewStorageClassRepository(db),
		pvc:             NewPvcRepository(db),
		project:         NewProjectRepository(db),
		projectTemplate: NewProjectTemplateRepository(db),
		webhook:         NewWebhookRepository(db),
		event:           NewEventRepository(db),
		eventHistory:    NewEventHistoryRepository(db),
		cronJob:         NewCronjobRepository(db),
		configMap:       NewConfigMapRepository(db),
		configData:      NewConfigDataRepository(db),
		wechatUser:      NewWechatUserRepository(db),
		permission:      NewPermissionRepository(db),
		role:            NewRoleRepository(db),
		projectJenkins:  NewProjectJenkins(db),
		consul:          NewConsulReporitory(db),
		dockerfile:      NewDockerfileRepository(db),
		noticeMember:    NewNoticeMemberRepository(db),
		configEnv:       NewConfigEnvRepository(db),
	}
}

func (c *store) Build() BuildRepository {
	return c.build
}

func (c *store) Notice() NoticesRepository {
	return c.notice
}

func (c *store) Proclaim() NoticesRepository {
	return c.proclaim
}

func (c *store) NoticeReceive() NoticeReceiveRepository {
	return c.noticeReceive
}

func (c *store) Namespace() NamespaceRepository {
	return c.namespace
}

func (c *store) Member() MemberRepository {
	return c.member
}

func (c *store) Template() TemplateRepository {
	return c.template
}

func (c *store) Groups() GroupsRepository {
	return c.groups
}

func (c *store) StorageClass() StorageClassRepository {
	return c.storageClass
}

func (c *store) Pvc() PvcRepository {
	return c.pvc
}

func (c *store) Project() ProjectRepository {
	return c.project
}

func (c *store) ProjectTemplate() ProjectTemplateRepository {
	return c.projectTemplate
}

func (c *store) Webhook() WebhookRepository {
	return c.webhook
}

func (c *store) Event() EventRepository {
	return c.event
}

func (c *store) CronJob() CronjobRepository {
	return c.cronJob
}

func (c *store) EventHistory() EventHistoryRepository {
	return c.eventHistory
}

func (c *store) ConfigMap() ConfigMapRepository {
	return c.configMap
}

func (c *store) ConfigData() ConfigDataRepository {
	return c.configData
}

func (c *store) WechatUser() WechatUserRepository {
	return c.wechatUser
}

func (c *store) Permission() PermissionRepository {
	return c.permission
}

func (c *store) Role() RoleRepository {
	return c.role
}

func (c *store) ProjectJenkins() ProjectJenkinsRepository {
	return c.projectJenkins
}

func (c *store) Consul() ConsulRepository {
	return c.consul
}

func (c *store) Dockerfile() DockerfileRepository {
	return c.dockerfile
}

func (c *store) NoticeMember() NoticeMemberRepository {
	return c.noticeMember
}

func (c *store) ConfigEnv() ConfigEnvRepository {
	return c.configEnv
}
