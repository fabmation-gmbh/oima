package internal

import (
	"github.com/spf13/viper"

	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/config"
)


func GetConfig() config.Configuration {
	var conf config.Configuration

	err := viper.Unmarshal(&conf)
	if err != nil {
		Log.PanicF("unable to decode into struct, %v", err.Error())
	}

	return conf
}