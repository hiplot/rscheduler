package config

import (
	"github.com/spf13/viper"
	"log"
)

func Init() {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	err := v.ReadInConfig()
	if err != nil {
		panic("read config failed, err: " + err.Error())
	}
	err = v.Unmarshal(&Config)
	if err != nil {
		panic("unmarshal config failed, err: " + err.Error())
	}
	log.Println("配置初始化成功")
}
