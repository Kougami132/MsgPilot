package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                string `mapstructure:"APP_ENV"`
	Port                  string `mapstructure:"PORT"`
	ContextTimeout        int    `mapstructure:"CONTEXT_TIMEOUT"`
	AccessTokenExpiryHour int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret     string `mapstructure:"ACCESS_TOKEN_SECRET"`
	FrontendPath		  string `mapstructure:"FRONTEND_PATH"`
}

func NewEnv() *Env {
	// 设置默认值
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("CONTEXT_TIMEOUT", 10)
	viper.SetDefault("ACCESS_TOKEN_EXPIRY_HOUR", 720)
	viper.SetDefault("ACCESS_TOKEN_SECRET", "secret")
	viper.SetDefault("FRONTEND_PATH", "./static")

	// 尝试读取 .env 文件，如果不存在则使用默认值
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("未找到 .env 文件，使用默认配置")
	}

	// 也可以从环境变量中读取（优先级更高）
	viper.AutomaticEnv()

	env := Env{}
	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("环境变量无法加载: ", err)
	}

	return &env
}
