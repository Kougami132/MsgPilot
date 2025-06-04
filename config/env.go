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
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("未找到 .env 文件: ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("环境变量无法加载: ", err)
	}

	return &env
}
