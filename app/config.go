package app

import (
	"github.com/sirupsen/logrus"
	//"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Logger is a stub type for now. One day,
// it will assist with logging to multiple targets
// i.e. file AND terminal simultaneously.
type Logger struct {
	Term *logrus.Logger
}

// DumpSettings will print all currently set configuration values
// using the given logger.
func DumpSettings(l *logrus.Logger) {

	l.WithFields(logrus.Fields{
		"Settings": viper.AllSettings(),
	}).Debug("All settings")
}

// LoadConfig sets Viper up to read from
// the env vars, config file, and CLI flags.
func LoadConfig(l *logrus.Logger) {

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("GENESIS")
	viper.AutomaticEnv()

	viper.SetDefault("mapDir", "maps")
	viper.SetDefault("extDirs", "ext")

	viper.SetConfigName(".genesis")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		l.WithError(err).Error("Config load failed")
	}
}
