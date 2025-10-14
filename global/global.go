package global

import (
	"blog-server/config"

	"gorm.io/gorm"
)

var (
	AppConfig config.Config
	DB        *gorm.DB
)
