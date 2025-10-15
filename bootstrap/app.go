package bootstrap

import (
	"blog-server/global"
	"blog-server/internal/core"
)

func Initialize() {
	global.AppConfig = core.InitConfig()
	global.Logger = core.InitLogger()
	global.DB = core.InitDB()

	global.Logger.Info("🚀 应用初始化完成")
}
