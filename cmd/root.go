package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"

	conf "github.com/sighupio/opa-notary-connector/config"
	"github.com/sighupio/opa-notary-connector/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	// flags set for the root command and all its subcommands
	rootCmd.PersistentFlags().StringVarP(&globalConfig.TrustRootDir, "trust-root-dir", "d", "/etc/opa-notary-connector/.trust", "Notary trust local cache directory.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.ConfigPath, "config", "c", "/etc/opa-notary-connector/trust.yaml", "Config file location.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.LogLevel, "verbosity", "v", "info", "Log level (one of fatal, error, warn, info or debug)")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.BindAddress, "listen-address", "l", ":8443", "Address the service should bind to.")
	rootCmd.AddCommand(defaultConfig)
	rootCmd.AddCommand(versionCmd)
}

var (
	globalConfig = conf.NewGlobalConfig()

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
		},
	}

	defaultConfig = &cobra.Command{
		Use:   "defaultConfig",
		Short: "Print default config",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := yaml.Marshal(conf.Config{
				Repositories: conf.Repositories{
					conf.Repository{
						Name: "registry.test/.*",
						Trust: conf.Trust{
							Enabled:     true,
							TrustServer: "https://notary-server:4443",
							Signers: []*conf.Signer{
								{
									Role:      "targets/releases",
									PublicKey: "BASE64_PUBLIC_KEY_HERE",
								},
							},
						},
					},
				},
			})
			if err != nil {
				logrus.WithError(err).Fatal("Error marshalling default config", err.Error())
			}
			fmt.Printf("\n%s\n", string(config))
		},
	}

	rootCmd = &cobra.Command{
		Use:   "opa-notary-connector",
		Short: "Start the server, loading configs",

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// needed settings
			viper.SetConfigFile(globalConfig.ConfigPath)
			level, err := logrus.ParseLevel(globalConfig.LogLevel)
			if err != nil {
				logrus.WithField("logLevel", globalConfig.LogLevel).WithError(err).Fatal("Log level not parsable")
			}
			logrus.SetLevel(level)
			logrus.SetReportCaller(true)
			logrus.SetFormatter(new(logrus.JSONFormatter))
		},
		// main function, reading config, setting up the router and starting the server
		Run: func(cmd *cobra.Command, args []string) {

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
			r := gin.New()
			r.Use(ginLogger())
			r.Use(recoveryLogger())
			//TODO move to customRecovery to log with logrus on panic, will be available in next gin release
			//r.Use(gin.CustomRecovery())
			r.POST("/checkImage", handlers.CheckImageHandlerBuilder(globalConfig))
			r.GET("/healthz", func(c *gin.Context) {
				c.String(http.StatusOK, "this is fine")
			})

			err = r.Run(globalConfig.BindAddress)
		},
	}
)

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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Error executing root command")
	}
}
