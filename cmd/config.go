package cmd

import (
	"github.com/fsnotify/fsnotify"
	conf "github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var globalConfig = conf.NewGlobalConfig()

func reloadConfig(e fsnotify.Event) {
	reloadLogger := logrus.WithField("phase", "reloadConfig")
	reloadLogger.WithField("file", globalConfig.ConfigPath).Info("Config file modified.")
	newConfig := conf.Config{}
	if err := viper.Unmarshal(&newConfig); err != nil {
		logrus.WithError(err).Error("error unmarshalling new config")
		return
	}
	if err := newConfig.Validate(reloadLogger); err != nil {
		reloadLogger.WithError(err).Error("Error validating new config reloaded, fallbacking to previous one.")
		return
	}
	globalConfig.SetConfig(&newConfig)
	reloadLogger.WithField("config", globalConfig.GetConfig()).Printf("New config loaded")
}
