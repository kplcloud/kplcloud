/**
 * @Time : 2019-08-22 10:09
 * @Author : solacowa@gmail.com
 * @File : root
 * @Software: GoLand
 */

package server

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/go-sql-driver/mysql"
	"github.com/icowan/config"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/jinzhu/gorm"
	"github.com/kplcloud/kplcloud/src/logging"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"github.com/spf13/cobra"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	generateCmd = &cobra.Command{
		Use:               "generate command <args> [flags]",
		Short:             "生成命令",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `## 生成命令
可用的配置类型：
[table, init-data]

kplcloud generate -h
`,
	}

	genTableCmd = &cobra.Command{
		Use:               `table <args> [flags]`,
		Short:             "生成数据库表",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kplcloud generate table all
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = rds.Close()
				}
			}()

			if len(args) > 0 && args[0] == "all" {
				return nil
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runPre()
		},
	}

	genInitDataCmd = &cobra.Command{
		Use:               `init-data <args> [flags]`,
		Short:             "生成数据",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kplcloud generate init-data
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = level.Debug(logger).Log("redis", "close", "err", rds.Close())
				}
			}()

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runPre()
		},
	}

	genAdminUserCmd = &cobra.Command{
		Use:               `admin-user <args> [flags]`,
		Short:             "生成管理员用户",
		SilenceErrors:     false,
		DisableAutoGenTag: false,
		Example: `
kplcloud generate admin-user
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 关闭资源连接
			defer func() {
				_ = level.Debug(logger).Log("db", "close", "err", db.Close())
				if rds != nil {
					_ = level.Debug(logger).Log("redis", "close", "err", rds.Close())
				}
			}()

			if strings.EqualFold(adminEmail, "") {
				adminEmail = args[0]
			}
			if strings.EqualFold(adminPassword, "") {
				adminPassword = args[1]
			}

			if validateFormat(adminEmail) != nil {
				return errors.New("邮箱不正确")
			}

			_, err = store.Member().Find(adminEmail)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("用户已存在")
			}

			var nss []types.Namespace
			if ns, err := store.Namespace().FindAll(); err == nil {
				for _, v := range ns {
					nss = append(nss, *v)
				}
			}

			role, err := store.Role().FindById(1)
			if err != nil {
				return err
			}
			var roles []types.Role
			roles = append(roles, *role)
			member := &types.Member{
				Email:      adminEmail,
				Username:   adminEmail,
				Password:   null.StringFrom(encode.EncodePassword(adminPassword, cf.GetString("server", "app_key"))),
				Roles:      roles,
				Namespaces: nss,
			}
			return store.Member().CreateMember(member)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runPre()
		},
	}
)

func runPre() error {
	cf, err = config.NewConfig(configPath)
	if err != nil {
		return err
	}
	_ = cf.SetValue("server", "kube_config", kubeConfig)

	logger = logging.SetLogging(logger, cf.GetString("server", "log_path"), cf.GetString("server", "log_level"))

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		cf.GetString(config.SectionMysql, "user"),
		cf.GetString(config.SectionMysql, "password"),
		cf.GetString(config.SectionMysql, "host"),
		cf.GetString(config.SectionMysql, "port"),
		cf.GetString(config.SectionMysql, "database"))

	// 连接数据库
	db, err = mysqlclient.NewMysql(dbUrl, cf.GetBool(config.SectionServer, "debug"))
	if err != nil {
		_ = level.Error(logger).Log("db", "connect", "err", err)
		return err
	}

	store = repository.NewRepository(db)

	return nil
}

var (
	ErrBadFormat = errors.New("invalid format")
	//ErrUnresolvableHost = errors.New("unresolvable host")

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrBadFormat
	}
	return nil
}

func importToDb(appKey string) error {

	if validateFormat(adminEmail) != nil {
		return errors.New("邮箱不正确")
	}

	if _, err := store.Member().Find(adminEmail); err != nil {
		switch err.(type) {
		case *mysql.MySQLError:
			e := err.(*mysql.MySQLError)
			if e.Number == 1146 {
				goto CREATE
			}
		}
		return nil
	}

	return nil

CREATE:
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Member{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Build{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Casbin{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.ConfigMap{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.ConfigData{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Consul{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Cronjob{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Dockerfile{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Event{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.EventHistory{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Groups{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.GroupsCronjobs{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.GroupsMemberss{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Namespace{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.NamespacesMembers{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.NoticeMember{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.NoticeReceive{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Notices{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Permission{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.PersistentVolumeClaim{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Project{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.ProjectJenkins{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.ProjectTemplate{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Role{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.StorageClass{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Template{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.Webhook{}).Error)
	_ = level.Debug(logger).Log("create table", db.CreateTable(types.WechatUser{}).Error)
	_ = level.Debug(logger).Log("create table", db.Exec("CREATE TABLE `casbin_rule` (`p_type` varchar(100) DEFAULT NULL,`v0` varchar(100) DEFAULT NULL,`v1` varchar(100) DEFAULT NULL, `v2` varchar(100) DEFAULT NULL, `v3` varchar(100) DEFAULT NULL, `v4` varchar(100) DEFAULT NULL, `v5` varchar(100) DEFAULT NULL) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").Error)

	path, err := filepath.Abs(sqlPath)
	if err != nil {
		_ = level.Warn(logger).Log("filepath", "Abs", "err", err.Error())
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		_ = level.Warn(logger).Log("ioutil", "ReadFile", "err", err.Error())
		return err
	}

	ds := strings.Split(string(data), ");\n")
	for _, v := range ds {
		if !strings.Contains(v, ");") {
			v += ");"
		}
		db.Exec(v)
	}
	u := strings.Split(adminEmail, "@")

	var nss []types.Namespace
	if ns, err := store.Namespace().FindAll(); err == nil {
		for _, v := range ns {
			nss = append(nss, *v)
		}
	}

	if role, err := store.Role().FindById(1); err == nil && role != nil {
		var roles []types.Role
		roles = append(roles, *role)
		member := &types.Member{
			Email:      adminEmail,
			Username:   u[0],
			Password:   null.StringFrom(encode.EncodePassword(adminPassword, appKey)),
			Roles:      roles,
			Namespaces: nss,
		}
		return store.Member().CreateMember(member)
	}

	return err
}
