package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"

	conf "github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	globalConfig = conf.NewGlobalConfig()
	versionCmd   = &cobra.Command{
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

		// main function, reading config, setting up the router and starting the server
		Run: rootCmdFunc,
	}
)

func Execute() {
	// flags set for the root command and all its subcommands
	rootCmd.PersistentFlags().StringVarP(&globalConfig.TrustRootDir, "trust-root-dir", "d", "/etc/opa-notary-connector/.trust", "Notary trust local cache directory.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.ConfigPath, "config", "c", "/etc/opa-notary-connector/trust.yaml", "Config file location.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.LogLevel, "verbosity", "v", "info", "Log level (one of fatal, error, warn, info or debug)")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.BindAddress, "listen-address", "l", ":8443", "Address the service should bind to.")
	rootCmd.PersistentFlags().StringVarP(&globalConfig.Mode, "mode", "m", gin.ReleaseMode, fmt.Sprintf("Set mode for gin and logger (%s, %s)", gin.ReleaseMode, gin.TestMode))
	rootCmd.AddCommand(defaultConfig)
	rootCmd.AddCommand(versionCmd)

	logrus.SetReportCaller(true)
	viper.SetConfigFile(globalConfig.ConfigPath)
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
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Error executing root command")
	}
}
