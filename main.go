package main

import (
	"github.com/apex/log"
	genesis "github.com/therealfakemoot/genesis/app"
	cmd "github.com/therealfakemoot/genesis/cmd"
)

func main() {

	genesis.SetupLogging()
	log.Info("Starting genesis")
	//logCtx := log.WithFields(log.Fields{
	//"name": "butts",
	//})

	cmd.Execute()
}
