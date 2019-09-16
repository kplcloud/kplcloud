/**
 * @Time : 2019-08-22 10:09
 * @Author : solacowa@gmail.com
 * @File : root
 * @Software: GoLand
 */

package server

import (
	"errors"
	"github.com/go-kit/kit/log/level"
	"github.com/go-sql-driver/mysql"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

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
