package cmd

import (
	"fmt"
	"os"

	"github.com/JNSAPH/db-seed-runner/internal"
	"github.com/JNSAPH/db-seed-runner/internal/db/postgres"
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

		logrus.Info("Starting db-seed-runner...")

		level, err := logrus.ParseLevel(internal.AppConfig.Log.Level)
		if err != nil {
			logrus.Warnf("Invalid log level '%s', defaulting to info", internal.AppConfig.Log.Level)
			level = logrus.InfoLevel
		}
		logrus.SetLevel(level)

		// Get the selected database engine
		engine := internal.AppConfig.Database.Engine
		logrus.Infof("Using database engine: %s", engine)

		// DB Creds from env vars
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")

		// Get Files in SQL Directory
		directory := internal.AppConfig.Seed.SqlFilesDirectory

		// Get all .sql files from the directory
		sqlFiles, err := internal.GetSQLFiles(directory)
		if err != nil {
			logrus.Fatalf("Error reading SQL files from directory '%s': %v", directory, err)
		}
		if len(sqlFiles) == 0 {
			logrus.Warnf("No SQL files found in directory '%s'", directory)
		} else {
			logrus.Infof("Found %d SQL files in directory '%s' (%v)", len(sqlFiles), directory, sqlFiles)
		}

		failed_files := []string{}

		switch engine {
		case "postgres":
			client := postgres.NewPostgres(&postgres.Credentials{
				Host:     internal.AppConfig.Database.Host,
				Port:     internal.AppConfig.Database.Port,
				User:     dbUser,
				Password: dbPassword,
				Database: internal.AppConfig.Database.Database,
				Sslmode:  internal.AppConfig.Database.Sslmode,
			})

			err := client.Connect()
			if err != nil {
				logrus.Fatalf("Failed to connect to the database: %v", err)
			}
			logrus.Info("Successfully connected to the PostgreSQL database")
			defer client.Close()

			for _, file := range sqlFiles {
				logrus.Infof("Executing SQL file: %s", file)
				if err := client.RunScript(file); err != nil {
					logrus.Errorf("Failed to execute SQL file '%s': %v", file, err)
					failed_files = append(failed_files, file)
					continue
				}
				logrus.Infof("Successfully executed SQL file: %s", file)
			}
		default:
			logrus.Fatalf("Unsupported database engine: %s", engine)
		}

		logrus.Info("Database seeding completed.")

		logrus.Print("===================================")
		if len(failed_files) > 0 {
			logrus.Errorf("Some files failed to execute (%d):", len(failed_files))
			for _, f := range failed_files {
				logrus.Errorf(" - %s", f)
			}
			logrus.Print("===================================")
			os.Exit(1) // mark kubernetes job as failed
		} else {
			logrus.Info("All files executed successfully!")
		}
		logrus.Print("===================================")

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
