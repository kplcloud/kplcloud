/**
 * @Time : 7/20/21 6:23 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package install

import (
	"context"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Middleware func(Service) Service

type Service interface {
	// 0
	Init(ctx context.Context, appName string) (err error)
	// 1. 初始化数据库,重新加载配置文件,连接数据库
	InitDb(ctx context.Context, drive, host string, port int, username, password, database string) (err error)
	// 2. 初始化Redis
	InitRedis(ctx context.Context, hosts, auth string, db int) (err error)
	// 初始化Jenkins构建机器
	InitJenkins(ctx context.Context) (err error)
	// 初始化MQ 可以考虑直接用redis
	InitMq(ctx context.Context) (err error)
	// 初始化镜像仓库
	InitRepo(ctx context.Context) (err error)
	// 初始化k8s集群
	InitK8sCluster(ctx context.Context) (err error)

	configReload(ctx context.Context) (err error)
	setValue(ctx context.Context, section, key, value string) bool
	database(ctx context.Context, drive, host string, port int, username, password, database string) (err error)
}

type service struct {
	logger     log.Logger
	cfg        *config.Config
	cfgPath    string
	repository repository.Repository
	db         *gorm.DB
}

func (s *service) Init(ctx context.Context, appName string) (err error) {
	_ = s.setValue(ctx, "server", "app.name", appName)
	_ = s.setValue(ctx, "server", "app.debug", "true")

	return
}

func (s *service) InitDb(ctx context.Context, drive, host string, port int, username, password, database string) (err error) {
	err = s.database(ctx, drive, host, port, username, password, database)
	if err != nil {
		err = errors.Wrap(err, "database.init")
		return
	}
	if !strings.EqualFold(drive, "mysql") {
		err = errors.New("暂不支持其他数据库.")
		return
	}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		username, password, host, port, database)

	// 连接数据库
	db, err := mysqlclient.NewMysql(dbUrl, s.cfg.GetBool(config.SectionServer, "app.debug"))
	if err != nil {
		_ = level.Error(s.logger).Log("db", "connect", "err", err)
		return err
	}

	s.db = db
	s.repository = repository.New(db, s.logger, "traceId", nil, nil)

	// 初始化数据
	return s.initData(ctx)
}

func (s *service) initData(ctx context.Context) (err error) {
	_ = s.logger.Log("create", "table", "SysRole", s.db.CreateTable(types.SysRole{}).Error)
	_ = s.logger.Log("create", "table", "SysUser", s.db.CreateTable(types.SysUser{}).Error)
	_ = s.logger.Log("create", "table", "SysPermission", s.db.CreateTable(types.SysPermission{}).Error)
	_ = s.logger.Log("create", "table", "SysSetting", s.db.CreateTable(types.SysSetting{}).Error)
	_ = s.logger.Log("create", "table", "SysNamespace", s.db.CreateTable(types.SysNamespace{}).Error)

	// 初始化数据
	authRsaPublicKey, authRsaPrivateKey, err := util.GenRsaKey()
	if err != nil {
		_ = level.Error(s.logger).Log("util", "GenRsaKey", "err", err.Error())
		return
	}

	// 生成公私钥
	publicKey := strings.TrimSpace(string(authRsaPublicKey))
	privateKey := strings.TrimSpace(string(authRsaPrivateKey))
	publicKey = strings.Trim(publicKey, "\n")
	privateKey = strings.Trim(privateKey, "\n")
	_ = s.logger.Log("add", "data", "publicKey", s.repository.SysSetting().Add(ctx, "AUTH_RSA_PUBLIC_KEY", publicKey, "公钥"))
	_ = s.logger.Log("add", "data", "privateKey", s.repository.SysSetting().Add(ctx, "AUTH_RSA_PRIVATE_KEY", privateKey, "私钥"))

	// TODO: 写入角色数据
	// TODO: 写入权限数据
	// TODO: 写入其他系统基础设置

	return
}

func (s *service) InitRedis(ctx context.Context, hosts, auth string, db int) (err error) {
	panic("implement me")
}

func (s *service) InitJenkins(ctx context.Context) (err error) {
	panic("implement me")
}

func (s *service) InitMq(ctx context.Context) (err error) {
	panic("implement me")
}

func (s *service) InitRepo(ctx context.Context) (err error) {
	panic("implement me")
}

func (s *service) InitK8sCluster(ctx context.Context) (err error) {
	panic("implement me")
}

func (s *service) database(ctx context.Context, drive, host string, port int, username, password, database string) (err error) {
	_ = s.setValue(ctx, "database", "drive", drive)
	_ = s.setValue(ctx, "database", "host", host)
	_ = s.setValue(ctx, "database", "port", strconv.Itoa(port))
	_ = s.setValue(ctx, "database", "username", username)
	_ = s.setValue(ctx, "database", "password", password)
	_ = s.setValue(ctx, "database", "database", database)
	return s.configReload(ctx)
}

func (s *service) configReload(ctx context.Context) (err error) {
	if err = goconfig.SaveConfigFile(s.cfg.ConfigFile, s.cfgPath); err != nil {
		return err
	}
	return s.cfg.ConfigFile.Reload()
}

func (s *service) setValue(ctx context.Context, section, key, value string) bool {
	return s.cfg.ConfigFile.SetValue(section, key, value)
}

func New(logger log.Logger, cfg *config.Config, cfgPath string) Service {
	return &service{cfg: cfg, cfgPath: cfgPath, logger: logger}
}
