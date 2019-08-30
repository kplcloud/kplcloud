/**
 * @Time : 2019/6/27 10:15 AM
 * @Author : yuntinghu1003@gmail.com
 * @File : event
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository/types"
)

type EventRepository interface {
	FindByIds(ids []int) (list []*types.Event, err error)
	FindAllEvents() (list []*types.Event, err error)
	FindByKind(kind EventsKind) (res *types.Event, err error)
}

type event struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &event{db: db}
}

type EventsKind string

const (
	BuildEvent     EventsKind = "Build"
	ApplyEvent     EventsKind = "Apply"
	AuditEvent     EventsKind = "Audit"
	DeleteEvent    EventsKind = "Delete"
	MemberEvent    EventsKind = "Member"
	RollbackEvent  EventsKind = "Rollback"
	LoggingEvent   EventsKind = "Logging"
	RebootEvent    EventsKind = "Reboot"
	CommandEvent   EventsKind = "Command"
	StorageEvent   EventsKind = "Storage"
	IngressGateway EventsKind = "Gateway"
	Expansion      EventsKind = "Expansion"      //扩容
	Extend         EventsKind = "Extend"         //伸缩
	SwitchModel    EventsKind = "SwitchModel"    //调整服务模式
	ReadinessProbe EventsKind = "ReadinessProbe" //修改探针
	TestEvent      EventsKind = "Test"
	VolumeConfig   EventsKind = "VolumeConfig"
	VolumeHosts    EventsKind = "VolumeHosts"
)

func (c EventsKind) String() string {
	return string(c)
}

func (c *event) FindByIds(ids []int) (list []*types.Event, err error) {
	err = c.db.Where("id in (?)", ids).Find(&list).Error
	return
}

func (c *event) FindAllEvents() (list []*types.Event, err error) {
	err = c.db.Find(&list).Error
	return
}

func (c *event) FindByKind(kind EventsKind) (res *types.Event, err error) {
	var event types.Event
	err = c.db.Where("name = ?", kind.String()).Preload("WebHook").Find(&event).Error
	return &event, err
}
