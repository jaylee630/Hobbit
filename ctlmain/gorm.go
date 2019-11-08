package ctlmain

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	gormzap "github.com/wantedly/gorm-zap"

	"golang-admin-basic-master/utils/log"
)

var (
	MySQLTablePrefix = "tbl_"
)

func NewGORM(l *log.Logger, logLevel string, dialect, ds string) (*gorm.DB, error) {

	db, err := gorm.Open(dialect, ds)
	if err != nil {
		return nil, err
	}

	if logLevel == "debug" {
		db.LogMode(true)
	} else {
		db.SetLogger(gormzap.New(l.L()))
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return MySQLTablePrefix + defaultTableName
	}

	db.DB().SetMaxIdleConns(50)
	db.DB().SetConnMaxLifetime(time.Hour)

	return db, nil
}
