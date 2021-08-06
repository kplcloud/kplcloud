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
	redisclient "github.com/icowan/redis-client"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type Middleware func(Service) Service

type Service interface {
	// 1. 初始化数据库,重新加载配置文件,连接数据库
	InitDb(ctx context.Context, drive, host string, port int, username, password, database string) (err error)
	// 2. platform设置
	InitPlatform(ctx context.Context, appName, adminName, adminPassword, appKey, domain, domainSuffix, logPath, logLevel, uploadPath string, debug bool) (err error)
	// 8. logo 设置
	InitLogo(ctx context.Context, f *multipart.FileHeader) (err error)
	// 9. 跨域配置
	InitCors(ctx context.Context, allow bool, origin, methods, headers string) (err error)
	// 3. 初始化Redis
	InitRedis(ctx context.Context, hosts, auth string, db int, prefix string) (err error)
	// 4. 初始化Jenkins构建机器
	InitJenkins(ctx context.Context) (err error)
	// 5. 初始化MQ 可以考虑直接用redis
	InitMq(ctx context.Context) (err error)
	// 6. 初始化镜像仓库
	InitRepo(ctx context.Context) (err error)
	// 7. 初始化k8s集群
	InitK8sCluster(ctx context.Context) (err error)
	// 10. 配置写入文件
	StoreToConfig(ctx context.Context) (err error)

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
	rds        *redisclient.RedisClient
}

func (s *service) InitPlatform(ctx context.Context, appName, adminName, adminPassword, appKey, domain, domainSuffix, logPath, logLevel, uploadPath string, debug bool) (err error) {
	if strings.EqualFold(appKey, "") {
		appKey = string(util.Krand(12, util.KC_RAND_KIND_ALL))
	}
	_ = s.repository.SysSetting().Add(ctx, "server", "name", appName, "平台名称")
	_ = s.repository.SysSetting().Add(ctx, "server", "key", appKey, "平台Key")
	_ = s.repository.SysSetting().Add(ctx, "server", "domain", domain, "平台域名")
	_ = s.repository.SysSetting().Add(ctx, "server", "domain.suffix", domainSuffix, "平台域名后缀,生成对名域名用")
	_ = s.repository.SysSetting().Add(ctx, "server", "log.path", logPath, "平台日志路径,不填的话打在控制台")
	_ = s.repository.SysSetting().Add(ctx, "server", "log.level", logLevel, "平台输出的日志级别,支持五个级别 all,error,warn,info,debug")
	_ = s.repository.SysSetting().Add(ctx, "server", "upload.path", uploadPath, "平台文件上传路径")
	_ = s.repository.SysSetting().Add(ctx, "server", "web.path", "/usr/local/kplcloud/web/v2", "Web路径")
	_ = s.repository.SysSetting().Add(ctx, "server", "debug", strconv.FormatBool(debug), "是否输出Debug日志")

	// 保存管理员账号到SysUser
	// 查询角色
	roles, err := s.repository.SysRole().FindByIds(ctx, []int64{1})
	if err != nil {
		err = errors.Wrap(err, "repository.SysRole.FindByIds")
		return
	}
	err = s.repository.SysUser().Save(ctx, &types.SysUser{
		Username:  adminName,
		LoginName: adminName,
		Email:     adminName,
		Password:  util.EncodePassword(adminPassword, appKey),
		SysRoles:  roles,
	})

	return
}

func (s *service) InitLogo(ctx context.Context, f *multipart.FileHeader) (err error) {
	st, err := s.repository.SysSetting().Find(ctx, "server", "web.path")
	if err != nil {
		_ = level.Error(s.logger).Log("repository.SysSetting", "Find", "err", err.Error())
		err = encode.ErrInstallUploadPath.Error()
		return
	}

	file, err := f.Open()
	if err != nil {
		_ = level.Error(s.logger).Log("f", "Open", "err", err.Error())
		err = encode.ErrInstallUpload.Wrap(err)
		return
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		_ = level.Error(s.logger).Log("ioutil", "ReadAll", "err", err.Error())
		err = encode.ErrInstallUpload.Wrap(err)
		return
	}

	logoPath := fmt.Sprintf("%s/images/logo%s", st.Value, ".png" /*path.Ext(f.Filename)*/)

	if err = ioutil.WriteFile(logoPath, b, 0666); err != nil {
		_ = level.Error(s.logger).Log("ioutil", "WriteFile", "err", err.Error())
		err = encode.ErrInstallUpload.Wrap(err)
		return
	}

	_ = s.repository.SysSetting().Add(ctx, "server", "logo", logoPath, "平台logo")
	return
}

