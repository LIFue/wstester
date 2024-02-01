package data

import (
	"path/filepath"
	"time"
	"wstester/internal/entity"
	"wstester/internal/schema"
	"wstester/pkg/dir"
	"wstester/pkg/log"

	"gorm.io/driver/mysql"
	//"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Data struct {
	DB *gorm.DB
}

func NewData(db *gorm.DB) (*Data, error) {
	if err := initTables(db); err != nil {
		return nil, err
	}
	return &Data{DB: db}, nil
}

func initTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Platform{}, &entity.SshHost{}, &entity.User{}, &entity.Message); err != nil {
		return err
	}
	return nil
}

func NewDB(debug bool, dbConf *Database) (db *gorm.DB, err error) {
	if dbConf.Driver == "" {
		dbConf.Driver = string(schema.MYSQL)
	} else if dbConf.Driver == string(schema.SQLITE) {
		//dbConf.Driver = "sqlite"
		dbFileDir := filepath.Dir(dbConf.Connection)
		if err := dir.CreateDirIfNotExist(dbFileDir); err != nil {
			log.Errorf("create database dir failed: %s", err)
		}
		dbConf.MaxOpenConn = 1
	}

	switch dbConf.Driver {
	case string(schema.MYSQL):
		db, err = gorm.Open(mysql.Open(dbConf.Connection), &gorm.Config{})
		//case string(schema.SQLITE):
		//	db, err = gorm.Open(sqlite.Open(dbConf.Connection), &gorm.Config{})
	}

	if err != nil {
		return
	}

	if debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConn)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConn)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Duration(dbConf.ConnMaxLifeTime))

	return
}
