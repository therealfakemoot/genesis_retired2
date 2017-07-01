package main

import (
	"github.com/apex/log"
	genesis "github.com/therealfakemoot/genesis/app"
	cmd "github.com/therealfakemoot/genesis/cmd"
)

func main() {
	log.Info("Starting genesis")

	genesis.SetupLogging()
	genesis.LoadConfig()
	//logCtx := log.WithFields(log.Fields{
	//"name": "butts",
	//})

	cmd.Execute()
}
