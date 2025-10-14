package main

import (
	"blog-server/core"
	"blog-server/global"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	global.AppConfig = core.InitConfig()
	global.DB = core.InitDB()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(fmt.Sprintf(":%v", global.AppConfig.App.Port))
}
