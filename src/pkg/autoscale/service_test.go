/**
 * @Time : 2021/12/13 10:14 AM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */

package autoscale

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	mysqlclient "github.com/icowan/mysql-client"
	"github.com/kplcloud/kplcloud/src/kubernetes"
	"github.com/kplcloud/kplcloud/src/middleware"
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

func TestService_Sync(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.ContextKeyClusterId, 21)
	ctx = context.WithValue(ctx, middleware.ContextKeyClusterName, "c5")
	err := initSvc().Sync(ctx, 21, "kpaas")
	if err != nil {
		t.Error(err)
	}
}
