package config

import "github.com/spf13/viper"

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
}
