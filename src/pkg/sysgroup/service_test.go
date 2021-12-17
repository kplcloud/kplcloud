/**
 * @Time : 2021/12/17 1:51 PM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */

package sysgroup

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/kplcloud/kplcloud/src/repository"
	"os"
	"testing"
)

func initSvc() Service {
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		portStr  = os.Getenv("DB_PORT")
		database = os.Getenv("DB_DATABASE")
	)
	logger := log.NewNopLogger()
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
		user,
		password,
		host,
		portStr,
		database)
	db, err := mysqlclient.NewMysql(dbUrl, true)
	if err != nil {
		panic(err)
	}
	store := repository.New(db, logger, "traceId", nil, nil, nil)
	return New(logger, "traceId", store)
}

func TestService_List(t *testing.T) {
	svc := initSvc()
	list, total, err := svc.List(context.Background(), 1, []int64{1}, "kpaas", "", 1, 10)
	if err != nil {
		t.Error()
		return
	}
	t.Log(total)
	for _, v := range list {
		t.Log(v)
	}
}
