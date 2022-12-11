package console

import (
	"os"

	"github.com/luckyAkbar/yori/internal/config"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd :nodoc:
var RootCmd = &cobra.Command{
	Use: "Himatro API",
}

// Execute execute console command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)

		os.Exit(1)
	}
}

func init() {
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	if config.Env() == "development" {
		formatter = runtime.Formatter{
			ChildFormatter: &logrus.TextFormatter{
				ForceColors:   true,
				FullTimestamp: true,
			},
			Line: true,
			File: true,
		}
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)

}
