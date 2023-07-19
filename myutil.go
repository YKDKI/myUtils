package myUtils

import "github.com/spf13/viper"

// -conf 指定config文件路径 默认 config/local.yml
func New() (config *viper.Viper, log *Logger) {
	config = NewConfig()
	log = NewLog(config)
	InitHttpClient(config, log)
	return
}
