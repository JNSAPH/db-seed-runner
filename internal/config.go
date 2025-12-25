package internal

import (
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ConfigPath string
var ConfigName string

const DefaultConfigName = "config"
const DefaultConfigPath = "./config"

type config struct {
	Seed struct {
		Enabled           bool   `mapstructure:"enabled" validate:"required"`
		Debug             bool   `mapstructure:"debug"`
		SqlFilesDirectory string `mapstructure:"sqlMountPath" validate:"required_if=Enabled true"`
	} `mapstructure:"seed" validate:"required"`
	Database struct {
		Engine string `mapstructure:"engine" validate:"required,oneof=postgres"`
		Host   string `mapstructure:"host" validate:"required"`
		Port   int    `mapstructure:"port" validate:"required,min=1"`
	} `mapstructure:"database" validate:"required"`
	Log struct {
		Level string `mapstructure:"level" validate:"required,oneof=debug INFO warn error fatal panic trace"`
	} `mapstructure:"logging" validate:"required"`
}

var AppConfig *config

func SetupConfig() {
	ConfigPath, ConfigName = splitPath(ConfigPath)

	logrus.Info("Config path: ", ConfigPath)
	logrus.Info("Config name: ", ConfigName)

	viper.SetConfigName(ConfigName)
	viper.AddConfigPath(ConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Error reading config file: ", err)
	}

	var c config

	if err := viper.Unmarshal(&c); err != nil {
		logrus.Fatal("Error unmarshaling config: ", err)
	}

	validateConfig(c)

	AppConfig = &c
}

func validateConfig(c config) {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		logrus.Fatal("Error validating config: ", err)
	}

	logrus.Info("Configuration validated successfully")
}

func splitPath(fullPath string) (string, string) {
	dir := filepath.Dir(fullPath)
	file := strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))
	return dir, file
}
