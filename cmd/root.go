package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/sighupio/opa-notary-connector/internal/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "opa-notary-connector",
		Short: "Start the server, loading configs",

		// main function, reading config, setting up the router and starting the server
		RunE: func(cmd *cobra.Command, args []string) error {

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
			return r.Run(globalConfig.BindAddress)
		},
	}
)

func init() {
	// flags set for the root command and all its subcommands
	rootCmd.PersistentFlags().StringVarP(&globalConfig.TrustRootDir, "trust-root-dir", "d", "/etc/opa-notary-connector/.trust", "Notary trust local cache directory.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.ConfigPath, "config", "c", "/etc/opa-notary-connector/trust.yaml", "Config file location.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.LogLevel, "verbosity", "v", "info", "Log level (one of fatal, error, warn, info or debug)")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.BindAddress, "listen-address", "l", ":8443", "Address the service should bind to.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.Mode, "mode", "m", gin.ReleaseMode, fmt.Sprintf("Set mode for gin and logger (%s, %s)", gin.ReleaseMode, gin.TestMode))
	viper.SetConfigFile(globalConfig.ConfigPath)

	cobra.OnInitialize(func() {
		logrus.SetReportCaller(true)
		level, err := logrus.ParseLevel(globalConfig.LogLevel)
		if err != nil {
			logrus.WithField("logLevel", globalConfig.LogLevel).WithError(err).Fatal("Log level not parsable")
			return
		}
		logrus.SetLevel(level)

		switch globalConfig.Mode {
		case gin.DebugMode:
			logrus.SetFormatter(new(logrus.TextFormatter))
			gin.SetMode(gin.DebugMode)
		case gin.ReleaseMode:
			logrus.SetFormatter(new(logrus.JSONFormatter))
			gin.SetMode(gin.ReleaseMode)
		default:
			logrus.SetFormatter(new(logrus.JSONFormatter))
			gin.SetMode(gin.ReleaseMode)
		}
	})
}

// Execute runs the root command from cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Error executing root command")
	}
}
