package main

import (
	"github.com/apex/log"
	cmd "github.com/therealfakemoot/genesis/cmd"
	"github.com/therealfakemoot/genesis/lib"
)

func main() {

	genesis.SetupLogging()
	log.Info("Starting genesis")
	//logCtx := log.WithFields(log.Fields{
	//"name": "butts",
	//})

	cmd.Execute()
}
