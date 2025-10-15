package global

import (
	"blog-server/config"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	AppConfig config.Config
	DB        *gorm.DB
	Logger    *logrus.Logger
)
