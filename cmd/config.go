package cmd

import (
	"github.com/fsnotify/fsnotify"
	conf "github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var globalConfig = conf.NewGlobalConfig()

func reloadConfig(e fsnotify.Event) {
	logrus.WithField("file", globalConfig.ConfigPath).Info("Config file modified.")
	newConfig := conf.Config{}
	if err := viper.Unmarshal(&newConfig); err != nil {
		logrus.WithError(err).Error("error unmarshalling new config")
		return
	}
	reloadLogger := logrus.WithField("phase", "reloadConfig")
	if err := newConfig.Validate(reloadLogger); err != nil {
		reloadLogger.WithError(err).Error("Error validating new config reloaded, fallbacking to previous one.")
		return
	}
	globalConfig.SetConfig(&newConfig)
	logrus.WithField("config", globalConfig.GetConfig()).Printf("New config loaded")
}
