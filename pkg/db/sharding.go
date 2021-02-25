package db

import (
	"log"
	"os"
	. "server/pkg/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// map[ds_0:0xc000366240 ds_1:0xc0003a2570 ds_2:0xc0003669f0 ds_3:0xc000367500]
var ShardingDBs map[string]*gorm.DB

func init() {
	ShardingDBs = map[string]*gorm.DB{}
	for _, ds := range ShardingConfig.ShardingDatasources {
		ShardingDBs[ds.Name] = initShardingDB(GlobalConfig.Env, ds.Write, ds.Read)
	}
}

func initShardingDB(env string, write string, read string) *gorm.DB {
	// https://gorm.io/zh_CN/docs/dbresolver.html
	// 建立读写分离连接池
	dsn_master := write
	dsn_slave := read

	// log level
	var level logger.LogLevel
	if env == "local" || env == "dev" || env == "test" {
		// level = logger.Info
		level = logger.Silent
	} else {
		level = logger.Silent
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second / 4, // Slow SQL threshold
			LogLevel:      level,           // Log level
		},
	)
	db, _ := gorm.Open(mysql.Open(dsn_master), &gorm.Config{Logger: newLogger})
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsn_master)},
		Replicas: []gorm.Dialector{mysql.Open(dsn_slave)},
		Policy:   dbresolver.RandomPolicy{},
	}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(100).
		SetMaxOpenConns(200))
	return db
}
