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
	Log struct {
		Level string `mapstructure:"level" validate:"required,oneof=debug info warn error fatal panic"`
	} `mapstructure:"log" validate:"required"`
	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required,min=1"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
	} `mapstructure:"database" validate:"required"`
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
