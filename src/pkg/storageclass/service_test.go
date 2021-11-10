/**
 * @Time : 2021/11/9 4:54 PM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */

package storageclass

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/kplcloud/kplcloud/src/kubernetes"
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
		_ = level.Error(logger).Log("db", "connect", "err", err)
		panic(err)
	}
	store := repository.New(db, logger, "traceId", nil, nil, nil)
	k8sClient, err := kubernetes.NewClient(store)
	if err != nil {
		panic(err)
	}
	return New(logger, "traceId", store, k8sClient)
}

func TestService_List(t *testing.T) {
	list, total, err := initSvc().List(context.Background(), 13, 1, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("total ===> ", total)
	for _, v := range list {
		t.Log(v.Name)
	}
}
