package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sighupio/opa-notary-connector/internal/handlers"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks image passed as argument against configuration.",
	Long:  `Returns 0 and the image if image is signed, 1 if it is not.`,
	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]

		log := logrus.NewEntry(&logrus.Logger{})
		_, image, err := handlers.CheckImage(image, globalConfig.GetConfig(), globalConfig.TrustRootDir, log)

		if err != nil {
			log.WithError(err).Fatalf("there was an error while processing %+v", image)
		}
		fmt.Print(image)
	},
}
