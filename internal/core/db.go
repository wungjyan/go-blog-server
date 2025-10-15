package core

import (
	"blog-server/global"
	"blog-server/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	logger := global.Logger
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
		logger.Fatalf("❌ 数据库连接失败: %v", err)
	}

	setPool(db)
	migrate(db)

	logger.Info("✅ 数据库实例初始化成功")

	return db
}

func setPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		global.Logger.Errorf("❌ 获取底层数据库连接失败: %v", err)
		return
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.Logger.Infof("✅ 数据库连接池配置完成 (MaxOpen=%d, MaxIdle=%d, Lifetime=%v)",
		100, 10, time.Hour)
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		global.Logger.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	global.Logger.Info("✅ 数据库自动迁移完成")
}
