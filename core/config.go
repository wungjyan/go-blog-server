package core

import (
	"blog-server/config"
	"blog-server/core/internal"
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitConfig() (conf config.Config) {
	configPath := getConfigPath()

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 启用环境变量
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件发生变化: %s", e.Name)
		if err := v.Unmarshal(&conf); err != nil {
			fmt.Printf("重新解析配置文件失败: %v", err)
			return
		}
		fmt.Println("配置文件热更新成功")
	})

	// 解析配置到结构体
	if err := v.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", err))
	}

	fmt.Printf("配置文件加载成功: %s \n", configPath)

	return
}

func getConfigPath() (configPath string) {
	flag.StringVar(&configPath, "c", "", "choose config file.")
	flag.Parse()

	// 命令行传递
	if configPath != "" {
		fmt.Printf("使用命令行 '-c' 参数传递的配置文件: %s\n", configPath)
		return
	}

	// 判断环境变量
	if env := os.Getenv(internal.ConfigEnv); env != "" {
		configPath = env
		fmt.Printf("使用 %s 环境变量配置文件: %s\n", internal.ConfigEnv, configPath)
		return
	}

	// 根据 gin 模式选择配置文件
	switch gin.Mode() {
	case gin.DebugMode:
		configPath = internal.ConfigDebugFile
	case gin.ReleaseMode:
		configPath = internal.ConfigReleaseFile
	case gin.TestMode:
		configPath = internal.ConfigTestFile
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); err != nil || os.IsNotExist(err) {
		configPath = internal.ConfigDefaultFile
		fmt.Printf("指定的配置文件不存在，使用默认配置文件: %s\n", configPath)
	} else {
		fmt.Printf("使用 gin %s 模式配置文件: %s\n", gin.Mode(), configPath)
	}

	return
}
