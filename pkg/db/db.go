package db

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/latonaio/authenticator/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	Con  *gorm.DB
	Info Info
}

type Info struct {
	DBName    string
	TableName string
}

var ConPool = &DB{}

func NewDBConPool(ctx context.Context, config configs.Configs) error {
	cfgs := config.Get()
	ConPool.Info = Info{
		DBName:    cfgs.Database.Name,
		TableName: cfgs.Database.TableName,
	}
	dsn := fmt.Sprintf(XXXXXXXXX/XXXXXXXXXXXparseTime=True&loc=Local",
		cfgs.Database.UserName,
		cfgs.Database.UserPassword,
		cfgs.Database.HostName,
		cfgs.Database.Port,
		cfgs.Database.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	ConPool.Con = db.WithContext(ctx)
	mysqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = mysqlDB.Ping()
	if err != nil {
		panic(err)
	}
	mysqlDB.SetConnMaxIdleTime(24 * time.Hour)
	mysqlDB.SetMaxOpenConns(cfgs.Database.MaxOpenCon)
	mysqlDB.SetMaxIdleConns(cfgs.Database.MaxIdleCon)
	return nil
}
