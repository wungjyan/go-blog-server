package bootstrap

import (
	"blog-server/core"
	"blog-server/global"
)

func Initialize() {
	global.AppConfig = core.InitConfig()
	global.Logger = core.InitLogger()
	global.DB = core.InitDB()

	global.Logger.Info("🚀 应用初始化完成")
}
