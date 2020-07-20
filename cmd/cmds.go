package cmd

import (
	"os"

	"github.com/fsnotify/fsnotify"
	conf "github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sighupio/opa-notary-connector/internal/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ()

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

func rootCmdFunc(cmd *cobra.Command, args []string) {

	// start execution
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Error reading config file")
	}

	// load config
	if err = viper.Unmarshal(globalConfig.GetConfig()); err != nil {
		logrus.WithError(err).Fatal("Error unmarshalling config into struct")
	} else {
		logrus.WithFields(logrus.Fields{"config": globalConfig.GetConfig(), "file": globalConfig.ConfigPath}).Info("Read config at startup.")
	}

	startupLogger := logrus.WithField("phase", "startup")
	if err := globalConfig.GetConfig().Validate(startupLogger); err != nil {
		startupLogger.WithError(err).Fatal("Validation error for config")
	}
	// setup watch for config change and reload if config parseable
	viper.WatchConfig()
	viper.OnConfigChange(reloadConfig)

	// trustRootDir is the location in which notary library will store its local cache
	if err := os.MkdirAll(globalConfig.TrustRootDir, 0700); err != nil {
		logrus.WithError(err).WithField("directory", globalConfig.TrustRootDir).Fatal("Error creating directory.")
	}

	// setup the router
	r := handlers.SetupServer(globalConfig)
	err = r.Run(globalConfig.BindAddress)
}
