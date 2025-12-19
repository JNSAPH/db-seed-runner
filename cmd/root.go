package cmd

import (
	"fmt"

	"github.com/JNSAPH/db-seed-runner/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var devFlag bool
var prodFlag bool

var rootCmd = &cobra.Command{
	Use:   "db-seed-runner",
	Short: "Go-based PostgreSQL seeding runner with Helm chart support",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		internal.SetupLogger()
		internal.SetupConfig()

		level, err := logrus.ParseLevel(internal.AppConfig.Log.Level)
		if err != nil {
			logrus.Warnf("Invalid log level '%s', defaulting to info", internal.AppConfig.Log.Level)
			level = logrus.InfoLevel
		}
		logrus.SetLevel(level)

		engine := internal.AppConfig.Database.Engine
		logrus.Infof("Using database engine: %s", engine)

		switch engine {
		default:
			logrus.Fatalf("Unsupported database engine: %s", engine)
		}

		logrus.Info("Starting db-seed-runner...")

	},
}

func init() {
	rootCmd.Flags().StringVarP(&internal.ConfigPath, "config", "c", fmt.Sprintf("%v/%v", internal.DefaultConfigPath, internal.DefaultConfigName), "Path to the config.yaml")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
