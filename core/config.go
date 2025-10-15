package core

import (
	"blog-server/config"
	"blog-server/core/internal"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var once sync.Once

func InitConfig() config.Config {
	var conf config.Config
	once.Do(func() {
		configPath := resolveConfigPath()

		v := viper.New()
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")
		v.AutomaticEnv() // 允许通过环境变量覆盖配置项

		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Errorf("❌ 读取配置文件失败: %w", err))
		}

		if err := v.Unmarshal(&conf); err != nil {
			panic(fmt.Errorf("❌ 解析配置文件失败: %w", err))
		}

		fmt.Printf("✅ 配置文件加载成功: %s\n", configPath)

		// 启动热更新监听
		watchConfig(v, &conf)
	})

	return conf
}

func resolveConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "c", "", "指定配置文件路径")
	flag.Parse()

	// 优先级 1：命令行参数
	if configPath != "" {
		fmt.Printf("📦 使用命令行参数指定配置文件: %s\n", configPath)
		return configPath
	}

	// 优先级 2：环境变量
	if envPath := os.Getenv(internal.ConfigEnv); envPath != "" {
		fmt.Printf("🌍 使用环境变量 %s 指定配置文件: %s\n", internal.ConfigEnv, envPath)
		return envPath
	}

	// 优先级 3：根据 gin 模式选择默认配置文件
	switch gin.Mode() {
	case gin.ReleaseMode:
		configPath = internal.ConfigReleaseFile
	case gin.TestMode:
		configPath = internal.ConfigTestFile
	default:
		configPath = internal.ConfigDebugFile
	}

	// 优先级 4：如果不存在则使用 fallback 默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  配置文件 [%s] 不存在，使用默认配置: %s\n", configPath, internal.ConfigDefaultFile)
		configPath = internal.ConfigDefaultFile
	} else {
		fmt.Printf("🔧 使用 gin %s 模式配置文件: %s\n", gin.Mode(), configPath)
	}

	return configPath
}

// watchConfig 监听配置文件变化
func watchConfig(v *viper.Viper, conf *config.Config) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("🔄 检测到配置文件变更: %s\n", e.Name)
		var newConf config.Config
		if err := v.Unmarshal(&newConf); err != nil {
			fmt.Printf("❌ 配置文件热更新解析失败: %v\n", err)
			return
		}

		// 安全更新：仅替换整个结构体
		*conf = newConf
		fmt.Println("✅ 配置文件热更新成功")
	})
}
