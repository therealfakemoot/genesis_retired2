package genesis

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	//"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

func LoadConfig() {
	var err error
	viper.SetDefault("mapDir", "maps")
	viper.SetDefault("extDirs", "ext")

	viper.SetConfigName(".genesis")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()

	if err != nil {
		log.WithError(err).Error("config load failed")
	}
}

func SetupLogging() {
	var err error

	logFileName := "log/" + time.Now().Format("2017-01-01-00:00:00") + ".log"
	logFile, err := os.Open(logFileName)
	defer logFile.Close()

	if err != nil {
		log.WithError(err).Error("opening log file failed")
	}

	log.SetHandler(text.New(logFile))
}
