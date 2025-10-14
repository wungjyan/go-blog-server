package core

import (
	"blog-server/global"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (db *gorm.DB) {
	dbConf := global.AppConfig.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.Username,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("打开数据库失败：", err)
		panic(err)
	}

	setPool(db)

	fmt.Println("数据库实例初始化成功")

	return
}

func setPool(db *gorm.DB) {
	sqlDB, err := db.DB()

	if err != nil {
		fmt.Println("设置连接池失败：", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
