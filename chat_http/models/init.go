package models

import (
	"conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(mysql.Open(conf.MysqlConnectString()), &gorm.Config{
		PrepareStmt:            true, //禁止预编译缓存命令
		SkipDefaultTransaction: true, //禁止默认事务
		//Logger:                 logger.Default.LogMode(logger.Info), //打印sql语句
	})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&User{}, &UserInfo{})
	if err != nil {
		panic(err)
	}
}
