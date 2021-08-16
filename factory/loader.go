package factory

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhigang/retailcloud-admin/config"
)

var globalConfig *config.Config
var onceLoad sync.Once

// GlobalConfig return a config form config/*.yml
func GlobalConfig() *config.Config {
	onceLoad.Do(func() {
		globalConfig = &config.Config{}
		config := viper.New()
		config.AddConfigPath("./config")
		config.SetConfigName("config")

		if err := config.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := config.Unmarshal(&globalConfig); err != nil {
			log.Error(err)
		}

		if level, err := log.ParseLevel(globalConfig.Log.Level); err == nil {
			log.SetLevel(level)
		} else {
			log.SetLevel(log.DebugLevel)
		}

		log.Debugf("Use config: %+v", globalConfig)
	})
	return globalConfig
}
