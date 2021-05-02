package commons

import (
	"github.com/spf13/viper"
)

func ViperInit() {
	viper.AddConfigPath("./conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig() 
}