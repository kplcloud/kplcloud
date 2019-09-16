package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kplcloud/kplcloud/src/config"
	"time"
)

var initDb *gorm.DB

func NewDb(cf *config.Config) (*gorm.DB, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci", cf.GetString("mysql", "mysql_user"), cf.GetString("mysql", "mysql_password"), cf.GetString("mysql", "mysql_host"), cf.GetString("mysql", "mysql_port"), cf.GetString("mysql", "mysql_database"))
	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)
	db.LogMode(cf.GetBool("server", "debug"))
	initDb = db
	return db, nil
}

func GetDb() *gorm.DB {
	return initDb
}
