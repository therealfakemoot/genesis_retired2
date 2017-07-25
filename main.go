package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	genesis "github.com/therealfakemoot/genesis/app"
	cmd "github.com/therealfakemoot/genesis/cmd"
	"os"
)

func main() {

	l := logrus.New()
	l.Out = os.Stdout

	genesis.LoadConfig(l)
	//logCtx := log.WithFields(log.Fields{
	//"name": "butts",
	//})

	l.WithFields(logrus.Fields{
		"Settings": viper.AllSettings(),
	}).Info("All settings")

	cmd.Execute()
}
