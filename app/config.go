package genesis

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

// LoadConfig sets Viper up to read from
// the env vars, config file, and CLI flags.
func LoadConfig(l *logrus.Logger) {

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("GENESIS")

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
