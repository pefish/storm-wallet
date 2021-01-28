package main

import (
	"github.com/pefish/go-commander"
	go_logger "github.com/pefish/go-logger"
	"wallet-storm-wallet/cmd/storm-wallet/default"
)

func main() {
	myCommander := commander.NewCommander("storm-wallet", "v0.0.1", "")
	myCommander.RegisterDefaultSubcommand(&_default.DefaultCommand{})

	err := myCommander.Run()
	if err != nil {
		go_logger.Logger.Error(err)
	}
}
