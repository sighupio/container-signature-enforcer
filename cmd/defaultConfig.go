package cmd

import (
	"fmt"

	conf "github.com/sighupio/opa-notary-connector/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(defaultConfig)
}

var defaultConfig = &cobra.Command{
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