func (s *service) InitCors(ctx context.Context, allow bool, origin, methods, headers string) (err error) {
	_ = s.repository.SysSetting().Add(ctx, "cors", "allow", strconv.FormatBool(allow), "是否允许跨域")
	_ = s.repository.SysSetting().Add(ctx, "cors", "origin", origin, "跨域来源 '*' 表示所有来源")
	_ = s.repository.SysSetting().Add(ctx, "cors", "methods", methods, "允许跨域Method")
	_ = s.repository.SysSetting().Add(ctx, "cors", "headers", headers, "允许跨域Headers")

	return
}

func (s *service) StoreToConfig(ctx context.Context) (err error) {
	res, err := s.repository.SysSetting().FindAll(ctx)
	if err != nil {
		return
	}

	for _, v := range res {
		s.cfg.ConfigFile.SetKeyComments(v.Section, v.Key, v.Description)
		s.cfg.ConfigFile.SetValue(v.Section, v.Key, v.Value)
	}

	return s.configReload(ctx)
}

func (s *service) InitDb(ctx context.Context, drive, host string, port int, username, password, database string) (err error) {
	err = s.database(ctx, drive, host, port, username, password, database)
	if err != nil {
		err = encode.ErrInstallDbConnect.Wrap(errors.Wrap(err, "database.init"))
		return
	}
	if !strings.EqualFold(drive, "mysql") {
		err = encode.ErrInstallDbDrive.Error()
		return
	}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		username, password, host, port, database)

	// 连接数据库
	db, err := mysqlclient.NewMysql(dbUrl, s.cfg.GetBool(config.SectionServer, "debug"))
	if err != nil {
		_ = level.Error(s.logger).Log("db", "connect", "err", err)
		err = encode.ErrInstallDbConnect.Wrap(err)
		return err
	}

	s.db = db
	s.repository = repository.New(db, s.logger, "traceId", nil, nil)

	// 初始化数据
	return s.initData(ctx)
}

func (s *service) initData(ctx context.Context) (err error) {
	_ = s.logger.Log("create", "table", "SysRole", s.db.AutoMigrate(types.SysRole{}).Error)
	_ = s.logger.Log("create", "table", "SysUser", s.db.AutoMigrate(types.SysUser{}).Error)
	_ = s.logger.Log("create", "table", "SysPermission", s.db.AutoMigrate(types.SysPermission{}).Error)
	_ = s.logger.Log("create", "table", "SysSetting", s.db.AutoMigrate(types.SysSetting{}).Error)
	_ = s.logger.Log("create", "table", "SysNamespace", s.db.AutoMigrate(types.SysNamespace{}).Error)

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
	_ = s.logger.Log("add", "data", "publicKey", s.repository.SysSetting().Add(ctx, "server", "rsa.public.key", publicKey, "公钥"))
	_ = s.logger.Log("add", "data", "privateKey", s.repository.SysSetting().Add(ctx, "server", "rsa.private.key", privateKey, "私钥"))

	// TODO: 写入角色数据
	// TODO: 写入权限数据
	// TODO: 写入其他系统基础设置

	return
}

func (s *service) InitRedis(ctx context.Context, hosts, auth string, db int, prefix string) (err error) {
	_ = s.repository.SysSetting().Add(ctx, "redis", "hosts", hosts, "单点 hosts: 127.0.0.1:6379;集群 hosts: 127.0.0.1:6380 用\",\"隔开")
	_ = s.repository.SysSetting().Add(ctx, "redis", "auth", auth, "密码")
	_ = s.repository.SysSetting().Add(ctx, "redis", "db", strconv.Itoa(db), "db")
	_ = s.repository.SysSetting().Add(ctx, "redis", "prefix", prefix, "前缀")

	rds, err := redisclient.NewRedisClient(hosts, auth, prefix, db)
	if err != nil {
		_ = level.Error(s.logger).Log("redisclient", "NewRedisClient", "err", err.Error())
		return
	}
	if err = rds.Set(ctx, "hello", "world", time.Second*5); err != nil {
		_ = level.Error(s.logger).Log("rds", "Set", "err", err.Error())
		return
	}
	s.rds = &rds

	return nil
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

func New(logger log.Logger, cfg *config.Config, cfgPath string, store repository.Repository) Service {
	return &service{cfg: cfg, cfgPath: cfgPath, logger: logger, repository: store}
}
